package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Connduit/frazzledrip/server/backend/catalog"
	"github.com/Connduit/frazzledrip/server/backend/network"
)

func toolCatalogHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodGet) {
		return
	}
	json.NewEncoder(w).Encode(catalog.List())
}
func simulateToolHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodPost) {
		return
	}
	var request struct {
		AgentID string `json:"agent_id"`
		ToolID  string `json:"tool_id"`
		Note    string `json:"note"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil || request.AgentID == "" {
		controlError(w, http.StatusBadRequest, "agent_id and tool_id are required; only a note is supported")
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		controlError(w, http.StatusBadRequest, "expected one JSON object")
		return
	}
	if err := catalog.Validate(request.ToolID, request.Note); err != nil {
		controlError(w, http.StatusBadRequest, err.Error())
		return
	}
	exists := false
	for _, client := range network.GetAllClients() {
		if client.ID == request.AgentID {
			exists = true
			break
		}
	}
	if !exists {
		controlError(w, http.StatusNotFound, "session no longer exists")
		return
	}
	store := network.PersistentState()
	if store == nil {
		controlError(w, http.StatusServiceUnavailable, "persistent store unavailable")
		return
	}
	task, err := catalog.Simulate(store, request.AgentID, request.ToolID, request.Note)
	if err != nil {
		controlError(w, http.StatusInternalServerError, "could not record simulation")
		return
	}
	json.NewEncoder(w).Encode(task)
}
