package network

import (
	"net"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func TestDisconnectRetainsSessionData(t *testing.T) {
	conn, peer := net.Pipe()
	client := &models.Client{ID: "retained", IPAddress: "test:offline", Conn: conn, IsActive: true, SystemInfo: models.SystemInfo{Hostname: "WIN-LAB", User: "lab-user"}, Messages: []models.ClientMessage{{Content: "saved response"}}, CommandResults: []models.CommandResult{{Command: "whoami", Output: "lab-user"}}}
	AddClient(client.IPAddress, client)
	defer RemoveClient(client.IPAddress)
	defer peer.Close()
	clientsMutex.Lock()
	pendingCommands[client.IPAddress] = pendingCommand{ID: 77, Command: "tasklist"}
	clientsMutex.Unlock()
	done := make(chan struct{})
	go func() { handleClientConnection(client); close(done) }()
	peer.Close()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("disconnect handler did not finish")
	}
	saved, exists := GetClient(client.IPAddress)
	if !exists || saved.IsActive || saved.Conn != nil || saved.SystemInfo.Hostname != "WIN-LAB" || len(saved.Messages) != 1 || len(saved.CommandResults) != 1 {
		t.Fatalf("session not retained offline: %+v", saved)
	}
	if err := SendCommandToClient(client.IPAddress, "whoami"); err == nil {
		t.Fatal("offline command accepted")
	}
	clientsMutex.Lock()
	_, pending := pendingCommands[client.IPAddress]
	clientsMutex.Unlock()
	if pending {
		t.Fatal("offline session retained pending command")
	}
}

func TestOldDisconnectDoesNotClearReplacement(t *testing.T) {
	old := &models.Client{IPAddress: "test:replacement", IsActive: true}
	replacement := &models.Client{IPAddress: old.IPAddress, IsActive: true}
	AddClient(old.IPAddress, old)
	AddClient(replacement.IPAddress, replacement)
	defer RemoveClient(old.IPAddress)
	DisconnectClient(old)
	saved, _ := GetClient(old.IPAddress)
	if !saved.IsActive {
		t.Fatal("stale disconnect marked replacement offline")
	}
}
