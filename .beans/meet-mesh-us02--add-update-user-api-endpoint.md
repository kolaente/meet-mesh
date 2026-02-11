---
# meet-mesh-us02
title: Add PUT /auth/me endpoint to OpenAPI spec and regenerate
status: completed
type: task
priority: normal
created_at: 2026-02-11T20:01:00Z
updated_at: 2026-02-11T19:18:31Z
parent: meet-mesh-us01
---

# Add PUT /auth/me Endpoint to OpenAPI Spec

**Goal:** Define a new `PUT /auth/me` endpoint in the OpenAPI spec that allows the authenticated organizer to update their profile (currently just `name`), then regenerate the Go server code.

**Architecture:** Add the endpoint definition to `api/openapi.yaml` under the existing `/auth/me` path. The endpoint accepts a JSON body with an optional `name` field, is protected by `cookieAuth`, and returns the updated `User` schema. After modifying the spec, run `make generate` to regenerate the ogen Go types and interfaces.

---

## Files

- Modify: `api/openapi.yaml` (add PUT method to `/auth/me` path)
- Regenerate: `api/gen/` (auto-generated, do not hand-edit)

---

## Step 1: Add PUT method to /auth/me path in openapi.yaml

Open `api/openapi.yaml` and find the `/auth/me` path (currently around line 406). It currently only has a `get` method. Add a `put` method to the same path object.

Add the following YAML block after the existing `get` method's closing response definition (after line 424), still inside the `/auth/me` path:

```yaml
    put:
      operationId: updateCurrentUser
      summary: Update current user profile
      security:
        - cookieAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                  description: Display name for the organizer
      responses:
        '200':
          description: Updated user
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '401':
          description: Not authenticated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
```

The full `/auth/me` path should now look like:

```yaml
  /auth/me:
    get:
      operationId: getCurrentUser
      summary: Get current user
      security:
        - cookieAuth: []
      responses:
        '200':
          description: Current user
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '401':
          description: Not authenticated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

    put:
      operationId: updateCurrentUser
      summary: Update current user profile
      security:
        - cookieAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                  description: Display name for the organizer
      responses:
        '200':
          description: Updated user
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '401':
          description: Not authenticated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
```

---

## Step 2: Regenerate Go server code

Run from the repo root:

```bash
make generate
```

This runs `go install github.com/ogen-go/ogen/cmd/ogen@latest` and then `cd api && go generate ./...`.

Expected: The `api/gen/` directory updates with new files/changes. Specifically, `oas_server_gen.go` should now contain an `UpdateCurrentUser` method in the `Handler` interface, and `oas_unimplemented_gen.go` should have a stub.

Verify: `grep -n "UpdateCurrentUser" api/gen/oas_server_gen.go`

Expected output (something like):

```
UpdateCurrentUser(ctx context.Context, req *UpdateCurrentUserReq) (UpdateCurrentUserRes, error)
```

---

## Step 3: Regenerate frontend TypeScript types

Run:

```bash
cd frontend && pnpm generate:api
```

Expected: `frontend/src/lib/api/generated/schema.d.ts` updates with the new `PUT /auth/me` path and `updateCurrentUser` operation.

---

## Step 4: Verify the Go code compiles (it will have a build error, which is expected)

Run:

```bash
cd api && go build ./cmd 2>&1 | head -5
```

Expected: A compile error mentioning `UpdateCurrentUser` is not implemented. This is expected -- the next task (meet-mesh-us03) implements the handler.

---

## Step 5: Commit

```bash
git add api/openapi.yaml api/gen/ frontend/src/lib/api/generated/
git commit -m "feat(api): add PUT /auth/me endpoint to OpenAPI spec for user profile updates"
```

## Summary of Changes

- Added PUT method to /auth/me path in openapi.yaml with updateCurrentUser operationId
- Request body accepts optional 'name' field
- Returns updated User schema on success
- Ran make generate to regenerate Go code
- Ran pnpm generate:api to regenerate frontend TypeScript types
- Verified UpdateCurrentUser is in oas_server_gen.go interface
