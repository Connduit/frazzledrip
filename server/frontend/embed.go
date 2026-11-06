package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:generate npm ci
//go:generate npm run build

//go:embed dist
var assets embed.FS

// Handler serves the frontend bundled into the server binary.
func Handler() http.Handler {
	files, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(files))
}
