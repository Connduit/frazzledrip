package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLinkSessionValidation(t *testing.T) {
	for _, test := range []struct {
		method, body string
		status       int
	}{
		{"GET", "", http.StatusMethodNotAllowed},
		{"POST", `{}`, http.StatusBadRequest},
		{"POST", `{"saved_id":"gone","live_id":"missing"}`, http.StatusConflict},
		{"POST", `not json`, http.StatusBadRequest},
		{"POST", `{"saved_id":"gone","live_id":"missing","unexpected":true}`, http.StatusBadRequest},
		{"OPTIONS", "", http.StatusOK},
	} {
		request := httptest.NewRequest(test.method, "/link-session", strings.NewReader(test.body))
		response := httptest.NewRecorder()
		linkSessionHandler(response, request)
		if response.Code != test.status {
			t.Errorf("%s %s: got %d want %d", test.method, test.body, response.Code, test.status)
		}
	}
}
