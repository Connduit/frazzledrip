package network

import (
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
	"github.com/Connduit/frazzledrip/server/backend/state"
)

func TestTrackedTaskWithMockClient(t *testing.T) {
	store, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	previous := persistentState
	persistentState = store
	defer func() { persistentState = previous; store.Close() }()
	conn, peer := net.Pipe()
	defer conn.Close()
	defer peer.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	peer.SetDeadline(time.Now().Add(3 * time.Second))
	client := &models.Client{ID: "test-persistent", IPAddress: "test:persistent", Conn: conn, IsActive: true, BasicInfoStatus: "complete"}
	AddClient(client.IPAddress, client)
	defer RemoveClient(client.IPAddress)
	sent := make(chan error, 1)
	go func() { sent <- SendCommandToClient(client.IPAddress, "test-command") }()
	request, err := readBinaryMessage(peer)
	if err != nil {
		t.Fatal(err)
	}
	if err := <-sent; err != nil {
		t.Fatal(err)
	}
	// Test-only response, not command execution. Legacy clients return 105.
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_ERROR_REPORT, MessageID: 105}, Data: []byte("TEST failure")})
	tasks, err := store.Tasks(client.ID)
	if err != nil || len(tasks) != 1 || tasks[0].Status != "failed" || tasks[0].RequestID != request.Header.MessageID || tasks[0].Output != "TEST failure" {
		t.Fatalf("task not correlated: %+v %v", tasks, err)
	}
	DisconnectSession(client.ID)
	records, _ := store.LoadSessions()
	if len(records) != 1 || len(records[0].Messages) != 1 {
		t.Fatal("session history was not persisted")
	}
}
