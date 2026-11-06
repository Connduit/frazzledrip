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

	// New sessions collect facts sequentially before accepting manual commands.
	peer.SetDeadline(time.Now().Add(3 * time.Second))
	initial := []struct{ command, output string }{
		{"whoami", "LAB\\tester\r\n"},
		{"ipconfig /all", "Host Name . . . : WIN-LAB\nIPv4 Address. . : 10.0.0.48(Preferred)\n"},
		{"systeminfo", "OS Name: Microsoft Windows 10 Pro\nOS Version: 10.0.19045 N/A Build 19045\n"},
	}
	for _, query := range initial {
		request, err := readBinaryMessage(peer)
		if err != nil || string(request.Data) != query.command {
			t.Fatalf("initial query %q: %v", query.command, err)
		}
		output := []byte(query.output)
		_, err = peer.Write(serializeBinaryMessage(models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: 105, DataSize: uint32(len(output))}, Data: output}))
		if err != nil {
			t.Fatal(err)
		}
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
		if len(session.Messages) == 4 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if session.BasicInfoStatus != "complete" || session.SystemInfo.Hostname != "WIN-LAB" || session.SystemInfo.User != `LAB\tester` {
		t.Fatalf("initial metadata not collected: %+v", session)
	}
	if !session.IsActive || len(session.Messages) != 4 {
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
	if offline.IsActive || len(offline.Messages) != 4 {
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
