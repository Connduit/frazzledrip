package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/utils"
)

func linkSessionHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}
	var request struct {
		SavedID string `json:"saved_id"`
		LiveID  string `json:"live_id"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil || request.SavedID == "" || request.LiveID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "saved_id and live_id are required"})
		return
	}
	if err := network.LinkSavedSession(request.SavedID, request.LiveID); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "Saved session information linked to the active connection."})
}
