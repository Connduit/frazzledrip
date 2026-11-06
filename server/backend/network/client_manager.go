package network

import (
	"fmt"
	"github.com/Connduit/frazzledrip/server/backend/models"
	"net"
	"sync"
)

// ==============================
//   Connected Client Registry
// ==============================
//
// This file maintains a thread-safe registry of all connected clients.
// Clients are stored in a map keyed by their IP address (string → *Client).
// The mutex ensures concurrent read/write access is safe.

var (
	// clients holds all active clients mapped by their address.
	clients = make(map[string]*models.Client)

	// clientsMutex protects the clients map from concurrent access.
	clientsMutex sync.Mutex
)

// AddClient registers a new client in the global client map.
// A unique offline session from the same peer IP resumes under its original ID.
func AddClient(addr string, client *models.Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	host, _, err := net.SplitHostPort(addr)
	if err == nil && net.ParseIP(host) != nil && client.IsActive {
		var saved *models.Client
		matches := 0
		for oldAddr, candidate := range clients {
			oldHost, _, splitErr := net.SplitHostPort(oldAddr)
			if splitErr == nil && net.ParseIP(oldHost).Equal(net.ParseIP(host)) && !candidate.IsActive {
				saved = candidate
				matches++
			}
		}
		if matches == 1 {
			restoreSessionData(saved, client)
			client.ID = saved.ID
			delete(clients, saved.IPAddress)
			delete(pendingCommands, saved.IPAddress)
		}
	}
	clients[addr] = client
}

// RemoveClient deletes a client entry from the map using its address.
// If the address does not exist, nothing happens.
func RemoveClient(addr string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	delete(basicInfoQueues, addr)
	delete(clients, addr)
	delete(pendingCommands, addr)
}

// DisconnectClient keeps collected facts and history available for inspection.
// An older socket must not mark a replacement connection offline.
func DisconnectClient(client *models.Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if clients[client.IPAddress] != client {
		return
	}
	interruptBasicInfo(client)
	client.IsActive = false
	client.Conn = nil
	delete(pendingCommands, client.IPAddress)
}

// GetClient returns a pointer to the client with the given address.
// The second return value indicates whether the client was found.
func GetClient(addr string) (*models.Client, bool) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	client, exists := clients[addr]
	if !exists {
		return nil, false
	}
	return snapshotClient(client), true
}

// GetAllClients returns a slice of all clients currently connected.
// This is a snapshot taken under lock; changes after return are not reflected.
func GetAllClients() []*models.Client {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clientList := make([]*models.Client, 0, len(clients))
	for _, client := range clients {
		clientList = append(clientList, snapshotClient(client))
	}
	return clientList
}

// Call with clientsMutex held so HTTP encoding never reads mutable session data.
func snapshotClient(client *models.Client) *models.Client {
	copy := *client
	copy.Messages = append([]models.ClientMessage(nil), client.Messages...)
	copy.CommandResults = append([]models.CommandResult(nil), client.CommandResults...)
	copy.SystemInfo.IPAddresses = append([]string(nil), client.SystemInfo.IPAddresses...)
	copy.SystemInfo.MACAddresses = append([]string(nil), client.SystemInfo.MACAddresses...)
	return &copy
}

// DisconnectSession closes one selected socket but keeps its saved record.
func DisconnectSession(id string) error {
	clientsMutex.Lock()
	var selected *models.Client
	for _, client := range clients {
		if client.ID == id {
			selected = client
			break
		}
	}
	if selected == nil {
		clientsMutex.Unlock()
		return fmt.Errorf("session no longer exists")
	}
	conn := selected.Conn
	interruptBasicInfo(selected)
	selected.IsActive = false
	selected.Conn = nil
	delete(pendingCommands, selected.IPAddress)
	clientsMutex.Unlock()
	if conn != nil {
		conn.Close()
	}
	return nil
}
