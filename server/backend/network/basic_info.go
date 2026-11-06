package network

import (
	"log"

	"github.com/Connduit/frazzledrip/server/backend/models"
	"github.com/Connduit/frazzledrip/server/backend/utils"
)

// Protected by clientsMutex; only populated for newly accepted sessions.
var basicInfoQueues = make(map[string][]string)

func beginBasicInfo(client *models.Client) {
	clientsMutex.Lock()
	if clients[client.IPAddress] != client || !client.IsActive || client.BasicInfoStatus != "" {
		clientsMutex.Unlock()
		return
	}
	client.BasicInfoStatus = "collecting"
	basicInfoQueues[client.IPAddress] = []string{"whoami", "ipconfig /all", "systeminfo"}
	clientsMutex.Unlock()
	sendNextBasicInfo(client)
}

func sendNextBasicInfo(client *models.Client) {
	clientsMutex.Lock()
	queue, exists := basicInfoQueues[client.IPAddress]
	if clients[client.IPAddress] != client || !client.IsActive || !exists {
		clientsMutex.Unlock()
		return
	}
	if _, pending := pendingCommands[client.IPAddress]; pending {
		clientsMutex.Unlock()
		return
	}
	if len(queue) == 0 {
		client.BasicInfoStatus = "complete"
		delete(basicInfoQueues, client.IPAddress)
		clientsMutex.Unlock()
		return
	}
	command := queue[0]
	basicInfoQueues[client.IPAddress] = queue[1:]
	id := utils.GenerateMessageID()
	pendingCommands[client.IPAddress] = pendingCommand{ID: id, Command: command, Automatic: true}
	clientsMutex.Unlock()
	msg := models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_EXECUTE_COMMAND, MessageID: id, DataSize: uint32(len(command))}, Data: []byte(command)}
	if err := sendBinaryMessage(client, msg); err != nil {
		clientsMutex.Lock()
		if clients[client.IPAddress] == client && pendingCommands[client.IPAddress].ID == id {
			delete(pendingCommands, client.IPAddress)
			delete(basicInfoQueues, client.IPAddress)
			client.BasicInfoStatus = "failed"
		}
		clientsMutex.Unlock()
		log.Printf("Initial host info for %s failed: %v", client.ID, err)
	}
}

// Call with clientsMutex held. Never replay commands after a reconnect.
func interruptBasicInfo(client *models.Client) {
	if client.BasicInfoStatus == "collecting" {
		client.BasicInfoStatus = "interrupted"
	}
	delete(basicInfoQueues, client.IPAddress)
}
