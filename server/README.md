# Server

The backend embeds its frontend assets in the executable with Go's `embed`
package. The React frontend uses npm and Vite. Its source is in
`frontend/src`; Vite bundles React, JavaScript, and CSS into
`frontend/dist`. No browser CDN or runtime Babel compiler is needed.

## Prerequisites

- Go 1.22 or newer.
- Node.js 20.19+ or 22.12+ with npm.

All commands below start from the repository root unless stated otherwise.

## Development

Run the backend and Vite in separate terminals. Vite provides frontend hot
reload and forwards API requests to the Go backend.

**Terminal 1 — prepare assets and start the backend:**

```sh
cd server
go generate ./frontend
go run ./backend
```

`go generate ./frontend` runs `npm ci` and `npm run build`. This initial build is
required even during development because Go's `embed` directive needs an existing
`frontend/dist` directory at compile time. Repeat it after a fresh checkout or
when you want to refresh the frontend served directly by Go.

**Terminal 2 — start the frontend development server:**

```sh
cd server/frontend
npm run dev
```

Open `http://localhost:5173` for development. Vite proxies API requests to
`http://127.0.0.1:8081`. If port 5173 is occupied, use the port printed by Vite.
Frontend edits appear through hot reload. After changing Go code, stop and rerun
`go run ./backend` from `server`; Go does not automatically reload.

If port 8081 is already in use, choose another backend address:

```sh
# Terminal 1, from server
go run ./backend -addr 127.0.0.1:8082

# Terminal 2, from server/frontend (Bash)
BACKEND_URL=http://127.0.0.1:8082 npm run dev
```

In PowerShell, set `$env:BACKEND_URL = 'http://127.0.0.1:8082'` before
`npm run dev`. Restart Vite after changing its backend address.

## Final build: one executable

Build the production frontend first, then compile it into the Go server:

```sh
cd server
go generate ./frontend
go build -trimpath -o ../frazzledrip-server ./backend
cd ..
./frazzledrip-server
```

Open `http://localhost:8081`. This executable serves both the frontend and API
and can run from any working directory. Deploy only the executable; Node.js,
`node_modules`, and the frontend source/build directories are not needed at
runtime. Stop the development backend first if it is still using port 8081.

To use another port for the combined executable, run
`./frazzledrip-server -addr 127.0.0.1:8082` (or
`./frazzledrip-server.exe -addr 127.0.0.1:8082` on Windows), then open
`http://localhost:8082`. The embedded frontend uses the same origin automatically.

For a Windows executable, run this build command from `server` after generating
the frontend assets (Bash syntax shown):

```sh
GOOS=windows GOARCH=amd64 go build -trimpath -o ../frazzledrip-server.exe ./backend
```

On Windows, launch it with `./frazzledrip-server.exe` in PowerShell. The default
build targets the build machine's operating system and architecture; the command
above targets Windows x64.

Rebuild the frontend and Go executable whenever frontend changes need to be
included in the final binary. `go build` does not run `go generate` automatically.
Generated `frontend/dist` and `frontend/node_modules` are ignored by Git;
`package-lock.json` pins npm dependencies.

`npm run preview` (from `server/frontend`) serves only the built frontend for
inspection. Use the Go executable for the combined frontend/API workflow.

## Session dashboard

Search sessions by host, user, or address and sort by host or last seen. Select
**Inspect**, a graph node, or a map marker to open the same session inspector.
The graph shows direct controller connections; it does not infer pivots. The
inventory refreshes every five seconds and reports failed refreshes. Disconnected
sessions remain visible as offline, with metadata and command history preserved
in the persistent server database. Offline sessions cannot receive commands.

The map bundles Natural Earth country boundaries for offline use. The server has
no IP geolocation data. Select a session and enter a label, latitude, and longitude
to place it on the map. Locations are manually assigned and saved in this
browser's local storage, keyed by host address without the remote port. Sessions
from the same address share the location; other browsers do not. Use **Clear**
to remove it. Addresses alone are not treated as verified physical locations.
No session addresses are sent to external mapping or geolocation services.

### System information

System-info responses update the session's OS version/build, architecture,
logical processor count, and physical/available memory. Both `SYS_INFO` responses
and the current client's version-keyed `COMMAND_RESULT` format are supported.
The original response stays in the message history. The table picks up changes
on its next refresh; the selected inspector refreshes every two seconds.
Missing or invalid fields preserve existing values. Hostname, user, and product
type remain unknown when the client does not report them. Memory is shown in GiB;
Windows versions are shown as reported version/build numbers rather than guessed
marketing names. Metadata remains available when the connection drops, in the persistent server database.

