package network

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Connduit/frazzledrip/server/backend/models"
	"github.com/Connduit/frazzledrip/server/backend/utils"
)

type DirectoryEntry struct {
	ModifiedMS  *int64  `json:"modified_ms,omitempty"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	IsDirectory bool    `json:"is_directory"`
	IsLink      bool    `json:"is_link"`
	Size        *uint64 `json:"size,omitempty"`
}
type DirectoryListing struct {
	Path      string           `json:"path"`
	Parent    string           `json:"parent"`
	Truncated bool             `json:"truncated"`
	Entries   []DirectoryEntry `json:"entries"`
}
type directoryResult struct {
	Listing *DirectoryListing
	Err     error
}
type directoryRequest struct {
	ID     uint32
	Client *models.Client
	Result chan directoryResult
}

// Protected by clientsMutex. Requests are correlated by ID, never by response order.
var directoryRequests = make(map[string]*directoryRequest)

func BrowseDirectory(ctx context.Context, addr, path string) (*DirectoryListing, error) {
	if len(path) > 2048 || strings.ContainsRune(path, 0) {
		return nil, fmt.Errorf("invalid directory path")
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
	request := &directoryRequest{ID: utils.GenerateMessageID(), Client: client, Result: make(chan directoryResult, 1)}
	directoryRequests[addr] = request
	clientsMutex.Unlock()
	defer func() {
		clientsMutex.Lock()
		if directoryRequests[addr] == request {
			delete(directoryRequests, addr)
		}
		clientsMutex.Unlock()
	}()
	message := models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_LIST_DIRECTORY, MessageID: request.ID, DataSize: uint32(len(path))}, Data: []byte(path)}
	if err := sendBinaryMessage(client, message); err != nil {
		return nil, err
	}
	select {
	case result := <-request.Result:
		return result.Listing, result.Err
	case <-ctx.Done():
		return nil, fmt.Errorf("directory request timed out or was cancelled; the client must support directory browsing: %w", ctx.Err())
	}
}

func finishDirectoryRequest(client *models.Client, message *models.InternalMessage) {
	if message.Header.MessageType != models.MESSAGE_DIRECTORY_LISTING && message.Header.MessageType != models.MESSAGE_ERROR_REPORT {
		return
	}
	clientsMutex.Lock()
	request := directoryRequests[client.IPAddress]
	if request == nil || request.Client != client || message.Header.MessageID != request.ID {
		clientsMutex.Unlock()
		return
	}
	delete(directoryRequests, client.IPAddress)
	clientsMutex.Unlock()
	result := directoryResult{}
	if message.Header.MessageType == models.MESSAGE_ERROR_REPORT {
		result.Err = fmt.Errorf("%s", string(message.Data))
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
}

// Call with clientsMutex held. Wake pending HTTP requests on disconnect/removal.
func cancelDirectoryRequest(addr string) {
	if request := directoryRequests[addr]; request != nil {
		delete(directoryRequests, addr)
		request.Result <- directoryResult{Err: fmt.Errorf("session disconnected")}
	}
}
