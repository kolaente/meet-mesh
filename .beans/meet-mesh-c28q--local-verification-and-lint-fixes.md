---
# meet-mesh-c28q
title: Local verification and lint fixes
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:30:22Z
updated_at: 2026-02-11T18:46:47Z
parent: meet-mesh-wzh8
blocked_by:
    - meet-mesh-ksvx
---

# Local verification and lint fixes

Run all the new make targets locally and fix any issues before pushing.

## Steps

### 1. Verify all make targets work

```bash
make lint
make test
make check
make frontend-build
make verify-codegen
```

All should exit 0. If not, fix the issues.

### 2. Fix golangci-lint issues

If `make lint` reports issues, fix them in the Go source files. Common issues:
- Unused variables/imports
- Unchecked errors
- Simplifiable code

Commit fixes separately: `git commit -m "fix: resolve golangci-lint issues"`

### 3. Push and verify on GitHub

Push the branch and verify both CI jobs (api and frontend) appear and pass on GitHub.

## Acceptance criteria

- All make targets pass locally
- Any lint issues are fixed and committed
- CI runs green on GitHub

## Summary of Changes

Fixed all golangci-lint issues:
- Added explicit error ignoring with \ for mailer send calls (fire-and-forget pattern)
- Added explicit error ignoring for \ calls (crypto/rand never fails on Linux)
- Added explicit error ignoring for \ call
- Added explicit error ignoring for \ call
- Removed unused \ function

All make targets now pass:
- \ - passes
- \ - passes
- \ - passes (5 warnings, 0 errors)
- \ - passes
- \ - passes
