package handlers

import "net/http"

// downloadFilePlaceholderHandler deliberately performs no client or file operations.
func downloadFilePlaceholderHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodPost) {
		return
	}
	controlError(w, http.StatusNotImplemented, "Download to server is not implemented. No request was sent to the client.")
}
