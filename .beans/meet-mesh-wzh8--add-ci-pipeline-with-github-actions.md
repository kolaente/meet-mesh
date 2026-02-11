---
# meet-mesh-wzh8
title: CI pipeline with GitHub Actions
status: completed
type: epic
priority: normal
created_at: 2026-02-11T18:18:23Z
updated_at: 2026-02-11T18:46:54Z
---

# CI Pipeline with GitHub Actions - Epic

**Goal:** Set up a GitHub Actions CI pipeline that runs linting, type checking, tests, and building on every PR and push to main.

## Architecture

- **Single workflow file:** `.github/workflows/ci.yml` with two parallel jobs (`api` and `frontend`)
- **Make-driven:** All CI steps use `make` targets so CI mirrors local dev exactly
- **Codegen freshness:** Both jobs verify that generated code matches the OpenAPI spec

## Make targets needed

The Makefile currently has: `build`, `build-api`, `frontend-dist`, `generate`, `clean`, `dev`.

New targets to add:
- `make lint` - runs golangci-lint on the api (requires `api/.golangci.yml`)
- `make test` - runs `go test -race ./...` in api
- `make check` - runs `pnpm check` in frontend (svelte-check)
- `make frontend-build` - runs `pnpm build` in frontend (without copying to embed dir)
- `make verify-codegen` - runs generators and checks for git diff

## CI matrix

| Check           | API job              | Frontend job              |
|-----------------|----------------------|---------------------------|
| Lint            | `make lint`          | (none yet)                |
| Type check      | (via go build)       | `make check`              |
| Tests           | `make test`          | (none yet)                |
| Build           | `make build-api`     | `make frontend-build`     |
| Codegen fresh   | `make verify-codegen`| `make verify-codegen`     |

## Child tasks

1. **Add missing make targets** - Add `lint`, `test`, `check`, `frontend-build`, `verify-codegen` targets to the Makefile; create `api/.golangci.yml`
2. **Create `.github/workflows/ci.yml`** - Single workflow file with `api` and `frontend` jobs using make targets
3. **Local verification and lint fixes** - Run all make targets locally, fix any lint issues found by golangci-lint

## Summary of Changes

All child tasks completed:
1. Added missing make targets (lint, test, check, frontend-build, verify-codegen) and golangci-lint config
2. Created .github/workflows/ci.yml with parallel api and frontend jobs
3. Fixed all golangci-lint issues in the codebase

The CI pipeline is now ready for use on GitHub.
