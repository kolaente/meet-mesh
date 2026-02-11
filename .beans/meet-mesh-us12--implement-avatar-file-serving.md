---
# meet-mesh-us12
title: Implement avatar file serving (GET /api/avatars/{filename})
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:33:37Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us10
---

# Implement Avatar File Serving

**Goal:** Serve avatar images from the filesystem at `GET /api/avatars/{filename}` with appropriate caching headers. This endpoint is public (no auth required) because avatars are displayed on public booking/poll pages.

**Architecture:** Add a `Serve` method to the `AvatarHandler` struct. It extracts the filename from the URL path, sanitizes it to prevent directory traversal, and serves the file from the avatars directory. Since filenames are content-hashed, we set aggressive cache headers (`Cache-Control: public, max-age=31536000, immutable`).

---

## Files

- Modify: `api/handler_avatar.go` (add Serve method)

---

## Step 1: Add the Serve method to AvatarHandler

Open `api/handler_avatar.go`. Add the following method:

```go
// Serve handles GET /api/avatars/{filename}
// This endpoint is public (no auth) since avatars appear on public pages.
func (h *AvatarHandler) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract filename from path: /api/avatars/{filename}
	filename := filepath.Base(r.URL.Path)

	// Sanitize: only allow alphanumeric + dot + hyphen
	if filename == "" || filename == "." || filename == ".." {
		http.NotFound(w, r)
		return
	}

	// Ensure filename doesn't contain path separators
	if filepath.Dir(filename) != "." {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(h.avatarsPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Set cache headers - filenames are content-hashed so we can cache forever
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Content-Type", "image/jpeg")

	http.ServeFile(w, r, filePath)
}
```

---

## Step 2: Add a ServeHTTP-style mux handler

To make routing cleaner, add a method that acts as a combined handler dispatching based on HTTP method. This will be used when registering on the mux:

```go
// HandleAvatars is the main handler for /api/avatars/ that dispatches by method.
func (h *AvatarHandler) HandleAvatars(w http.ResponseWriter, r *http.Request) {
	// /api/avatars/ with no filename = upload or delete
	// /api/avatars/{filename} = serve
	filename := strings.TrimPrefix(r.URL.Path, "/api/avatars/")
	filename = strings.TrimPrefix(filename, "/api/avatars")

	switch {
	case r.Method == http.MethodPost:
		h.Upload(w, r)
	case r.Method == http.MethodDelete:
		h.Delete(w, r)
	case r.Method == http.MethodGet && filename != "" && filename != "/":
		h.Serve(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
```

Add `"strings"` to the imports if not already present.

---

## Step 3: Verify

Run:

```bash
cd api && go vet ./...
```

Expected: No errors.

---

## Step 4: Commit

```bash
git add api/handler_avatar.go
git commit -m "feat(api): implement avatar file serving with cache headers"
```

## Summary of Changes

- Added Serve method for GET /api/avatars/{filename}
  - Public (no auth) since avatars appear on public pages
  - Extracts filename from path and sanitizes
  - Sets aggressive cache headers (immutable, 1 year)
  - Serves file with Content-Type: image/jpeg
- Added HandleAvatars dispatcher method
  - Routes POST -> Upload, DELETE -> Delete, GET with filename -> Serve

Build passes with go vet
