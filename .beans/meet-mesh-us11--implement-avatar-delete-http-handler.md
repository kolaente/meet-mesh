---
# meet-mesh-us11
title: Implement avatar delete HTTP handler (DELETE /api/avatars)
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:32:47Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us10
---

# Implement Avatar Delete HTTP Handler

**Goal:** Add a `DELETE /api/avatars` endpoint that removes the organizer's avatar file from disk and clears the `AvatarFilename` in the database.

**Architecture:** Add a `Delete` method to the existing `AvatarHandler` struct in `api/handler_avatar.go`. The handler authenticates the request, looks up the user, removes the file from disk, and clears the database field.

---

## Files

- Modify: `api/handler_avatar.go` (add Delete method)

---

## Step 1: Add the Delete method to AvatarHandler

Open `api/handler_avatar.go`. Add the following method:

```go
// Delete handles DELETE /api/avatars
func (h *AvatarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate the user
	userID, ok := h.authenticateRequest(r)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Get current user
	var user User
	if err := h.db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Nothing to delete
	if user.AvatarFilename == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"avatar_url":""}`)
		return
	}

	oldFilename := user.AvatarFilename

	// Clear avatar in database
	if err := h.db.Model(&user).Update("avatar_filename", "").Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Remove file from disk (best-effort, don't fail if file is already gone)
	os.Remove(filepath.Join(h.avatarsPath, oldFilename))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"avatar_url":""}`)
}
```

Make sure `"fmt"` and `"os"` and `"path/filepath"` are in the imports (they should be already from the Upload handler).

---

## Step 2: Verify

Run:

```bash
cd api && go vet ./...
```

Expected: No errors.

---

## Step 3: Commit

```bash
git add api/handler_avatar.go
git commit -m "feat(api): implement avatar delete handler"
```

## Summary of Changes

- Added Delete method to AvatarHandler
- Authenticates user via session cookie
- Clears avatar_filename in database
- Removes file from disk (best-effort)
- Returns empty avatar_url JSON

Build passes with go vet
