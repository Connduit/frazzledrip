package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/network"
)

func browseFilesHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodPost) {
		return
	}
	var request struct {
		Client string `json:"client"`
		Path   string `json:"path"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192)).Decode(&request); err != nil || request.Client == "" {
		controlError(w, http.StatusBadRequest, "client and a valid path are required")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	listing, err := network.BrowseDirectory(ctx, request.Client, request.Path)
	if err != nil {
		controlError(w, http.StatusBadGateway, err.Error())
		return
	}
	json.NewEncoder(w).Encode(listing)
}
