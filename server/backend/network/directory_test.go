package network

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func directoryClient(t *testing.T) (*models.Client, net.Conn) {
	t.Helper()
	conn, peer := net.Pipe()
	client := &models.Client{ID: t.Name(), IPAddress: "test:" + t.Name(), IsActive: true, Conn: conn, BasicInfoStatus: "complete"}
	AddClient(client.IPAddress, client)
	peer.SetDeadline(time.Now().Add(3 * time.Second))
	t.Cleanup(func() { RemoveClient(client.IPAddress); conn.Close(); peer.Close() })
	return client, peer
}

func TestBrowseDirectoryCorrelation(t *testing.T) {
	client, peer := directoryClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := make(chan directoryResult, 1)
	go func() {
		listing, err := BrowseDirectory(ctx, client.IPAddress, `C:\Lab`)
		result <- directoryResult{Listing: listing, Err: err}
	}()
	request, err := readBinaryMessage(peer)
	if err != nil || request.Header.MessageType != models.MESSAGE_LIST_DIRECTORY || string(request.Data) != `C:\Lab` {
		t.Fatalf("bad browse request: %v", err)
	}
	if err := SendCommandToClient(client.IPAddress, "whoami"); err == nil {
		t.Fatal("shell command accepted while browsing")
	}
	payload := []byte(`{"path":"C:\\Lab","parent":"C:\\","entries":[{"name":"hello.txt","path":"C:\\Lab\\hello.txt","is_directory":false,"size":5,"modified_ms":1700000000000}]}`)
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_DIRECTORY_LISTING, MessageID: request.Header.MessageID + 1}, Data: payload})
	select {
	case <-result:
		t.Fatal("unmatched result consumed request")
	default:
	}
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_DIRECTORY_LISTING, MessageID: request.Header.MessageID}, Data: payload})
	got := <-result
	if got.Err != nil || got.Listing.Path != `C:\Lab` || len(got.Listing.Entries) != 1 || *got.Listing.Entries[0].Size != 5 || got.Listing.Entries[0].ModifiedMS == nil || *got.Listing.Entries[0].ModifiedMS != 1700000000000 {
		t.Fatalf("listing not returned: %+v", got)
	}
}

func TestBrowseDirectoryFailures(t *testing.T) {
	for _, scenario := range []string{"permission", "invalid JSON", "cancel", "disconnect"} {
		t.Run(scenario, func(t *testing.T) {
			client, peer := directoryClient(t)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			done := make(chan error, 1)
			go func() { _, err := BrowseDirectory(ctx, client.IPAddress, ""); done <- err }()
			request, err := readBinaryMessage(peer)
			if err != nil {
				t.Fatal(err)
			}
			switch scenario {
			case "permission":
				processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_ERROR_REPORT, MessageID: request.Header.MessageID}, Data: []byte("Access denied")})
			case "invalid JSON":
				processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_DIRECTORY_LISTING, MessageID: request.Header.MessageID}, Data: []byte("invalid")})
			case "cancel":
				cancel()
			case "disconnect":
				DisconnectSession(client.ID)
			}
			if err := <-done; err == nil {
				t.Fatal("failure was reported as success")
			}
			clientsMutex.Lock()
			_, pending := directoryRequests[client.IPAddress]
			clientsMutex.Unlock()
			if pending {
				t.Fatal("request not cleaned up")
			}
		})
	}
	if _, err := BrowseDirectory(context.Background(), "missing", "x"); err == nil {
		t.Fatal("offline browse accepted")
	}
	if _, err := BrowseDirectory(context.Background(), "missing", strings.Repeat("x", 2049)); err == nil {
		t.Fatal("oversized path accepted")
	}
}
