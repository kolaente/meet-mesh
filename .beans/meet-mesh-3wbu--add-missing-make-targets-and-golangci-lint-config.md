---
# meet-mesh-3wbu
title: Add missing make targets and golangci-lint config
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:30:02Z
updated_at: 2026-02-11T18:44:01Z
parent: meet-mesh-wzh8
---

# Add missing make targets and golangci-lint config

## 1. Create `api/.golangci.yml`

Minimal config that excludes the `gen/` directory (generated ogen code):

```yaml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck

issues:
  exclude-dirs:
    - gen
```

## 2. Add make targets to `Makefile`

Add these new targets:

- **`lint`** - `cd api && golangci-lint run ./...`
- **`test`** - `cd api && CGO_ENABLED=1 go test -race ./...`
- **`check`** - `cd frontend && pnpm check` (runs svelte-kit sync + svelte-check)
- **`frontend-build`** - `cd frontend && pnpm build` (build only, no copy to embed dir)
- **`verify-codegen`** - Runs both `make generate` and `cd frontend && pnpm generate:api`, then checks for git diff. Fails if generated code is out of date.

All targets should be added to `.PHONY`.

## Acceptance criteria

- `make lint` runs golangci-lint against api code (skipping gen/)
- `make test` runs Go tests with race detector
- `make check` runs svelte-check on frontend
- `make frontend-build` builds the frontend without copying to embed dir
- `make verify-codegen` detects stale generated code

## Summary of Changes

- Created \ with minimal config excluding the \ directory
- Added new make targets to \:
  - \ - runs golangci-lint on api code
  - \ - runs Go tests with race detector
  - \ - runs svelte-check on frontend
  - \ - builds frontend without copying to embed dir
  - \ - runs generators and checks for git diff
