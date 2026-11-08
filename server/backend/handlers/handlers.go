package handlers

import (
	"encoding/json"
	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/utils"
	"net/http"
)

// RegisterHandlers connects frontend API routes to their backend handlers.
func RegisterHandlers() {

	// Placeholder for SCREENSHOT: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/screenshot", screenshotPlaceholderHandler)

	// Placeholder for CD: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/cd", cdPlaceholderHandler)

	// Placeholder for CP: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/cp", cpPlaceholderHandler)

	// Placeholder for DRIVES: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/drives", drivesPlaceholderHandler)

	// Placeholder for EXIT: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/exit", exitPlaceholderHandler)

	// Placeholder for GETPRIVS: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/getprivs", getprivsPlaceholderHandler)

	// Placeholder for GETUID: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/getuid", getuidPlaceholderHandler)

	// Placeholder for KERBEROS_CCACHE_USE: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/kerberos-ccache-use", kerberosCcacheUsePlaceholderHandler)

	// Placeholder for KERBEROS_TICKET_PURGE: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/kerberos-ticket-purge", kerberosTicketPurgePlaceholderHandler)

	// Placeholder for KERBEROS_TICKET_USE: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/kerberos-ticket-use", kerberosTicketUsePlaceholderHandler)

	// Placeholder for JOBKILL: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/jobkill", jobkillPlaceholderHandler)

	// Placeholder for KILL: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/kill", killPlaceholderHandler)

	// Placeholder for LINK: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/link", linkPlaceholderHandler)

	// Placeholder for LS: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/ls", lsPlaceholderHandler)

	// Placeholder for MAKE_TOKEN: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/make-token", makeTokenPlaceholderHandler)

	// Placeholder for MKDIR: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/mkdir", mkdirPlaceholderHandler)

	// Placeholder for MV: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/mv", mvPlaceholderHandler)

	// Placeholder for PPID: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/ppid", ppidPlaceholderHandler)

	// Placeholder for PS: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/ps", psPlaceholderHandler)

	// Placeholder for PWD: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/pwd", pwdPlaceholderHandler)

	// Placeholder for REG_QUERY: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/reg-query", regQueryPlaceholderHandler)

	// Placeholder for REG_QUERYV: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/reg-queryv", regQueryvPlaceholderHandler)

	// Placeholder for REV2SELF: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/rev2self", rev2selfPlaceholderHandler)

	// Placeholder for RM: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/rm", rmPlaceholderHandler)

	// Placeholder for RPORTFWD: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/rportfwd", rportfwdPlaceholderHandler)

	// Placeholder for SETENV: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/setenv", setenvPlaceholderHandler)

	// Placeholder for SOCKS: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/socks", socksPlaceholderHandler)

	// Placeholder for STEAL_TOKEN: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/steal-token", stealTokenPlaceholderHandler)

	// Placeholder for TIMESTOMP: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/timestomp", timestompPlaceholderHandler)

	// Placeholder for UNLINK: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/unlink", unlinkPlaceholderHandler)

	// List the non-executable attack and script placeholders available in the tools catalog.
	http.HandleFunc("/tool-catalog", toolCatalogHandler)

	// Record a simulated tool task and test result without contacting a client.
	http.HandleFunc("/simulate-tool", simulateToolHandler)

	// Retrieve saved task summaries or the details and output of one task.
	http.HandleFunc("/task-history", taskHistoryHandler)

	// Retrieve recent server audit events for the Event Viewer and audit history.
	http.HandleFunc("/audit-history", auditHistoryHandler)

	// Read cached directory listings or request a new listing from the selected client.
	http.HandleFunc("/browse-files", browseFilesHandler)

	// Request that the selected client upload a file to the server; acknowledge dispatch.
	http.HandleFunc("/download-file", downloadFilePlaceholderHandler)

	// Explicitly associate a live connection with a previously saved session.
	http.HandleFunc("/link-session", linkSessionHandler)

	// Stop accepting new TCP connections while preserving existing sessions.
	http.HandleFunc("/stop-listener", stopListenerHandler)

	// Report TCP listener status, bind address, and bound port.
	http.HandleFunc("/listener-status", listenerStatusHandler)

	// Close a selected session's connection while retaining its saved information.
	http.HandleFunc("/disconnect-session", disconnectSessionHandler)

	// Start the default TCP listener and report any bind failure.
	http.HandleFunc("/start-listener", startListenerHandler)

	// List snapshots of active and saved offline client sessions.
	http.HandleFunc("/get-clients", getClientsHandler)

	// Retrieve the selected client's recent incoming messages and session metadata.
	http.HandleFunc("/get-client-messages", getClientMessagesHandler)

	// Dispatch a shell command through the client's EXECUTE_COMMAND message handler.
	http.HandleFunc("/send-command", sendCommandHandler)

	// Request native host information through the existing SYS_INFO message handler.
	http.HandleFunc("/my-send-command", mySendCommandHandler)

	// Accept a file from the browser and send it to the selected client.
	http.HandleFunc("/send-binary", sendBinaryHandler)

	// Accept shellcode from the browser and dispatch it to the client's shellcode handler.
	http.HandleFunc("/send-shellcode", sendShellcodeHandler)

	// Retrieve the selected client's collected system information.
	http.HandleFunc("/get-client-info", getClientInfoHandler)
}

func getClientsHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)
	clients := network.GetAllClients()
	json.NewEncoder(w).Encode(clients)
}
