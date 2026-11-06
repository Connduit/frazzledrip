package network

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func TestLinkSavedSession(t *testing.T) {
	conn, peer := net.Pipe()
	defer conn.Close()
	defer peer.Close()
	now := time.Now()
	saved := &models.Client{ID: "saved-link", IPAddress: "test:old", SystemInfo: models.SystemInfo{Hostname: "WIN-LAB", User: "old-user"}, CommandResults: []models.CommandResult{{Command: "whoami", Output: "old-user", ReceivedAt: now.Add(-time.Minute)}, {Command: "tasklist", Output: "saved processes", ReceivedAt: now.Add(-time.Minute)}}}
	live := &models.Client{ID: "live-link", IPAddress: "test:new", Conn: conn, IsActive: true, LastSeen: now, SystemInfo: models.SystemInfo{User: "new-user"}, Messages: []models.ClientMessage{{Timestamp: now, Content: "new message"}}, CommandResults: []models.CommandResult{{Command: "WHOAMI", Output: "new-user", ReceivedAt: now}}}
	for i := 0; i < 101; i++ {
		saved.Messages = append(saved.Messages, models.ClientMessage{Timestamp: now.Add(time.Duration(i-101) * time.Second), Content: fmt.Sprint(i)})
	}
	AddClient(saved.IPAddress, saved)
	AddClient(live.IPAddress, live)
	defer RemoveClient(saved.IPAddress)
	defer RemoveClient(live.IPAddress)
	clientsMutex.Lock()
	pendingCommands[live.IPAddress] = pendingCommand{ID: 123, Command: "ipconfig"}
	clientsMutex.Unlock()
	if err := LinkSavedSession(saved.ID, live.ID); err != nil {
		t.Fatal(err)
	}
	got, _ := GetClient(live.IPAddress)
	if got.ID != live.ID || got.Conn != conn || !got.IsActive || !got.LastSeen.Equal(now) {
		t.Fatal("live connection changed")
	}
	if got.SystemInfo.Hostname != "WIN-LAB" || got.SystemInfo.User != "new-user" {
		t.Fatalf("metadata priority: %+v", got.SystemInfo)
	}
	if len(got.Messages) != 100 || got.Messages[99].Content != "new message" || got.Messages[0].Content != "2" {
		t.Fatal("history order/retention incorrect")
	}
	if len(got.CommandResults) != 2 || got.CommandResults[0].Output != "new-user" {
		t.Fatal("latest command outputs not merged")
	}
	if _, exists := GetClient(saved.IPAddress); exists {
		t.Fatal("saved record still present")
	}
	clientsMutex.Lock()
	pending := pendingCommands[live.IPAddress]
	clientsMutex.Unlock()
	if pending.ID != 123 {
		t.Fatal("live pending command changed")
	}
	DisconnectClient(saved)
	if got, _ := GetClient(live.IPAddress); !got.IsActive {
		t.Fatal("old disconnect affected linked live connection")
	}
	if err := LinkSavedSession(saved.ID, live.ID); err == nil {
		t.Fatal("stale link accepted")
	}
}

func TestLinkRejectsInvalidPair(t *testing.T) {
	conn, peer := net.Pipe()
	defer conn.Close()
	defer peer.Close()
	one := &models.Client{ID: "link-one", IPAddress: "test:one", Conn: conn, IsActive: true}
	two := &models.Client{ID: "link-two", IPAddress: "test:two", Conn: conn, IsActive: true}
	AddClient(one.IPAddress, one)
	AddClient(two.IPAddress, two)
	defer RemoveClient(one.IPAddress)
	defer RemoveClient(two.IPAddress)
	for _, pair := range [][2]string{{one.ID, two.ID}, {one.ID, one.ID}, {"missing", two.ID}} {
		if err := LinkSavedSession(pair[0], pair[1]); err == nil {
			t.Fatalf("invalid pair accepted: %v", pair)
		}
	}
	clientsMutex.Lock()
	one.IsActive = false
	two.IsActive = false
	clientsMutex.Unlock()
	if err := LinkSavedSession(one.ID, two.ID); err == nil {
		t.Fatal("offline target accepted")
	}
	if _, exists := GetClient(one.IPAddress); !exists {
		t.Fatal("failed link removed saved record")
	}
}
