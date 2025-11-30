package network

import (
	"backend/models"
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
// If an entry already exists for this address, it will be overwritten.
func AddClient(addr string, client *models.Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clients[addr] = client
}

// RemoveClient deletes a client entry from the map using its address.
// If the address does not exist, nothing happens.
func RemoveClient(addr string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	delete(clients, addr)
}

// GetClient returns a pointer to the client with the given address.
// The second return value indicates whether the client was found.
func GetClient(addr string) (*models.Client, bool) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	client, exists := clients[addr]
	return client, exists
}

// GetAllClients returns a slice of all clients currently connected.
// This is a snapshot taken under lock; changes after return are not reflected.
func GetAllClients() []*models.Client {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clientList := make([]*models.Client, 0, len(clients))
	for _, client := range clients {
		clientList = append(clientList, client)
	}
	return clientList
}
