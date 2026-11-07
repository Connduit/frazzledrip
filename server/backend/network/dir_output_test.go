package network

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

const sampleDir = " Volume in drive C is Windows\r\n Directory of C:\\Lab\r\n\r\n11/06/2026  09:30 PM    <DIR>          .\r\n11/06/2026  09:30 PM    <DIR>          ..\r\n11/05/2026  12:00 AM    <DIR>          Docs\r\n11/06/2026  12:00 PM                5 report one.txt\r\n11/06/2026  09:30 PM    <SYMLINKD>     shortcut [C:\\Data]\r\n               1 File(s)              5 bytes\r\n"

func TestParseDirOutput(t *testing.T) {
	listing, err := parseDirOutput(sampleDir)
	if err != nil || listing.Path != `C:\Lab` || listing.Parent != `C:\` || len(listing.Entries) != 3 {
		t.Fatalf("parse: %+v %v", listing, err)
	}
	file := listing.Entries[1]
	if file.Name != "report one.txt" || file.Path != `C:\Lab\report one.txt` || *file.Size != 5 || file.ModifiedText != "11/06/2026 12:00 PM" || file.ModifiedSort != "2026-11-06T12:00:00.000Z" || file.ModifiedMS != nil {
		t.Fatalf("file metadata: %+v", file)
	}
	if listing.Entries[0].ModifiedSort != "2026-11-05T00:00:00.000Z" {
		t.Fatal("midnight parsed incorrectly")
	}
	if !listing.Entries[2].IsLink || !listing.Entries[2].IsDirectory || listing.Entries[2].Name != "shortcut" {
		t.Fatal("link parsed incorrectly")
	}
	if _, err := parseDirOutput("Access is denied."); err == nil {
		t.Fatal("access error reported as empty directory")
	}
	if _, err := parseDirOutput("Directory of C:\\Lab\n0 File(s) 0 bytes\n"); err != nil {
		t.Fatal("empty directory not recognized", err)
	}
	if _, err := parseDirOutput("Directory of C:\\Lab\nunsupported entry\n"); err == nil {
		t.Fatal("unsupported format reported as empty")
	}
	if windowsParent(`\\server\share`) != `\\server\share` || windowsParent(`C:\`) != `C:\` {
		t.Fatal("root parent should be itself")
	}
	unknown, _ := parseDirOutput("Directory of C:\\Lab\n31/12/2026  09:00  3 local.txt\n")
	if unknown.Entries[0].ModifiedText != "31/12/2026 09:00" || unknown.Entries[0].ModifiedSort != "" {
		t.Fatal("unsupported locale date should remain as raw text")
	}
}

func TestDirCommandQuoting(t *testing.T) {
	command, err := dirCommand(`C:\Research & Notes`)
	if err != nil || !strings.Contains(command, `"C:\Research & Notes"`) || !strings.Contains(command, "/-s") {
		t.Fatal("directory not quoted or recursion not disabled")
	}
	for _, path := range []string{"x\" & echo bad", `C:\%TEMP%`, "a\nb", `C:\*`, "a\x01"} {
		if _, err := dirCommand(path); err == nil {
			t.Fatalf("unsafe path accepted: %q", path)
		}
	}
}

func TestBrowseViaDirLegacyResponse(t *testing.T) {
	client, peer := directoryClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := make(chan directoryResult, 1)
	go func() {
		listing, err := BrowseDirectoryViaDir(ctx, client.IPAddress, `C:\Lab`)
		result <- directoryResult{Listing: listing, Err: err}
	}()
	request, err := readBinaryMessage(peer)
	if err != nil || request.Header.MessageType != models.MESSAGE_EXECUTE_COMMAND || !strings.Contains(string(request.Data), "dir /a") {
		t.Fatalf("not a shell dir request: %v", err)
	}
	processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: 105}, Data: []byte(sampleDir)})
	response := <-result
	if response.Err != nil || len(response.Listing.Entries) != 3 {
		t.Fatalf("legacy dir result not parsed: %v", response.Err)
	}
	snapshot, _ := GetClient(client.IPAddress)
	if snapshot.Messages[0].Type != "DIRECTORY_LISTING" {
		t.Fatal("directory result misclassified as a shell result")
	}
}
