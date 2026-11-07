package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/utils"
)

func requireControlMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	utils.EnableCORS(&w, r)
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		return false
	}
	if r.Method != method {
		w.Header().Set("Allow", method+", OPTIONS")
		controlError(w, http.StatusMethodNotAllowed, "method not allowed")
		return false
	}
	return true
}
func controlError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
func startListenerHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodPost) {
		return
	}
	if err := network.StartReverseTCPListener(4444); err != nil {
		controlError(w, http.StatusConflict, err.Error())
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "✅ Listener started on port 4444"})
}
func stopListenerHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodPost) {
		return
	}
	if err := network.StopReverseTCPListener(); err != nil {
		controlError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "Listener stopped. Existing sessions remain connected."})
}
func listenerStatusHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodGet) {
		return
	}
	running, port, bindAddress := network.ListenerDetails()
	json.NewEncoder(w).Encode(map[string]interface{}{"running": running, "port": port, "bind_address": bindAddress})
}
func disconnectSessionHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodPost) {
		return
	}
	var request struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&request); err != nil || request.ID == "" {
		controlError(w, http.StatusBadRequest, "session id is required")
		return
	}
	if err := network.DisconnectSession(request.ID); err != nil {
		controlError(w, http.StatusNotFound, err.Error())
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "Session disconnected; saved information retained. A reconnecting client may return."})
}
