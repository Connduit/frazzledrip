package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisteredCommandPlaceholders(t *testing.T) {
	original := http.DefaultServeMux
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	RegisterHandlers()
	http.DefaultServeMux = original
	commands := []struct{ path, command string }{
		{"/commands/screenshot", "SCREENSHOT"},
		{"/commands/cd", "CD"},
		{"/commands/cp", "CP"},
		{"/commands/drives", "DRIVES"},
		{"/commands/exit", "EXIT"},
		{"/commands/getprivs", "GETPRIVS"},
		{"/commands/getuid", "GETUID"},
		{"/commands/kerberos-ccache-use", "KERBEROS_CCACHE_USE"},
		{"/commands/kerberos-ticket-purge", "KERBEROS_TICKET_PURGE"},
		{"/commands/kerberos-ticket-use", "KERBEROS_TICKET_USE"},
		{"/commands/jobkill", "JOBKILL"},
		{"/commands/kill", "KILL"},
		{"/commands/link", "LINK"},
		{"/commands/ls", "LS"},
		{"/commands/make-token", "MAKE_TOKEN"},
		{"/commands/mkdir", "MKDIR"},
		{"/commands/mv", "MV"},
		{"/commands/ppid", "PPID"},
		{"/commands/ps", "PS"},
		{"/commands/pwd", "PWD"},
		{"/commands/reg-query", "REG_QUERY"},
		{"/commands/reg-queryv", "REG_QUERYV"},
		{"/commands/rev2self", "REV2SELF"},
		{"/commands/rm", "RM"},
		{"/commands/rportfwd", "RPORTFWD"},
		{"/commands/setenv", "SETENV"},
		{"/commands/socks", "SOCKS"},
		{"/commands/steal-token", "STEAL_TOKEN"},
		{"/commands/timestomp", "TIMESTOMP"},
		{"/commands/unlink", "UNLINK"},
	}
	for _, command := range commands {
		t.Run(command.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, command.path, strings.NewReader(`{"agent_id":"test","payload":"ignored"}`)))
			if response.Code != http.StatusNotImplemented {
				t.Fatalf("status %d, want 501", response.Code)
			}
			var body map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body["error"] != command.command+" is not implemented. No request was sent to the client." {
				t.Fatalf("unexpected response: %s", body["error"])
			}
			response = httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, command.path, nil))
			if response.Code != http.StatusMethodNotAllowed {
				t.Fatalf("GET status %d, want 405", response.Code)
			}
			response = httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(http.MethodOptions, command.path, nil))
			if response.Code != http.StatusOK {
				t.Fatalf("OPTIONS status %d, want 200", response.Code)
			}
		})
	}
}
