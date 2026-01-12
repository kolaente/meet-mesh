# API Route Standardization Plan

## Problem

The codebase has inconsistent API route prefixes:
- Backend mounts routes at `/api/v1/`
- OpenAPI spec and frontend expect `/api`

This causes routing mismatches where some endpoints work and others don't.

## Current State

| Location | Current Pattern | Status |
|----------|-----------------|--------|
| `api/cmd/main.go:64` | `/api/v1/` | Needs change |
| `api/mailer.go:64-65` | `/api/v1/actions/*` | Needs change |
| `config.example.yaml:13` | `/api/v1/auth/callback` | Needs change |
| `api/openapi.yaml:8` | `/api` | Correct |
| `frontend/src/lib/api/client.ts:5` | `/api` | Correct |
| `frontend/src/routes/(actions)/actions/approve/+page.svelte:21` | `/api/actions/*` | Correct |
| `frontend/src/routes/(actions)/actions/decline/+page.svelte:21` | `/api/actions/*` | Correct |
| `frontend/src/routes/auth/login/+page.svelte:6` | `/api/auth/login` | Correct |

## Implementation Steps

### Step 1: Update backend route mounting

**File**: `api/cmd/main.go`

Change:
```go
// API routes - the ogen server handles /api/v1/*
mux.Handle("/api/v1/", server)
```

To:
```go
// API routes - the ogen server handles /api/*
mux.Handle("/api/", server)
```

### Step 2: Update email action URLs

**File**: `api/mailer.go`

Change lines 64-65:
```go
approveURL := fmt.Sprintf("%s/api/v1/actions/approve?token=%s", m.baseURL, booking.ActionToken)
declineURL := fmt.Sprintf("%s/api/v1/actions/decline?token=%s", m.baseURL, booking.ActionToken)
```

To:
```go
approveURL := fmt.Sprintf("%s/api/actions/approve?token=%s", m.baseURL, booking.ActionToken)
declineURL := fmt.Sprintf("%s/api/actions/decline?token=%s", m.baseURL, booking.ActionToken)
```

### Step 3: Update OIDC redirect URI example

**File**: `config.example.yaml`

Change:
```yaml
redirect_uri: http://localhost:8080/api/v1/auth/callback
```

To:
```yaml
redirect_uri: http://localhost:8080/api/auth/callback
```

### Step 4: Check for actual config files

Verify if there's a `config.yaml` or environment variables with the old `/api/v1/` path that need updating. Users with existing configurations will need to update their OIDC redirect URIs.

## Summary

- **Files to modify**: 3
- **Lines to change**: ~5
- **Breaking change**: Yes - existing OIDC configurations will need their redirect URIs updated
