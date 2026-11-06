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
sessions disappear when removed from the backend registry.

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
marketing names. Metadata is kept for the lifetime of the connected session.

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
