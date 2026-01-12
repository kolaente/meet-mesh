package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/kolaente/meet-mesh/api/embedded"
)

func NewStaticHandler() http.Handler {
	// Get the dist subdirectory from the embedded FS
	distFS, err := fs.Sub(embedded.Frontend, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := r.URL.Path

		// Check if file exists
		f, err := distFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// For SPA: serve index.html (or 200.html) for non-existent paths
		// SvelteKit adapter-static creates 200.html for SPA fallback
		r.URL.Path = "/200.html"
		fileServer.ServeHTTP(w, r)
	})
}
