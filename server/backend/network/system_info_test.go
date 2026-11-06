package network

import (
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

const sampleSystemInfo = "dwMajorVersion: 10\ndwMinorVersion: 0\ndwBuildNumber: 19045\nwProductType: \x00\ntotalPhysicalMemory: 33444220928\navailablePhysicalMemory: 7638020096\nnumberOfProcessors: 16\nprocessorArchitecture: 9\n"

func TestSystemInfoResponseUpdatesSession(t *testing.T) {
	client := &models.Client{ID: "test", IPAddress: "test:1", SystemInfo: models.SystemInfo{Hostname: "lab", User: "tester"}}
	AddClient(client.IPAddress, client)
	defer RemoveClient(client.IPAddress)
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT}, Data: []byte(sampleSystemInfo)})
	snapshot, _ := GetClient(client.IPAddress)
	info := snapshot.SystemInfo
	if info.OS != "Windows 10.0 (build 19045)" || info.Arch != "x64" || info.Hostname != "lab" || info.User != "tester" {
		t.Fatalf("unexpected system info: %+v", info)
	}
	if info.TotalPhysicalMemory == nil || *info.TotalPhysicalMemory != 33444220928 || info.AvailablePhysicalMemory == nil || *info.AvailablePhysicalMemory != 7638020096 || info.NumberOfProcessors == nil || *info.NumberOfProcessors != 16 {
		t.Fatalf("missing hardware info: %+v", info)
	}
	if info.ProductType != nil {
		t.Fatal("empty product type must remain unknown")
	}
	if len(snapshot.Messages) != 1 || snapshot.Messages[0].Content != sampleSystemInfo {
		t.Fatal("raw response was not preserved")
	}
	payload, err := json.Marshal(GetAllClients())
	if err != nil || !strings.Contains(string(payload), `"build_number":19045`) {
		t.Fatalf("session JSON missing metadata: %s, %v", payload, err)
	}
	// A later response updates memory but preserves previous values for invalid fields.
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_SYS_INFO}, Data: []byte("availablePhysicalMemory: 0\nnumberOfProcessors: -1\ntotalPhysicalMemory: invalid\n")})
	updated, _ := GetClient(client.IPAddress)
	if *updated.SystemInfo.AvailablePhysicalMemory != 0 || *updated.SystemInfo.NumberOfProcessors != 16 || *updated.SystemInfo.TotalPhysicalMemory != 33444220928 {
		t.Fatal("partial update lost valid metadata")
	}
	if *snapshot.SystemInfo.AvailablePhysicalMemory != 7638020096 || len(snapshot.Messages) != 1 {
		t.Fatal("previous snapshot was mutated")
	}
}

func TestSystemInfoParsing(t *testing.T) {
	current := models.SystemInfo{OS: "Existing OS", Arch: "Existing arch"}
	for _, input := range []string{"ordinary output", "numberOfProcessors: 42", "dwMajorVersion: invalid\ndwMinorVersion: 0\ndwBuildNumber: 19045"} {
		got := mergeSystemInfo(current, []byte(input), true)
		if got.OS != current.OS || got.NumberOfProcessors != nil {
			t.Fatalf("ordinary/malformed output changed metadata: %+v", got)
		}
	}
	got := mergeSystemInfo(current, []byte(strings.ReplaceAll(sampleSystemInfo, "\n", "\\\r\n")), true)
	if got.BuildNumber == nil || *got.BuildNumber != 19045 {
		t.Fatal("CRLF/backslash line endings not parsed")
	}
	for code, want := range map[string]string{"0": "x86", "9": "x64", "12": "ARM64", "65535": "Unknown (65535)"} {
		info := mergeSystemInfo(models.SystemInfo{}, []byte("processorArchitecture: "+code), false)
		if info.Arch != want {
			t.Errorf("architecture %s: %s", code, info.Arch)
		}
	}
}

func TestConcurrentSessionSnapshots(t *testing.T) {
	client := &models.Client{ID: "race-test", IPAddress: "test:2"}
	AddClient(client.IPAddress, client)
	defer RemoveClient(client.IPAddress)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_SYS_INFO}, Data: []byte(sampleSystemInfo)})
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			json.Marshal(GetAllClients())
			snapshot, _ := GetClient(client.IPAddress)
			json.Marshal(snapshot)
		}
	}()
	wg.Wait()
}
