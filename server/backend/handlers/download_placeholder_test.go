package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDownloadPlaceholder(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			response := httptest.NewRecorder()
			downloadFilePlaceholderHandler(response, httptest.NewRequest(method, "/download-file", nil))
			want := http.StatusMethodNotAllowed
			if method == http.MethodPost {
				want = http.StatusNotImplemented
				if !strings.Contains(response.Body.String(), "No request was sent to the client") {
					t.Fatalf("missing placeholder explanation: %s", response.Body.String())
				}
			}
			if response.Code != want {
				t.Fatalf("status = %d, want %d", response.Code, want)
			}
		})
	}
}
