package network

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func TestDirectoryCachePersistenceAndIsolation(t *testing.T) {
	file := filepath.Join(t.TempDir(), "directories.json")
	cache, err := loadDirectoryCache(file)
	if err != nil {
		t.Fatal(err)
	}
	listing := &DirectoryListing{Path: `C:\Lab`, Parent: `C:\`, Entries: []DirectoryEntry{{Name: "first.txt"}}}
	if err := cache.put("agent", "192.0.2.1", "dir", "", listing); err != nil {
		t.Fatal(err)
	}
	if listing.CapturedAt.IsZero() {
		t.Fatal("missing capture timestamp")
	}
	loaded, err := loadDirectoryCache(file)
	if err != nil {
		t.Fatal(err)
	}
	for _, directory := range []string{"", `c:/LAB/`, `C:\Lab`} {
		got, ok := loaded.lookup("agent", "192.0.2.1", "dir", directory, false, nil)
		if !ok || !got.Cached || got.Entries[0].Name != "first.txt" {
			t.Fatalf("missing persisted alias %q", directory)
		}
		got.Entries[0].Name = "changed"
	}
	if _, ok := loaded.lookup("other", "192.0.2.2", "dir", "", false, nil); ok {
		t.Fatal("other session reused cache")
	}
	if _, ok := loaded.lookup("agent", "192.0.2.1", "native", "", false, nil); ok {
		t.Fatal("listing sources mixed")
	}
	if directoryKey("dir", `C:`) == directoryKey("dir", `C:\`) {
		t.Fatal("drive-relative and root paths conflated")
	}
	if directoryKey("dir", `\\host\share\folder`) != directoryKey("dir", "//HOST/share/folder/") {
		t.Fatal("UNC aliases differ")
	}
	data := []byte("not JSON")
	if err := os.WriteFile(file, data, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadDirectoryCache(file); err == nil {
		t.Fatal("corrupt cache ignored")
	}
	after, _ := os.ReadFile(file)
	if string(after) != string(data) {
		t.Fatal("corrupt cache overwritten")
	}
}

func TestDirectoryCacheRestartAdoptionAndAmbiguity(t *testing.T) {
	cache := &directoryCacheStore{Agents: make(map[string]*cachedDirectoryAgent)}
	cache.put("old", "192.0.2.1", "dir", "", &DirectoryListing{Path: `C:\Lab`, Entries: []DirectoryEntry{}})
	if _, ok := cache.lookup("new", "192.0.2.1", "dir", "", true, map[string]bool{"new": true}); !ok {
		t.Fatal("unique persisted host cache not restored")
	}
	cache.put("second", "192.0.2.1", "dir", "", &DirectoryListing{Path: `C:\Other`, Entries: []DirectoryEntry{}})
	if _, ok := cache.lookup("third", "192.0.2.1", "dir", "", true, nil); ok {
		t.Fatal("ambiguous cached hosts merged")
	}
	if _, ok := cache.lookup("third", "192.0.2.1", "dir", "", false, nil); ok {
		t.Fatal("IP matching used despite active-peer ambiguity")
	}
}

func TestDirectorySnapshotCacheHitRefreshAndOffline(t *testing.T) {
	previous := directorySnapshots
	directorySnapshots = &directoryCacheStore{Agents: make(map[string]*cachedDirectoryAgent)}
	defer func() { directorySnapshots = previous }()
	client, peer := directoryClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	request := func(refresh bool, output string) *DirectoryListing {
		done := make(chan directoryResult, 1)
		go func() {
			listing, err := BrowseDirectorySnapshot(ctx, client.IPAddress, `C:\Lab`, "dir", refresh)
			done <- directoryResult{Listing: listing, Err: err}
		}()
		command, err := readBinaryMessage(peer)
		if err != nil {
			t.Fatal(err)
		}
		processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: command.Header.MessageID}, Data: []byte("Directory of C:\\Lab\n11/06/2026 09:30 PM 5 " + output + "\n")})
		result := <-done
		if result.Err != nil {
			t.Fatal(result.Err)
		}
		return result.Listing
	}
	first := request(false, "first.txt")
	if first.Cached {
		t.Fatal("fresh response marked cached")
	}
	// No peer read occurs here: a network request would block and fail on context expiry.
	hit, err := BrowseDirectorySnapshot(ctx, client.IPAddress, `c:/lab/`, "dir", false)
	if err != nil || !hit.Cached || !hit.CapturedAt.Equal(first.CapturedAt) {
		t.Fatalf("cache miss: %v", err)
	}
	refreshed := request(true, "new.txt")
	if refreshed.Cached || refreshed.Entries[0].Name != "new.txt" {
		t.Fatal("explicit refresh did not replace snapshot")
	}
	DisconnectSession(client.ID)
	offline, err := BrowseDirectorySnapshot(ctx, client.IPAddress, `C:\Lab`, "dir", false)
	if err != nil || !offline.Cached || offline.Entries[0].Name != "new.txt" {
		t.Fatalf("offline snapshot unavailable: %v", err)
	}
	if _, err := BrowseDirectorySnapshot(ctx, client.IPAddress, `C:\Missing`, "dir", false); err == nil {
		t.Fatal("offline cache miss accepted")
	}
	if _, err := BrowseDirectorySnapshot(ctx, client.IPAddress, `C:\Lab`, "dir", true); err == nil {
		t.Fatal("offline refresh accepted")
	}
}

func TestDirectoryCacheMergeAndPersistenceFailure(t *testing.T) {
	store := &directoryCacheStore{Agents: make(map[string]*cachedDirectoryAgent)}
	store.put("saved", "192.0.2.1", "dir", `C:\Lab`, &DirectoryListing{Path: `C:\Lab`, Entries: []DirectoryEntry{{Name: "old.txt"}}})
	store.put("live", "192.0.2.2", "dir", `C:\Other`, &DirectoryListing{Path: `C:\Other`, Entries: []DirectoryEntry{{Name: "live.txt"}}})
	if err := store.merge("saved", "live", "192.0.2.2"); err != nil {
		t.Fatal(err)
	}
	for _, directory := range []string{`C:\Lab`, `C:\Other`} {
		if _, ok := store.lookup("live", "192.0.2.2", "dir", directory, false, nil); !ok {
			t.Fatalf("linked cache missing %s", directory)
		}
	}
	if _, ok := store.Agents["saved"]; ok {
		t.Fatal("consumed cache record retained")
	}
	collision := filepath.Join(t.TempDir(), "directory-not-file")
	if err := os.Mkdir(collision, 0700); err != nil {
		t.Fatal(err)
	}
	store.file = collision
	if err := store.put("live", "192.0.2.2", "dir", `C:\Lab`, &DirectoryListing{Path: `C:\Lab`, Entries: []DirectoryEntry{{Name: "memory.txt"}}}); err == nil {
		t.Fatal("persistence failure ignored")
	}
	snapshot, ok := store.lookup("live", "192.0.2.2", "dir", `C:\Lab`, false, nil)
	if !ok || snapshot.Entries[0].Name != "memory.txt" || snapshot.CacheWarning == "" {
		t.Fatal("failed disk write lost memory snapshot or warning")
	}
}
