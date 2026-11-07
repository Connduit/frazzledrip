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
			downloadFilePlaceholderHandler(response, httptest.NewRequest(method, "/download-file", strings.NewReader(`{"agent_id":"test","path":"C:\\Test Folder\\example.txt"}`)))
			want := http.StatusMethodNotAllowed
			if method == http.MethodPost {
				want = http.StatusNotFound
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

func TestDownloadPlaceholderRejectsInvalidPath(t *testing.T) {
	for _, body := range []string{`{`, `{}`, `{"path":" "}`, `{"agent_id":"test","path":"test"} {}`, `{"agent_id":"test","path":"` + strings.Repeat("a", 4097) + `"}`} {
		response := httptest.NewRecorder()
		downloadFilePlaceholderHandler(response, httptest.NewRequest(http.MethodPost, "/download-file", strings.NewReader(body)))
		if response.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", response.Code)
		}
	}
}
