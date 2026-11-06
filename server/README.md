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
`http://127.0.0.1:8080`. If port 5173 is occupied, use the port printed by Vite.
Frontend edits appear through hot reload. After changing Go code, stop and rerun
`go run ./backend` from `server`; Go does not automatically reload.

## Final build: one executable

Build the production frontend first, then compile it into the Go server:

```sh
cd server
go generate ./frontend
go build -trimpath -o ../frazzledrip-server ./backend
cd ..
./frazzledrip-server
```

Open `http://localhost:8080`. This executable serves both the frontend and API
and can run from any working directory. Deploy only the executable; Node.js,
`node_modules`, and the frontend source/build directories are not needed at
runtime. Stop the development backend first if it is still using port 8080.

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
