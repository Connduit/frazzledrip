package network

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func TestMergeCommandInfo(t *testing.T) {
	info := mergeCommandInfo(models.SystemInfo{}, "whoami", "LAB\\researcher\r\n")
	if info.User != `LAB\researcher` {
		t.Fatalf("identity: %+v", info)
	}
	info = mergeCommandInfo(info, "ipconfig /all", "Windows IP Configuration\r\n   Host Name . . . . . . . . . . . . : WIN-LAB\n   Primary Dns Suffix  . . . . . . . : lab.local\n   IPv4 Address. . . . . . . . . . . : 10.0.0.48(Preferred)\n   IPv4 Address. . . . . . . . . . . : 10.0.0.48(Preferred)\n   Link-local IPv6 Address . . . . . : fe80::123%4(Preferred)\n   Physical Address. . . . . . . . . : AA-BB-CC-DD-EE-FF\n   Default Gateway . . . . . . . . . : 10.0.0.1\n")
	if info.Hostname != "WIN-LAB" || info.Domain != "lab.local" || len(info.IPAddresses) != 2 || info.IPAddresses[0] != "10.0.0.48" || info.IPAddresses[1] != "fe80::123" || len(info.MACAddresses) != 1 {
		t.Fatalf("ipconfig: %+v", info)
	}
	info = mergeCommandInfo(info, "systeminfo", "Host Name: WIN-LAB\nOS Name: Microsoft Windows 10 Pro\nOS Version: 10.0.19045 N/A Build 19045\nSystem Manufacturer: Example\nSystem Model: Lab PC\nSystem Type: x64-based PC\nProcessor(s): 16 Processor(s) Installed.\nTotal Physical Memory: 31,895 MB\nAvailable Physical Memory: 7,284 MB\nSystem Boot Time: 11/06/2026, 10:00:00\nDomain: WORKGROUP\n")
	if info.OS != "Microsoft Windows 10 Pro" || *info.BuildNumber != 19045 || info.Arch != "x64" || *info.NumberOfProcessors != 16 || *info.TotalPhysicalMemory != 31895*(1<<20) || info.Model != "Lab PC" {
		t.Fatalf("systeminfo: %+v", info)
	}
	for _, command := range []string{"dir C:\\", "netstat -an", "tasklist", "net user", "whoami & echo other"} {
		next := mergeCommandInfo(info, command, "Host Name: WRONG\n")
		if next.Hostname != info.Hostname {
			t.Fatalf("unrecognized command overwrote facts: %s", command)
		}
	}
	for _, output := range []string{"", "ERROR: Access denied", "The command completed successfully.", "User Name    SID\nLAB\\tester  S-1-5"} {
		if got := mergeCommandInfo(info, "whoami", output); got.User != info.User {
			t.Fatalf("invalid identity replaced user: %q", output)
		}
	}
	if got := mergeCommandInfo(info, "ipconfig", "IPv4 Address: invalid\nPhysical Address: invalid"); len(got.IPAddresses) != 2 || len(got.MACAddresses) != 1 {
		t.Fatal("invalid addresses erased facts")
	}
}

func TestCommandCorrelationAndLegacyResult(t *testing.T) {
	server, peer := net.Pipe()
	defer server.Close()
	defer peer.Close()
	server.SetDeadline(time.Now().Add(5 * time.Second))
	peer.SetDeadline(time.Now().Add(5 * time.Second))
	client := &models.Client{ID: "correlation-test", IPAddress: "test:commands", Conn: server}
	AddClient(client.IPAddress, client)
	defer RemoveClient(client.IPAddress)
	sent := make(chan error, 1)
	go func() { sent <- SendCommandToClient(client.IPAddress, "whoami") }()
	command, err := readBinaryMessage(peer)
	if err != nil || string(command.Data) != "whoami" {
		t.Fatalf("command write: %v", err)
	}
	if err := <-sent; err != nil {
		t.Fatal(err)
	}
	if err := SendCommandToClient(client.IPAddress, "tasklist"); err == nil || !strings.Contains(err.Error(), "pending") {
		t.Fatalf("parallel legacy command not rejected: %v", err)
	}
	// Dedicated system info (legacy ID 99) must not consume a shell request.
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: 99}, Data: []byte(sampleSystemInfo)})
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: 105}, Data: []byte("LAB\\tester\r\n")})
	snapshot, _ := GetClient(client.IPAddress)
	if snapshot.SystemInfo.User != `LAB\tester` || len(snapshot.CommandResults) != 1 || snapshot.CommandResults[0].Command != "whoami" {
		t.Fatalf("uncorrelated legacy response: %+v", snapshot)
	}
	// Unmatched outputs remain logs and do not change metadata.
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: 105}, Data: []byte("WRONG\\user")})
	snapshot, _ = GetClient(client.IPAddress)
	if snapshot.SystemInfo.User != `LAB\tester` {
		t.Fatal("unmatched result changed identity")
	}
	// Modern clients can echo the actual request ID.
	clientsMutex.Lock()
	pendingCommands[client.IPAddress] = pendingCommand{ID: 12345, Command: "whoami"}
	clientsMutex.Unlock()
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: 12345}, Data: []byte("LAB\\newuser")})
	snapshot, _ = GetClient(client.IPAddress)
	if snapshot.SystemInfo.User != `LAB\newuser` || len(snapshot.CommandResults) != 1 {
		t.Fatal("request ID correlation/latest output failed")
	}
}

func TestWindowsIdentityWithSpaces(t *testing.T) {
	got := mergeCommandInfo(models.SystemInfo{}, "whoami", "nt authority\\system\r\n")
	if got.User != `nt authority\system` {
		t.Fatalf("system identity not recognized: %q", got.User)
	}
}

func TestLatestCommandOutputs(t *testing.T) {
	client := &models.Client{}
	for _, command := range []string{"dir C:\\", "netstat -an", "tasklist", "net user"} {
		storeCommandResult(client, command, "sample listing", time.Now())
	}
	storeCommandResult(client, "TASKLIST", "new listing", time.Now())
	if len(client.CommandResults) != 4 || client.CommandResults[0].Output != "new listing" {
		t.Fatal("latest results not deduplicated")
	}
}
