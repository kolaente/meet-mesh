---
# meet-mesh-us09
title: Add AvatarFilename to User model and avatar_url to OpenAPI User schema
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:29:10Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us08
---

# Add AvatarFilename to User Model and avatar_url to OpenAPI User Schema

**Goal:** Add an `AvatarFilename` field to the GORM User model and an `avatar_url` field to the OpenAPI User schema, so the backend can track which avatar file belongs to a user and the frontend can display it.

**Architecture:** The GORM model stores just the filename (e.g., `a1b2c3d4.jpg`). The API response returns a full URL (e.g., `/api/avatars/a1b2c3d4.jpg`). The handler constructs the URL from the filename when building the response. GORM AutoMigrate will add the column to the SQLite database automatically on next startup.

---

## Files

- Modify: `api/models.go` (add AvatarFilename field to User struct)
- Modify: `api/openapi.yaml` (add avatar_url to User schema)
- Modify: `api/handler_auth.go` (include avatar_url in GetCurrentUser response)
- Regenerate: `api/gen/` (run `make generate`)
- Regenerate: `frontend/src/lib/api/generated/` (run `pnpm generate:api`)

---

## Step 1: Add AvatarFilename to User model

Open `api/models.go`. Add the `AvatarFilename` field to the `User` struct:

```go
type User struct {
	ID             uint      `gorm:"primaryKey"`
	OIDCSub        string    `gorm:"column:oidc_sub;uniqueIndex;not null"`
	Email          string    `gorm:"not null"`
	Name           string
	AvatarFilename string
	CreatedAt      time.Time
	Calendars      []CalendarConnection `gorm:"foreignKey:UserID"`
	BookingLinks   []BookingLink        `gorm:"foreignKey:UserID"`
	Polls          []Poll               `gorm:"foreignKey:UserID"`
}
```

---

## Step 2: Add avatar_url to OpenAPI User schema

Open `api/openapi.yaml`. Find the `User` schema (around line 98). Add `avatar_url` as an optional string field:

```yaml
    User:
      type: object
      required: [id, email]
      properties:
        id:
          type: integer
        email:
          type: string
        name:
          type: string
        avatar_url:
          type: string
          description: URL to the user's avatar image, or empty if no avatar is set
```

---

## Step 3: Regenerate code

Run from repo root:

```bash
make generate
cd frontend && pnpm generate:api
```

Verify: `grep -n "AvatarURL" api/gen/oas_schemas_gen.go` should show the new field in the generated User struct.

---

## Step 4: Update GetCurrentUser handler to include avatar_url

Open `api/handler_auth.go`. In the `GetCurrentUser` method, update the return to include the avatar URL. Add a helper function to construct the avatar URL from a filename:

```go
// avatarURL constructs the full avatar URL from a filename.
// Returns empty string if no avatar is set.
func avatarURL(filename string) string {
	if filename == "" {
		return ""
	}
	return "/api/avatars/" + filename
}
```

Then update the `GetCurrentUser` response:

```go
return &gen.User{
    ID:        int(user.ID),
    Email:     user.Email,
    Name:      gen.NewOptString(user.Name),
    AvatarURL: gen.NewOptString(avatarURL(user.AvatarFilename)),
}, nil
```

Note: If `UpdateCurrentUser` exists from meet-mesh-us03, update its response too with the same `AvatarURL` field.

---

## Step 5: Update auth store type

Open `frontend/src/lib/stores/auth.svelte.ts`. Update the User type:

```typescript
type User = { id: number; email: string; name?: string; avatar_url?: string }
```

---

## Step 6: Verify

Run:

```bash
cd api && go build ./cmd
cd frontend && pnpm check
```

Expected: Both compile without errors.

---

## Step 7: Commit

```bash
git add api/models.go api/openapi.yaml api/handler_auth.go api/gen/ frontend/src/lib/api/generated/ frontend/src/lib/stores/auth.svelte.ts
git commit -m "feat(api): add avatar_url to User model and API schema"
```

## Summary of Changes

- Added AvatarFilename field to User struct in api/models.go
- Added avatar_url property to User schema in api/openapi.yaml
- Regenerated Go code with make generate
- Regenerated TypeScript types with pnpm generate:api
- Added avatarURL helper function in api/handler_auth.go
- Updated GetCurrentUser and UpdateCurrentUser to include AvatarURL in response
- Updated User type in frontend auth store to include avatar_url

Both Go and frontend builds pass
