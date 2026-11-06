package network

import (
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func TestReconnectSamePeerIPPreservesSession(t *testing.T) {
	old := &models.Client{ID: "original-agent", IPAddress: "192.0.2.25:1000", IsActive: true, SystemInfo: models.SystemInfo{Hostname: "WIN-LAB", User: "lab-user"}, Messages: []models.ClientMessage{{Content: "saved message"}}, CommandResults: []models.CommandResult{{Command: "whoami", Output: "lab-user"}}}
	AddClient(old.IPAddress, old)
	defer RemoveClient(old.IPAddress)
	DisconnectClient(old)
	incoming := &models.Client{ID: "temporary-agent", IPAddress: "192.0.2.25:2000", IsActive: true, Type: "TCP", Listener: "127.0.0.1:4444", LastSeen: time.Now()}
	AddClient(incoming.IPAddress, incoming)
	defer RemoveClient(incoming.IPAddress)
	got, _ := GetClient(incoming.IPAddress)
	if got.ID != "original-agent" || !got.IsActive || got.SystemInfo.User != "lab-user" || len(got.Messages) != 1 || len(got.CommandResults) != 1 {
		t.Fatalf("session not restored: %+v", got)
	}
	if got.Listener != incoming.Listener || got.IPAddress != incoming.IPAddress || !got.LastSeen.Equal(incoming.LastSeen) {
		t.Fatal("new connection fields not retained")
	}
	if _, exists := GetClient(old.IPAddress); exists {
		t.Fatal("old duplicate entry remains")
	}
	DisconnectClient(old)
	got, _ = GetClient(incoming.IPAddress)
	if !got.IsActive {
		t.Fatal("stale disconnect affected resumed connection")
	}
}

func TestReconnectDoesNotMergeDifferentActiveOrAmbiguousPeers(t *testing.T) {
	for _, test := range []struct {
		name      string
		oldActive bool
		ambiguous bool
		newIP     string
	}{
		{"different IP", false, false, "192.0.2.51:2000"},
		{"active peer", true, false, "192.0.2.50:2000"},
		{"ambiguous offline peers", false, true, "192.0.2.50:2000"},
	} {
		t.Run(test.name, func(t *testing.T) {
			old := &models.Client{ID: "saved", IPAddress: "192.0.2.50:1000", IsActive: test.oldActive}
			AddClient(old.IPAddress, old)
			defer RemoveClient(old.IPAddress)
			if test.ambiguous {
				other := &models.Client{ID: "other", IPAddress: "192.0.2.50:1001"}
				AddClient(other.IPAddress, other)
				defer RemoveClient(other.IPAddress)
			}
			incoming := &models.Client{ID: "new", IPAddress: test.newIP, IsActive: true}
			AddClient(incoming.IPAddress, incoming)
			defer RemoveClient(incoming.IPAddress)
			got, _ := GetClient(incoming.IPAddress)
			if got.ID != "new" {
				t.Fatal("distinct or ambiguous connection merged")
			}
			if _, exists := GetClient(old.IPAddress); !exists {
				t.Fatal("saved record lost")
			}
		})
	}
}

func TestReconnectIPv6PortChange(t *testing.T) {
	saved := &models.Client{ID: "ipv6-saved", IPAddress: "[2001:db8::7]:1000"}
	AddClient(saved.IPAddress, saved)
	defer RemoveClient(saved.IPAddress)
	live := &models.Client{ID: "ipv6-new", IPAddress: "[2001:db8::7]:2000", IsActive: true}
	AddClient(live.IPAddress, live)
	defer RemoveClient(live.IPAddress)
	if live.ID != "ipv6-saved" {
		t.Fatal("IPv6 peer not resumed")
	}
}
