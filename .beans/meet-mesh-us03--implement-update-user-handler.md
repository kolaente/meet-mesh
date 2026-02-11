---
# meet-mesh-us03
title: Implement UpdateCurrentUser handler in Go backend
status: completed
type: task
priority: normal
created_at: 2026-02-11T20:02:00Z
updated_at: 2026-02-11T19:19:14Z
parent: meet-mesh-us01
blocked_by:
    - meet-mesh-us02
---

# Implement UpdateCurrentUser Handler

**Goal:** Implement the `UpdateCurrentUser` method on the Handler struct so the `PUT /auth/me` endpoint actually persists the user's updated name to the database and returns the updated user.

**Architecture:** Follow the same pattern as `GetCurrentUser` in `api/handler_auth.go`: extract user ID from context, fetch user from DB, apply updates, save, return the API response type. The User model already has a `Name` field in `api/models.go`, so no schema migration is needed.

**Prerequisite:** meet-mesh-us02 (OpenAPI spec + code generation must be done first)

---

## Files

- Modify: `api/handler_auth.go` (add `UpdateCurrentUser` method)

---

## Step 1: Understand the existing pattern

Look at the existing `GetCurrentUser` method in `api/handler_auth.go` (lines 113-129):

```go
func (h *Handler) GetCurrentUser(ctx context.Context) (gen.GetCurrentUserRes, error) {
    userID, ok := GetUserID(ctx)
    if !ok {
        return &gen.Error{Message: "Not authenticated"}, nil
    }

    var user User
    if err := h.db.First(&user, userID).Error; err != nil {
        return &gen.Error{Message: "User not found"}, nil
    }

    return &gen.User{
        ID:    int(user.ID),
        Email: user.Email,
        Name:  gen.NewOptString(user.Name),
    }, nil
}
```

The new handler follows the same pattern but accepts a request body and updates the user.

---

## Step 2: Check the generated request/response types

After running `make generate` (from meet-mesh-us02), check what types ogen generated:

```bash
grep -n "UpdateCurrentUser" api/gen/oas_server_gen.go
grep -n "UpdateCurrentUserReq\|UpdateCurrentUserRes" api/gen/oas_schemas_gen.go
```

The generated request type will likely be something like `UpdateCurrentUserReq` with an optional `Name` field. The response type will be `UpdateCurrentUserRes` (an interface satisfied by `*gen.User` and `*gen.Error`).

---

## Step 3: Implement the UpdateCurrentUser method

Add the following method to `api/handler_auth.go`, after the existing `GetCurrentUser` method:

```go
// UpdateCurrentUser updates the authenticated user's profile
func (h *Handler) UpdateCurrentUser(ctx context.Context, req *gen.UpdateCurrentUserReq) (gen.UpdateCurrentUserRes, error) {
    userID, ok := GetUserID(ctx)
    if !ok {
        return &gen.Error{Message: "Not authenticated"}, nil
    }

    var user User
    if err := h.db.First(&user, userID).Error; err != nil {
        return &gen.Error{Message: "User not found"}, nil
    }

    // Apply updates
    if req.Name.Set {
        user.Name = req.Name.Value
    }

    if err := h.db.Save(&user).Error; err != nil {
        return &gen.Error{Message: "Failed to update user"}, nil
    }

    return &gen.User{
        ID:    int(user.ID),
        Email: user.Email,
        Name:  gen.NewOptString(user.Name),
    }, nil
}
```

**Important notes on the generated types:**

- The `req` parameter type name depends on what ogen generates. It might be `*gen.UpdateCurrentUserReq` or `*gen.OptUpdateCurrentUserReq`. Check `api/gen/oas_schemas_gen.go` for the exact type.
- ogen wraps optional fields in `gen.OptString` (with `.Set` boolean and `.Value` string). Access the name via `req.Name.Set` and `req.Name.Value`.
- If the generated type differs from the above, adjust the field access accordingly. The key pattern is: check if the field was provided (`Set`), then use its value.

---

## Step 4: Verify the build compiles

Run:

```bash
cd api && go build ./cmd
```

Expected: Build succeeds with no errors.

---

## Step 5: Test manually (optional but recommended)

Start the dev server:

```bash
cd api && go run ./cmd
```

Test with curl (replace the session cookie with a valid one):

```bash
curl -X PUT http://localhost:8080/api/auth/me \
  -H "Content-Type: application/json" \
  -H "Cookie: session=YOUR_SESSION_COOKIE" \
  -d '{"name": "Test User"}'
```

Expected response:

```json
{"id": 1, "email": "user@example.com", "name": "Test User"}
```

Then verify the name persisted by calling GET:

```bash
curl http://localhost:8080/api/auth/me \
  -H "Cookie: session=YOUR_SESSION_COOKIE"
```

Expected: The response includes `"name": "Test User"`.

---

## Step 6: Commit

```bash
git add api/handler_auth.go
git commit -m "feat(api): implement UpdateCurrentUser handler for PUT /auth/me"
```

## Summary of Changes

Implemented UpdateCurrentUser handler in api/handler_auth.go:
- Extracts user ID from context
- Fetches user from database
- Updates name field if provided in request
- Persists changes to database
- Returns updated User response

Build passes successfully.
