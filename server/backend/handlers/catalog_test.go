package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCatalogRequestValidation(t *testing.T) {
	for _, test := range []struct {
		body   string
		status int
	}{
		{`{}`, 400},
		{`{"agent_id":"x","tool_id":"unknown"}`, 400},
		{`{"agent_id":"x","tool_id":"host-survey","command":"whoami"}`, 400},
		{`{"agent_id":"x","tool_id":"host-survey","script":"test"}`, 400},
		{`{"agent_id":"x","tool_id":"host-survey"} {}`, 400},
		{`{"agent_id":"missing-session","tool_id":"host-survey"}`, 404},
	} {
		response := httptest.NewRecorder()
		simulateToolHandler(response, httptest.NewRequest(http.MethodPost, "/simulate-tool", strings.NewReader(test.body)))
		if response.Code != test.status {
			t.Fatalf("%s: got %d want %d", test.body, response.Code, test.status)
		}
	}
	response := httptest.NewRecorder()
	toolCatalogHandler(response, httptest.NewRequest(http.MethodGet, "/tool-catalog", nil))
	if response.Code != 200 || !strings.Contains(response.Body.String(), "host-survey") {
		t.Fatal("catalog not available")
	}
	response = httptest.NewRecorder()
	simulateToolHandler(response, httptest.NewRequest(http.MethodGet, "/simulate-tool", nil))
	if response.Code != 405 {
		t.Fatal("GET changed state")
	}
}
