package network

import (
	"net"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func TestInitialCollectionFailureAndInterruption(t *testing.T) {
	t.Run("write failure", func(t *testing.T) {
		conn, peer := net.Pipe()
		peer.Close()
		defer conn.Close()
		client := &models.Client{ID: "failed-init", IPAddress: "test:init-failure", IsActive: true, Conn: conn}
		AddClient(client.IPAddress, client)
		defer RemoveClient(client.IPAddress)
		beginBasicInfo(client)
		saved, _ := GetClient(client.IPAddress)
		if saved.BasicInfoStatus != "failed" {
			t.Fatalf("write failure status: %s", saved.BasicInfoStatus)
		}
		clientsMutex.Lock()
		_, pending := pendingCommands[client.IPAddress]
		_, queued := basicInfoQueues[client.IPAddress]
		clientsMutex.Unlock()
		if pending || queued {
			t.Fatal("failed collection left pending work")
		}
	})
	t.Run("disconnect and resume", func(t *testing.T) {
		conn, peer := net.Pipe()
		defer peer.Close()
		defer conn.Close()
		peer.SetDeadline(time.Now().Add(3 * time.Second))
		client := &models.Client{ID: "saved-init", IPAddress: "192.0.2.99:1000", IsActive: true, Conn: conn}
		AddClient(client.IPAddress, client)
		defer RemoveClient(client.IPAddress)
		done := make(chan struct{})
		go func() { beginBasicInfo(client); close(done) }()
		first, err := readBinaryMessage(peer)
		if err != nil || string(first.Data) != "whoami" {
			t.Fatalf("first query: %v", err)
		}
		<-done
		if err := SendCommandToClient(client.IPAddress, "tasklist"); err == nil {
			t.Fatal("manual command interrupted collection")
		}
		if err := DisconnectSession(client.ID); err != nil {
			t.Fatal(err)
		}
		saved, _ := GetClient(client.IPAddress)
		if saved.BasicInfoStatus != "interrupted" {
			t.Fatal("collection not marked interrupted")
		}
		replacement := &models.Client{ID: "new-init", IPAddress: "192.0.2.99:2000", IsActive: true}
		AddClient(replacement.IPAddress, replacement)
		defer RemoveClient(replacement.IPAddress)
		beginBasicInfo(replacement)
		if replacement.ID != "saved-init" || replacement.BasicInfoStatus != "interrupted" {
			t.Fatal("reconnect replayed initialization")
		}
	})
}

func TestInitialCollectionWaitsForExistingCommand(t *testing.T) {
	conn, peer := net.Pipe()
	defer peer.Close()
	defer conn.Close()
	peer.SetDeadline(time.Now().Add(3 * time.Second))
	client := &models.Client{ID: "init-race", IPAddress: "test:init-race", IsActive: true, Conn: conn}
	AddClient(client.IPAddress, client)
	defer RemoveClient(client.IPAddress)
	clientsMutex.Lock()
	pendingCommands[client.IPAddress] = pendingCommand{ID: 55, Command: "tasklist"}
	clientsMutex.Unlock()
	beginBasicInfo(client)
	done := make(chan struct{})
	go func() {
		processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: 105}, Data: []byte("process listing")})
		close(done)
	}()
	query, err := readBinaryMessage(peer)
	if err != nil || string(query.Data) != "whoami" {
		t.Fatalf("initial queue did not advance after existing command: %v", err)
	}
	<-done
	DisconnectSession(client.ID)
}
