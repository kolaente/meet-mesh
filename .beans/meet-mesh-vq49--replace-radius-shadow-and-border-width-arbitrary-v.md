---
# meet-mesh-vq49
title: Replace radius, shadow, and border-width arbitrary values
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:37:41Z
updated_at: 2026-02-11T19:15:08Z
parent: meet-mesh-gpjw
blocked_by:
    - meet-mesh-4rqq
---

# Replace radius, shadow, and border-width arbitrary values

Replace `rounded-[var(--...)]`, `shadow-[var(--...)]`, and `border-[length:var(--...)]` patterns across all affected files.

## Files

All files from the UI, booking, poll, dashboard, and route page groups that contain these patterns.

## Radius Search-and-Replace Mapping

| Old | New |
|-----|-----|
| `rounded-[var(--radius-sm)]` | `rounded-brutalist-sm` |
| `rounded-[var(--radius-md)]` | `rounded-brutalist-md` |
| `rounded-[var(--radius)]` | `rounded-brutalist` |
| `rounded-[var(--radius-lg)]` | `rounded-brutalist-lg` |
| `rounded-b-[var(--radius-lg)]` | `rounded-b-brutalist-lg` |
| `rounded-l-[var(--radius-md)]` | `rounded-l-brutalist-md` |
| `rounded-r-[var(--radius-md)]` | `rounded-r-brutalist-md` |

**Also update Svelte `class:` directive syntax:**

| Old | New |
|-----|-----|
| `class:rounded-l-[var(--radius-md)]` | `class:rounded-l-brutalist-md` |
| `class:rounded-r-[var(--radius-md)]` | `class:rounded-r-brutalist-md` |
| `class:border-[var(--border-color)]` | `class:border-border` |

## Shadow Search-and-Replace Mapping

| Old | New |
|-----|-----|
| `shadow-[var(--shadow)]` | `shadow-brutalist` |
| `shadow-[var(--shadow-md)]` | `shadow` |

## Border-Width Search-and-Replace Mapping

Note: `--border-width` was never defined in the original CSS. Use `2px` which matches the design system.

| Old | New |
|-----|-----|
| `border-[length:var(--border-width)]` | `border-2` |
| `border-b-[length:var(--border-width)]` | `border-b-2` |
| `border-t-[length:var(--border-width)]` | `border-t-2` |

**Known files with shadow/border-width patterns:**
- `frontend/src/lib/components/ui/Card.svelte`
- `frontend/src/lib/components/ui/Select.svelte`
- `frontend/src/lib/components/ui/Toast.svelte`

## Steps

1. Find-and-replace across all affected Svelte files using the mappings above.
2. Run: `cd frontend && pnpm build`
3. Commit: `git commit -m "refactor: replace arbitrary radius, shadow, and border-width values with theme utilities"`

## Summary of Changes

No additional changes needed - all radius, shadow, and border-width arbitrary values were already replaced in previous tasks:
- meet-mesh-sjc2: UI components 
- meet-mesh-agk9: booking components
- meet-mesh-w1j1: poll components
- meet-mesh-nrp6: dashboard and route pages

Verified with ripgrep that no \ patterns remain in the frontend/src directory.
