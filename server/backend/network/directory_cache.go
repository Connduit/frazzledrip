package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type cachedDirectoryAgent struct {
	Host     string                      `json:"host"`
	Listings map[string]DirectoryListing `json:"listings"`
	Aliases  map[string]string           `json:"aliases"`
}
type directoryCacheStore struct {
	mu     sync.Mutex
	file   string
	Agents map[string]*cachedDirectoryAgent `json:"agents"`
}

var directorySnapshots = &directoryCacheStore{Agents: make(map[string]*cachedDirectoryAgent)}

func DefaultDirectoryCacheDir() string {
	root, err := os.UserCacheDir()
	if err != nil {
		root = os.TempDir()
	}
	return filepath.Join(root, "frazzledrip")
}

// Configure at startup, before serving requests. Corrupt files are not overwritten.
func ConfigureDirectoryCache(directory string) error {
	store, err := loadDirectoryCache(filepath.Join(directory, "directories.json"))
	if err != nil {
		return err
	}
	directorySnapshots = store
	return nil
}
func loadDirectoryCache(file string) (*directoryCacheStore, error) {
	store := &directoryCacheStore{file: file, Agents: make(map[string]*cachedDirectoryAgent)}
	data, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, store); err != nil {
		return nil, fmt.Errorf("invalid directory cache: %w", err)
	}
	if store.Agents == nil {
		store.Agents = make(map[string]*cachedDirectoryAgent)
	}
	for _, agent := range store.Agents {
		if agent == nil {
			return nil, fmt.Errorf("invalid directory cache agent")
		}
		if agent.Listings == nil {
			agent.Listings = make(map[string]DirectoryListing)
		}
		if agent.Aliases == nil {
			agent.Aliases = make(map[string]string)
		}
	}
	return store, nil
}
func cloneDirectoryListing(listing DirectoryListing) DirectoryListing {
	entries := make([]DirectoryEntry, len(listing.Entries))
	copy(entries, listing.Entries)
	for i := range entries {
		if entries[i].Size != nil {
			value := *entries[i].Size
			entries[i].Size = &value
		}
		if entries[i].ModifiedMS != nil {
			value := *entries[i].ModifiedMS
			entries[i].ModifiedMS = &value
		}
	}
	listing.Entries = entries
	return listing
}
func directoryKey(source, directory string) string {
	if directory == "" {
		return source + ":<current>"
	}
	normalized := strings.ToLower(strings.ReplaceAll(directory, `\`, "/"))
	driveRoot := len(normalized) >= 3 && normalized[1] == ':' && normalized[2] == '/'
	unc := strings.HasPrefix(normalized, "//")
	normalized = path.Clean(normalized)
	if driveRoot && len(normalized) == 2 {
		normalized += "/"
	}
	if unc && !strings.HasPrefix(normalized, "//") {
		normalized = "/" + normalized
	}
	if normalized == "." {
		return source + ":<current>"
	}
	return source + ":" + normalized
}
func (store *directoryCacheStore) saveLocked() error {
	if store.file == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(store.file), 0700); err != nil {
		return err
	}
	data, err := json.Marshal(store)
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(store.file), "directories-*.tmp")
	if err != nil {
		return err
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err = temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err = temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err = temporary.Close(); err != nil {
		return err
	}
	return os.Rename(name, store.file)
}
func (store *directoryCacheStore) lookup(id, host, source, directory string, allowIPRestore bool, known map[string]bool) (*DirectoryListing, bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	agent := store.Agents[id]
	// Across a server restart IDs change. Reuse only one unclaimed cache for
	// this peer IP, following the existing single-host reconnect policy.
	if agent == nil && allowIPRestore {
		candidate := ""
		count := 0
		for oldID, old := range store.Agents {
			if old.Host == host && !known[oldID] {
				candidate = oldID
				count++
			}
		}
		if count == 1 {
			agent = store.Agents[candidate]
			store.Agents[id] = agent
			delete(store.Agents, candidate)
		}
	}
	if agent == nil {
		return nil, false
	}
	key := directoryKey(source, directory)
	if canonical, ok := agent.Aliases[key]; ok {
		key = canonical
	}
	listing, ok := agent.Listings[key]
	if !ok {
		return nil, false
	}
	listing = cloneDirectoryListing(listing)
	listing.Cached = true
	return &listing, true
}
func (store *directoryCacheStore) put(id, host, source, requested string, listing *DirectoryListing) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	agent := store.Agents[id]
	if agent == nil {
		agent = &cachedDirectoryAgent{Host: host, Listings: make(map[string]DirectoryListing), Aliases: make(map[string]string)}
		store.Agents[id] = agent
	}
	copy := cloneDirectoryListing(*listing)
	copy.CapturedAt = time.Now().UTC()
	copy.Cached = false
	copy.CacheWarning = ""
	key := directoryKey(source, copy.Path)
	agent.Listings[key] = copy
	agent.Aliases[directoryKey(source, requested)] = key
	*listing = cloneDirectoryListing(copy)
	if err := store.saveLocked(); err != nil {
		copy.CacheWarning = "Saved in server memory, but could not persist to disk: " + err.Error()
		agent.Listings[key] = copy
		return err
	}
	return nil
}
func (store *directoryCacheStore) merge(savedID, liveID, host string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	saved := store.Agents[savedID]
	if saved == nil {
		return nil
	}
	live := store.Agents[liveID]
	if live == nil {
		store.Agents[liveID] = saved
		saved.Host = host
	} else {
		for key, value := range saved.Listings {
			if existing, ok := live.Listings[key]; !ok || value.CapturedAt.After(existing.CapturedAt) {
				live.Listings[key] = value
			}
		}
		for key, value := range saved.Aliases {
			if _, ok := live.Aliases[key]; !ok {
				live.Aliases[key] = value
			}
		}
	}
	delete(store.Agents, savedID)
	return store.saveLocked()
}

func BrowseDirectorySnapshot(ctx context.Context, addr, directory, source string, refresh bool) (*DirectoryListing, error) {
	if source == "" {
		source = "dir"
	}
	if source != "dir" && source != "native" {
		return nil, fmt.Errorf("source must be dir or native")
	}
	client, ok := GetClient(addr)
	if !ok {
		return nil, fmt.Errorf("session no longer exists")
	}
	host, _, _ := net.SplitHostPort(addr)
	if !refresh {
		if snapshot, err := FindCachedDirectory(addr, directory, source); err != nil {
			return nil, err
		} else if snapshot != nil {
			return snapshot, nil
		}
	}
	if !client.IsActive {
		if refresh {
			return nil, fmt.Errorf("session is offline; open the saved directory without refreshing")
		}
		return nil, fmt.Errorf("session is offline and no cached listing is available for this path; reconnect to refresh")
	}
	browse := BrowseDirectoryViaDir
	if source == "native" {
		browse = BrowseDirectory
	}
	listing, err := browse(ctx, addr, directory)
	if err != nil {
		return nil, err
	}
	if err := directorySnapshots.put(client.ID, host, source, directory, listing); err != nil {
		listing.CacheWarning = "Saved in server memory, but could not persist to disk: " + err.Error()
	}
	return listing, nil
}

// A cache-only lookup never sends a message to the client.
func FindCachedDirectory(addr, directory, source string) (*DirectoryListing, error) {
	if source == "" {
		source = "dir"
	}
	if source != "dir" && source != "native" {
		return nil, fmt.Errorf("source must be dir or native")
	}
	client, ok := GetClient(addr)
	if !ok {
		return nil, fmt.Errorf("session no longer exists")
	}
	host, _, _ := net.SplitHostPort(addr)
	known := make(map[string]bool)
	active := 0
	for _, other := range GetAllClients() {
		known[other.ID] = true
		peer, _, _ := net.SplitHostPort(other.IPAddress)
		if peer == host && other.IsActive {
			active++
		}
	}
	snapshot, _ := directorySnapshots.lookup(client.ID, host, source, directory, active == 1 && host != "", known)
	return snapshot, nil
}
