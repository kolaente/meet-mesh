---
# meet-mesh-ksvx
title: Create .github/workflows/ci.yml with api and frontend jobs
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:30:15Z
updated_at: 2026-02-11T18:44:30Z
parent: meet-mesh-wzh8
blocked_by:
    - meet-mesh-3wbu
---

# Create .github/workflows/ci.yml

Single workflow file with two parallel jobs that use make targets.

## Workflow structure

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  api:
    runs-on: ubuntu-latest
    steps:
      - actions/checkout@v4
      - actions/setup-go@v5 (go-version: "1.25", cache-dependency-path: api/go.sum)
      - Install golangci-lint
      - make lint
      - make test
      - make build-api
      - make verify-codegen

  frontend:
    runs-on: ubuntu-latest
    steps:
      - actions/checkout@v4
      - pnpm/action-setup@v4 (version: 10)
      - actions/setup-node@v4 (node-version: "22", cache: pnpm, cache-dependency-path: frontend/pnpm-lock.yaml)
      - cd frontend && pnpm install --frozen-lockfile
      - make check
      - make frontend-build
      - make verify-codegen
```

## Key decisions

- **CGO_ENABLED=1** is required for SQLite (mattn/go-sqlite3) - set as env var on api job
- **golangci-lint:** Use `golangci/golangci-lint-action@v7` for caching, OR just `make lint` if golangci-lint is pre-installed. Prefer the action for better caching.
- **pnpm version 10** matches what the project currently uses
- **Node 22** is current LTS
- **`--frozen-lockfile`** ensures lockfile is not modified during install
- **`verify-codegen`** only needs to run the relevant generator per job, but using the shared make target is simpler and catches cross-dependency issues
- Both jobs run in parallel for faster CI

## Acceptance criteria

- Single `.github/workflows/ci.yml` file (not separate files per job)
- Both jobs use `make` targets for all build/lint/test steps
- CI passes on the current codebase

## Summary of Changes

Created .github/workflows/ci.yml with:
- Two parallel jobs: api and frontend
- API job: Go 1.25 setup, golangci-lint action, make test, make build-api, make verify-codegen
- Frontend job: pnpm 10, Node 22, pnpm install, make check, make frontend-build, make verify-codegen
- CGO_ENABLED=1 for SQLite support in API job