### Metadata from command responses

Shell commands sent through the controller are associated with their responses:

- `whoami`: current user identity.
- `ipconfig` / `ipconfig /all`: host name, DNS suffix, reported IPv4/IPv6 addresses,
  and MAC addresses when present.
- `systeminfo`: host name, OS name/version/build, architecture, manufacturer,
  model, domain, boot time, processor count, and physical memory when present.
- Other commands, including `dir C:\`, `netstat -an`, `tasklist`, and `net user`:
  the latest output is retained under **Collected command output** in the session
  inspector. The latest 20 distinct commands are retained per connection.

Metadata parsing supports the shown English Windows output formats. Empty,
failed, localized, or unrecognized fields leave existing facts intact; original
responses remain available. Reported interface IPs are distinct from the TCP
remote endpoint, and do not change map coordinates automatically.

The current Windows client returns the fixed response ID `105` for shell commands.
The server therefore allows one shell command in flight per session, rejecting a
second until the first responds. It also accepts clients that echo request IDs.
No client rebuild is required. If a command never responds, its slot clears when
the session disconnects; reconnect before sending another shell command.

### Disconnects

A lost connection marks its record offline instead of deleting collected facts,
messages, and command outputs. Offline graph links are dashed and their inspector
remains available. Pending commands are cleared when the connection ends.

A reconnect from the same observed peer IP reuses its unique offline session,
including the original agent ID, facts, messages, and command outputs. Remote
source-port changes do not create a duplicate row. The inspector follows the
agent ID onto the new endpoint. Active connections are kept separate; if several
offline records share the IP, use explicit linking rather than guessing. IP
matching is the chosen heuristic, not a verified device identity: NAT or WSL
forwarding may make different devices appear under one address. Session records are persisted and reload offline after a server restart. The existing 60-second inbound read
timeout still applies; clients that send no heartbeats can hit that timeout.

### Link a reconnect to saved information

Select the new **active** connection in the session table. In its inspector,
use **Link a saved session**, choose the corresponding offline record, and click
**Link saved session**. This explicit action transfers the saved facts, message
history (latest 100), and latest command outputs (20 distinct commands) into the
live connection, removing the consumed offline row. Existing live facts win over
saved values; duplicate command outputs keep the newest response. The active
connection retains its ID, endpoint, and pending command. Verify which device
the saved record belongs to before linking. Use this flow for ambiguous IPs
or when a device reconnects from a different address.

If the saved host has a manually assigned map location in this browser, it is
copied to the new host address only when the new address has no saved location.
No Windows client changes or rebuild are required. Linking and saved records
are stored in the persistent server database.

### Session table columns

The table shows Agent ID, Type, External, Listener, Internal, Domain, Computer,
User, OS, Process, PID, Tags, Last, and Sleep, in that order. Click the agent ID
to inspect a session; the adjacent dot indicates active/offline status.
External is the remote peer IP observed by the server, which may be a private,
loopback, or NAT address. Listener is the local endpoint that accepted the
connection; Type is its transport. Internal addresses and host facts come from
collected command responses. Last is the last message timestamp. Process, PID,
Tags, and Sleep currently remain unreported (`—`); tasklist output alone does
not identify the agent process, and UI polling intervals are not agent sleep
settings. Memory and other detailed hardware fields remain in the inspector.

### Connection controls and graph zoom

Use **Start listener** and **Stop listener** in the header. The status badge
reflects the actual listener state. Stopping closes only the listening socket:
existing sessions remain connected. New clients cannot connect until it starts
again. Bind errors are shown instead of reporting a successful start.

Select an active session and click **Disconnect session** to close its socket,
clear pending commands, and keep its metadata/history offline. This does not
stop the client program: a client that retries may reconnect, at which point
the same-IP restoration rule applies. Stop the listener first if you also want
to prevent reconnections through that listener.

Use the mouse wheel or **+** / **−** above the connection graph to zoom from
50% to 300%. Wheel zoom centers on the pointer; left-click and drag to pan. **Fit graph** resets to 100% and centers the graph in the viewport. Zoom is retained during session refreshes; nodes remain selectable.

### Automatic basic host queries

A new session automatically receives `whoami`, `ipconfig /all`, and `systeminfo`,
one at a time. Responses update the session facts and appear in its command
outputs and message history. Manual shell commands wait until this initial
sequence finishes. The inspector reports collection progress, completion,
interruption, or send failure. Completion means each query returned, not that
all fields were available; parsing still supports English Windows output.
Recognized reconnects preserve this state and do not replay the initial queries.
If collection is interrupted, collect missing information manually after reconnect.
No Windows client rebuild is required.

The session table uses wrapped, fixed-width columns to fit the workspace without
horizontal scrolling. Multiple internal IPv4 addresses each occupy their own
line. On narrow screens, sessions use labeled fields in a two-column layout.

### Workspace console

Select a session in the graph or table. The **Console** accepts an explicit list
of custom message commands; it does not accept shell or terminal commands.

- `SYS_INFO` requests native OS, memory, and processor information using the
  existing `SYS_INFO` message handler. Names are case insensitive.
- `help` displays supported console commands locally.
- `clear`, **Clear history**, or Ctrl+L clears the console view and recalled
  commands for the selected session. Saved server history is retained.

Press Enter or **Run** to submit; Up/Down recalls typed commands. `help` and
`clear` also work for offline sessions. Unknown commands, shell syntax, and
response-only message types are rejected locally without a network request.
The **Commands** panel also provides only the native `SYS_INFO` operation.

Native responses appear in the console, including the current client's legacy
`COMMAND_RESULT` responses with message ID 99. Shell-command output does not
appear here. Transfer status and directory errors remain in the Event Viewer.
This UI change requires no client rebuild or new backend endpoint. The existing
shell backend endpoint remains available to other existing server features;
removing shell entry from the console does not disable it globally.

### File Explorer

Select a session and open the **Files** tab. Leave the path empty and
click **Open directory** for the client's current directory, or enter an absolute
Windows path. Click a folder to navigate, use **Up** for its parent, and
**Refresh directory** to reload. Filter names locally and select a file to view
its full path and size. Directory links are labeled; folders sort first.

This is a read-only, nonrecursive listing of names and metadata. It does not
read file contents, download, delete, rename, or upload files. The existing File
Transfer tab remains separate. Listings display at most 2,000 entries and are clearly marked when truncated.
Native listings also cap metadata at about 8 MiB. Permission errors, invalid paths, disconnects, and the
20-second response timeout are surfaced without clearing the previous listing.
Wait for initial host collection or any current shell command before browsing.

The default **Windows dir** method uses the existing Execute Command message,
so it requires only rebuilding the server/frontend. The server sends a quoted,
nonrecursive `dir /a /-c /4 /t:w` request with formatting overrides, then parses
names, sizes, folders, links, and modification dates from its response. English
`Directory of` headers are required. Names remain read-only; invalid paths and
unrecognized formats produce an error instead of an empty directory. CMD paths
containing percent signs, exclamation marks, wildcards, quotes, or control
characters are rejected. Use **Native listing** for such names or unsupported
locales; it requires the updated Windows client and protocol messages 13/14.
The native client Windows build/runtime validation remains outstanding here.

The File Explorer detail columns are **Name**, **Size**, **Date modified**, and
**Type**. Click a header to toggle ascending/descending sorting, with folders
kept first and missing values last. Types describe the file extension or folder,
including link indicators. Sizes use B/KiB/MiB/GiB/TiB. For dir results, Date modified preserves the date/time text reported by the
client; no timezone is invented. Sorting understands month-first slash dates,
year-first ISO dates, and day-first dotted dates with supported times. Other
date formats remain visible as text but have no parsed sort key. Native results
continue to use modification timestamps displayed in the browser's local time.
The properties panel shows the same value.

### Workspace layout

The desktop workspace shows the graph and session table above one shared window
with **Console**, **Event Viewer**, and **Files** tabs. Selecting a session updates
the Console and Files tabs. Switching tabs preserves console drafts, file paths,
and event filters; the event feed continues refreshing while its tab is hidden
unless paused. Arrow keys, Home, and End navigate the tabs.

Drag the vertical divider to change graph/table widths, or the horizontal divider
to resize the shared tabbed window. Keyboard arrow keys also resize focused
dividers. Sizes are saved in this browser; **Reset panel sizes** restores defaults.
Narrow screens stack the graph, sessions, and tabbed window.

Listener controls and disconnect remain in the toolbar/context strip.
**Details & tools** opens host properties, saved-session linking, collected output,
and existing controls. **Map** opens session geography separately.

### Saved directory snapshots

Directory listings are cached on the server per agent, directory, and listing
method. Open, folder navigation, and Up use an existing snapshot; only a first
visit to an uncached path queries the client. **Refresh directory** explicitly
queries the client and replaces the snapshot after a successful response.
Failures keep the previous snapshot. Listings show their capture time and whether
they came from cache. Changing a client directory does not update its saved
listing until you refresh it.

Selecting a session first makes a cache-only request to restore its saved current
directory view; this never sends a client command. Cached paths can be opened
while offline, but uncached paths and refresh require a live session. Captures
survive reconnects and explicit saved-session linking.

Snapshots are written to `directories.json` in the OS user cache directory's
`frazzledrip` subfolder (for example, `$XDG_CACHE_HOME/frazzledrip` or
`~/.cache/frazzledrip` on Linux, and `%LocalAppData%\frazzledrip` on Windows).
Use `-cache-dir /path/to/folder` to choose another location. Writes replace the
file through a temporary file, with private file permissions where supported.
If persistence fails, the snapshot remains usable in memory and the UI reports
that it could not be saved to disk. An invalid existing cache file is reported at
startup rather than overwritten.

After a server restart, once a client reconnects, one unclaimed cached agent
matching its observed IP can be restored under the existing IP-based identity
policy. Ambiguous matches are not assigned automatically. Session records now persist separately in the server database described below. Cache size grows with visited directories; there is no automatic expiry.
No client rebuild is required for caching.

## Persistent server foundation (milestone 1)

The server now stores sessions, host facts, the existing bounded message/output
history, shell-task outcomes, and audit events in a transactional bbolt database.
This pure-Go store keeps Linux and Windows builds independent of CGO. The default
file is `server.db` alongside the directory cache. Use `-database /path/server.db`
to choose another file. The database is private to one server process; stop the
server before copying its database as a backup. Do not share one database between
concurrent server processes.

At startup, saved sessions appear offline. Reconnecting follows the existing
IP-based matching policy and retains the saved agent ID and collected facts.
Unfinished tasks become **interrupted**, with an audit event; they are never
replayed. This remains a legacy IP heuristic, not authenticated device identity.
The HTTP interface now defaults to `127.0.0.1:8081`.

Under **Details & tools**, **Shell task history** shows the latest 500 manual
shell and initial-host-info tasks for the selected agent. Click a command to
inspect its archived output (limited to 1 MiB, explicitly marked if truncated).
States are queued (recorded before socket dispatch), sent (written, not
acknowledged), completed, failed, and interrupted. These states record existing
execution; this milestone does not add a background job scheduler or automatic
retry. Existing clients with fixed response ID 105 still use one pending shell
command; test fixtures can echo the request ID.

**Audit** shows the latest 300 durable lifecycle events and POST request outcomes.
HTTP peer endpoints and task origins are recorded; request bodies and query
strings are excluded. Peer addresses are not authenticated operator identities.
Directory/native operations currently receive HTTP audit entries; shell commands
and initial host queries additionally have task lifecycle records.

This is the first server-side foundation, not a completed security rollout.
Operator login, authenticated stable agent identity, TLS provisioning, roles,
advanced task scheduling/cancellation, and engagement isolation remain follow-up
work. Keep the HTTP service on its loopback default while those are developed.
Existing APIs do not yet authenticate operators. The database audit is not a
cryptographically tamper-evident log.

The Windows client is unchanged. `client/test/fixtures/mock_agent.py` provides
TEST-only responses without executing commands or reading files. See its adjacent
README for an isolated test workflow. Rebuild the server/frontend for this phase;
no production client changes or rebuild are required.

## Details & tools and placeholder catalog

**Details & tools** opens a modal session workspace with a session picker and
Overview, Task history, Attacks & scripts, Commands & transfers, and Messages
sections. Close it with its Close button or Escape. The modal keeps keyboard
focus inside and returns focus to its opener. Overview includes saved-session
linking; offline sessions retain their inspectable history.

The Attacks & scripts catalog is development scaffolding only. Discovery,
privilege assessment, credential-access assessment, lateral-movement exercise,
persistence exercise, defense-validation, and script/report-template entries
have no executable implementation. Search or filter by category/kind, optionally
add a note, and choose **Simulate task**. The server writes a completed simulation
task and its audit events, with a synthetic TEST result. There is no client
contact, shell command, script execution, filesystem operation, or host change.
Simulation tasks have `simulation: true`, origin `simulation`, and a
`[SIMULATION]` command label; the history UI also shows a Simulation badge.

`GET /tool-catalog` returns the fixed catalog. `POST /simulate-tool` accepts only
`agent_id`, a whitelisted `tool_id`, and an optional note up to 500 characters.
Unknown fields (including commands and scripts) are rejected. The target session
must exist, but may be offline because simulation runs only on the server.
Existing command/transfer controls retain their previous behavior and are
separate from these placeholders. No production client changes are required.
