package network

import (
	"backend/models"
	"sync"
)

var (
	clients      = make(map[string]*models.Client)
	clientsMutex sync.Mutex
)

func AddClient(addr string, client *models.Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients[addr] = client
}

func RemoveClient(addr string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, addr)
}

func GetClient(addr string) (*models.Client, bool) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	client, exists := clients[addr]
	return client, exists
}

func GetAllClients() []*models.Client {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clientList := make([]*models.Client, 0, len(clients))
	for _, client := range clients {
		clientList = append(clientList, client)
	}
	return clientList
}
