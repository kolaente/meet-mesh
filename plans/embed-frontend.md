# Plan: Embed Frontend into API Binary

## Overview

Embed the SvelteKit static build into the Go API binary using Go's `//go:embed` directive, enabling single-binary deployment.

## Current State

- **Backend**: Go API using ogen-generated server at `/api/v1/*`
- **Frontend**: SvelteKit SPA using `adapter-static`, outputs to `frontend/build/`
- **Makefile**: Already copies frontend build to `api/embedded/dist` (but embedding not implemented)
- **Server**: `api/cmd/main.go` uses `http.ListenAndServe` with the ogen-generated handler

## Implementation Steps

### Step 1: Create the embedded package

Create `api/embedded/embedded.go` with the embed directive:

```go
package embedded

import "embed"

//go:embed dist/*
var Frontend embed.FS
```

Also create a placeholder file `api/embedded/dist/.gitkeep` so the directory exists for `go:embed` to work before the first build.

### Step 2: Create a static file server handler

Create `api/static.go` to serve the embedded SPA:

```go
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
```

### Step 3: Create a multiplexer to combine API and static handlers

Modify `api/cmd/main.go` to route requests:

```go
// Create a mux that routes /api/* to the API server
// and everything else to the static file server
mux := http.NewServeMux()

// API routes - the ogen server handles /api/v1/*
mux.Handle("/api/", server)

// Static file server for everything else
mux.Handle("/", api.NewStaticHandler())

log.Printf("Meet Mesh starting on port %d...", cfg.Server.Port)
log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), mux))
```

### Step 4: Update .gitignore

Add to `.gitignore`:
```
api/embedded/dist/
!api/embedded/dist/.gitkeep
```

### Step 5: Verify the build process

The existing Makefile already handles this correctly:
1. `frontend-dist` builds the SPA and copies to `api/embedded/dist`
2. `build-api` compiles the Go binary with embedded files

Run `make build` to create the single binary.

## File Changes Summary

| File | Action |
|------|--------|
| `api/embedded/embedded.go` | Create - embed directive |
| `api/embedded/dist/.gitkeep` | Create - placeholder for git |
| `api/static.go` | Create - static file handler |
| `api/cmd/main.go` | Modify - add mux routing |
| `.gitignore` | Modify - ignore dist except .gitkeep |

## Testing

1. Run `make build` to create the binary
2. Run `./ai-feed -config config.yaml`
3. Access `http://localhost:8080` - should serve the frontend
4. Access `http://localhost:8080/api/v1/...` - should serve API

## Considerations

### Cache headers
Consider adding cache headers for static assets:
- Immutable hashed assets (`/_app/immutable/*`): `Cache-Control: public, max-age=31536000, immutable`
- Other assets: `Cache-Control: public, max-age=3600`

### Compression
For production, consider adding gzip middleware. Options:
- Use `github.com/NYTimes/gziphandler`
- Or pre-compress files during build and serve `.gz` variants

### Development mode
For development, you'll still run frontend and API separately (`pnpm dev` + `go run`). The embedded version is for production only.

## Rollback

If issues occur, simply revert the changes to `api/cmd/main.go` and remove the new files. The API will return to serving only API endpoints.
