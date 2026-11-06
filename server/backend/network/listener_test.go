package network

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func TestListenerStopAndDisconnect(t *testing.T) {
	defer StopReverseTCPListener()
	if err := StartReverseTCPListener(0); err != nil {
		t.Fatal(err)
	}
	running, port := ListenerStatus()
	if !running || port == 0 {
		t.Fatal("listener not ready after start")
	}
	if err := StartReverseTCPListener(0); err == nil {
		t.Fatal("duplicate start accepted")
	}
	peer, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()
	addr := peer.LocalAddr().String()
	defer RemoveClient(addr)
	var session *models.Client
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if found, ok := GetClient(addr); ok {
			session = found
			break
		}
		time.Sleep(time.Millisecond)
	}
	if session == nil {
		t.Fatal("connection was not registered")
	}
	if err := StopReverseTCPListener(); err != nil {
		t.Fatal(err)
	}
	if IsListening() {
		t.Fatal("listener still running")
	}
	if err := StopReverseTCPListener(); err != nil {
		t.Fatal("second stop must be harmless")
	}
	if extra, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), time.Second); err == nil {
		extra.Close()
		t.Fatal("stopped listener accepted a connection")
	}
	// Existing connections still work after stopping accepts.
	peer.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := peer.Write(serializeBinaryMessage(models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_HEARTBEAT}})); err != nil {
		t.Fatal(err)
	}
	deadline = time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		session, _ = GetClient(addr)
		if len(session.Messages) > 0 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if !session.IsActive || len(session.Messages) != 1 {
		t.Fatal("stopping listener interrupted existing session")
	}
	if err := DisconnectSession(session.ID); err != nil {
		t.Fatal(err)
	}
	var buffer [1]byte
	if _, err := peer.Read(buffer[:]); err == nil {
		t.Fatal("session socket still open")
	}
	offline, _ := GetClient(addr)
	if offline.IsActive || len(offline.Messages) != 1 {
		t.Fatal("disconnect failed to retain offline session")
	}
	if err := DisconnectSession("missing-session"); err == nil {
		t.Fatal("unknown disconnect accepted")
	}
	if err := StartReverseTCPListener(0); err != nil {
		t.Fatal("listener could not restart", err)
	}
	// Give the previous accept loop time to clean up without clearing the new listener.
	time.Sleep(10 * time.Millisecond)
	if !IsListening() {
		t.Fatal("old accept loop cleared the new listener")
	}
}
