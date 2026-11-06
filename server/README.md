# Server

The backend embeds its frontend assets in the executable with Go's `embed`
package. The React frontend uses npm and Vite. Its source is in
`frontend/src`; Vite bundles React, JavaScript, and CSS into
`frontend/dist`. No browser CDN or runtime Babel compiler is needed.

With Go 1.22 or newer and Node.js 20.19+ or 22.12+:

```sh
cd server
go generate ./frontend
go build -o /tmp/frazzledrip-server ./backend
/tmp/frazzledrip-server
```

The executable serves the frontend and API on port 8080 from any working
directory. Frontend API requests use the same origin as the page. Node.js is
needed only to build the frontend; the resulting executable needs no frontend
files or npm installation at runtime.

`go generate ./frontend` runs `npm ci` and `npm run build`. Run it before the
first Go build and after changing frontend sources. Generated `dist` files and
`node_modules` are ignored by Git; `package-lock.json` pins dependencies.

For frontend development, start the Go backend as above, then in another terminal:

```sh
cd server/frontend
npm ci
npm run dev
```

Vite serves the frontend on port 5173 with hot reload and proxies the existing
API routes to the backend on port 8080. Use `npm run build` to regenerate the
assets embedded by the next Go build. `npm run preview` serves only the built
frontend for inspection; use the Go server to exercise the integrated API.
