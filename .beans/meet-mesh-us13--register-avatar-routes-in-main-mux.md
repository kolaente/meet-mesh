---
# meet-mesh-us13
title: Register avatar routes in main.go mux
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:34:21Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us10
    - meet-mesh-us11
    - meet-mesh-us12
---

# Register Avatar Routes in main.go Mux

**Goal:** Wire up the `AvatarHandler` in `api/cmd/main.go` so the avatar upload, delete, and serving endpoints are live.

**Architecture:** Create the `AvatarHandler` after the auth service is initialized. Register it on the HTTP mux before the ogen server handler, so `/api/avatars/` requests go to the avatar handler instead of ogen.

---

## Files

- Modify: `api/cmd/main.go`

---

## Step 1: Create AvatarHandler and register routes

Open `api/cmd/main.go`. After the line that creates the `handler` (currently `handler := api.NewHandler(db, auth, caldav, mailer, cfg)`), add:

```go
// Create avatar handler (plain HTTP, not ogen)
avatarHandler := api.NewAvatarHandler(db, auth, cfg.Storage.AvatarsPath)
```

Then, in the mux setup section, add the avatar route **before** the general `/api/` route:

```go
// Avatar routes (must be before /api/ to take priority)
mux.HandleFunc("/api/avatars/", avatarHandler.HandleAvatars)

// API routes - the ogen server handles /api/*
mux.Handle("/api/", server)
```

The order matters because Go's `ServeMux` uses longest-prefix matching, so `/api/avatars/` will match before `/api/`.

The full mux setup should look like:

```go
mux := http.NewServeMux()

// Avatar routes (plain HTTP, not ogen - must be registered before /api/)
mux.HandleFunc("/api/avatars/", avatarHandler.HandleAvatars)

// API routes - the ogen server handles /api/*
mux.Handle("/api/", server)

// Static file server for everything else
mux.Handle("/", api.NewStaticHandler())
```

---

## Step 2: Verify the full application compiles and starts

Run:

```bash
cd api && go build ./cmd
```

Expected: Compiles without errors.

If you have a `config.yaml` configured, you can test it briefly:

```bash
cd api && go run ./cmd
```

Then in another terminal:

```bash
# Test avatar serving (should return 404 since no avatar exists yet)
curl -v http://localhost:9090/api/avatars/test.jpg
```

Expected: 404 Not Found.

---

## Step 3: Commit

```bash
git add api/cmd/main.go
git commit -m "feat(api): register avatar upload/delete/serve routes in HTTP mux"
```

## Summary of Changes

- Created AvatarHandler after auth service initialization
- Registered /api/avatars/ route before /api/ route (longest-prefix matching)
- Route dispatcher sends POST to Upload, DELETE to Delete, GET with filename to Serve

Build passes with go build
