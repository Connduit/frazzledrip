package network

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
	"github.com/Connduit/frazzledrip/server/backend/utils"
)

type DirectoryEntry struct {
	ModifiedText string  `json:"modified_text,omitempty"`
	ModifiedSort string  `json:"modified_sort,omitempty"`
	ModifiedMS   *int64  `json:"modified_ms,omitempty"`
	Name         string  `json:"name"`
	Path         string  `json:"path"`
	IsDirectory  bool    `json:"is_directory"`
	IsLink       bool    `json:"is_link"`
	Size         *uint64 `json:"size,omitempty"`
}
type DirectoryListing struct {
	CapturedAt   time.Time        `json:"captured_at,omitempty"`
	Cached       bool             `json:"cached"`
	CacheWarning string           `json:"cache_warning,omitempty"`
	Path         string           `json:"path"`
	Parent       string           `json:"parent"`
	Truncated    bool             `json:"truncated"`
	Entries      []DirectoryEntry `json:"entries"`
}
type directoryResult struct {
	Listing *DirectoryListing
	Err     error
}
type directoryRequest struct {
	ViaDir bool
	ID     uint32
	Client *models.Client
	Result chan directoryResult
}

// Protected by clientsMutex. Requests are correlated by ID, never by response order.
var directoryRequests = make(map[string]*directoryRequest)

func BrowseDirectory(ctx context.Context, addr, path string) (*DirectoryListing, error) {
	return browseDirectory(ctx, addr, path, false)
}
func BrowseDirectoryViaDir(ctx context.Context, addr, path string) (*DirectoryListing, error) {
	return browseDirectory(ctx, addr, path, true)
}
func browseDirectory(ctx context.Context, addr, path string, viaDir bool) (*DirectoryListing, error) {
	if len(path) > 2048 || strings.ContainsRune(path, 0) {
		return nil, fmt.Errorf("invalid directory path")
	}
	payload := path
	messageType := uint32(models.MESSAGE_LIST_DIRECTORY)
	if viaDir {
		var err error
		payload, err = dirCommand(path)
		if err != nil {
			return nil, err
		}
		messageType = models.MESSAGE_EXECUTE_COMMAND
	}
	clientsMutex.Lock()
	client := clients[addr]
	if client == nil || !client.IsActive || client.Conn == nil {
		clientsMutex.Unlock()
		return nil, fmt.Errorf("session is offline")
	}
	if _, ok := directoryRequests[addr]; ok {
		clientsMutex.Unlock()
		return nil, fmt.Errorf("a directory request is already pending")
	}
	if _, ok := pendingCommands[addr]; ok || (client.BasicInfoStatus == "collecting" || (client.BasicInfoStatus == "" && client.Type == "TCP")) {
		clientsMutex.Unlock()
		return nil, fmt.Errorf("wait for the current command or initial host queries to finish")
	}
	request := &directoryRequest{ViaDir: viaDir, ID: utils.GenerateMessageID(), Client: client, Result: make(chan directoryResult, 1)}
	directoryRequests[addr] = request
	clientsMutex.Unlock()
	defer func() {
		clientsMutex.Lock()
		if directoryRequests[addr] == request {
			delete(directoryRequests, addr)
		}
		clientsMutex.Unlock()
	}()
	message := models.InternalMessage{Header: models.MessageHeader{MessageType: messageType, MessageID: request.ID, DataSize: uint32(len(payload))}, Data: []byte(payload)}
	if err := sendBinaryMessage(client, message); err != nil {
		return nil, err
	}
	select {
	case result := <-request.Result:
		return result.Listing, result.Err
	case <-ctx.Done():
		return nil, fmt.Errorf("directory request timed out or was cancelled: %w", ctx.Err())
	}
}

func finishDirectoryRequest(client *models.Client, message *models.InternalMessage) bool {
	clientsMutex.Lock()
	request := directoryRequests[client.IPAddress]
	if request == nil || request.Client != client {
		clientsMutex.Unlock()
		return false
	}
	expected := message.Header.MessageID == request.ID || (request.ViaDir && message.Header.MessageID == 105)
	validType := message.Header.MessageType == models.MESSAGE_ERROR_REPORT || (!request.ViaDir && message.Header.MessageType == models.MESSAGE_DIRECTORY_LISTING) || (request.ViaDir && message.Header.MessageType == models.MESSAGE_COMMAND_RESULT)
	if !expected || !validType {
		clientsMutex.Unlock()
		return false
	}
	delete(directoryRequests, client.IPAddress)
	clientsMutex.Unlock()
	result := directoryResult{}
	if message.Header.MessageType == models.MESSAGE_ERROR_REPORT {
		result.Err = fmt.Errorf("%s", string(message.Data))
	} else if request.ViaDir {
		result.Listing, result.Err = parseDirOutput(string(message.Data))
	} else {
		var listing DirectoryListing
		if err := json.Unmarshal(message.Data, &listing); err != nil || listing.Path == "" || len(listing.Entries) > 2000 {
			result.Err = fmt.Errorf("client returned an invalid directory listing")
		} else {
			if listing.Entries == nil {
				listing.Entries = []DirectoryEntry{}
			}
			result.Listing = &listing
		}
	}
	request.Result <- result
	return true
}

// Call with clientsMutex held. Wake pending HTTP requests on disconnect/removal.
func cancelDirectoryRequest(addr string) {
	if request := directoryRequests[addr]; request != nil {
		delete(directoryRequests, addr)
		request.Result <- directoryResult{Err: fmt.Errorf("session disconnected")}
	}
}
