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
		Client    string `json:"client"`
		Source    string `json:"source"`
		Path      string `json:"path"`
		Refresh   bool   `json:"refresh"`
		CacheOnly bool   `json:"cache_only"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192)).Decode(&request); err != nil || request.Client == "" {
		controlError(w, http.StatusBadRequest, "client and a valid path are required")
		return
	}
	if request.CacheOnly {
		listing, err := network.FindCachedDirectory(request.Client, request.Path, request.Source)
		if err != nil {
			controlError(w, http.StatusBadRequest, err.Error())
			return
		}
		if listing == nil {
			controlError(w, http.StatusNotFound, "No saved listing for this directory")
			return
		}
		json.NewEncoder(w).Encode(listing)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	listing, err := network.BrowseDirectorySnapshot(ctx, request.Client, request.Path, request.Source, request.Refresh)
	if err != nil {
		controlError(w, http.StatusBadGateway, err.Error())
		return
	}
	json.NewEncoder(w).Encode(listing)
}
