---
# meet-mesh-sjc2
title: Replace arbitrary CSS vars with theme utilities in UI components
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:37:40Z
updated_at: 2026-02-11T19:01:49Z
parent: meet-mesh-gpjw
blocked_by:
    - meet-mesh-4rqq
---

# Replace arbitrary values in UI components

Replace `[var(--...)]` arbitrary value patterns with proper Tailwind v4 theme utilities in all UI components.

## Files to Modify (9 files)

- `frontend/src/lib/components/ui/Button.svelte`
- `frontend/src/lib/components/ui/Card.svelte`
- `frontend/src/lib/components/ui/Dialog.svelte`
- `frontend/src/lib/components/ui/Input.svelte`
- `frontend/src/lib/components/ui/Select.svelte`
- `frontend/src/lib/components/ui/Skeleton.svelte`
- `frontend/src/lib/components/ui/Textarea.svelte`
- `frontend/src/lib/components/ui/Toast.svelte`
- `frontend/src/lib/components/ui/Popover.svelte`

## Search-and-Replace Mapping for Color Patterns

These apply with all Tailwind variants (e.g., `hover:`, `focus:`, `active:`, `data-[highlighted]:`):

| Old | New |
|-----|-----|
| `bg-[var(--bg-primary)]` | `bg-bg-primary` |
| `bg-[var(--bg-secondary)]` | `bg-bg-secondary` |
| `bg-[var(--bg-tertiary)]` | `bg-bg-tertiary` |
| `text-[var(--text-primary)]` | `text-text-primary` |
| `text-[var(--text-secondary)]` | `text-text-secondary` |
| `text-[var(--text-muted)]` | `text-text-muted` |
| `text-[var(--text-tertiary)]` | `text-text-tertiary` |
| `border-[var(--border-color)]` | `border-border` |
| `divide-[var(--border-color)]` | `divide-border` |
| `bg-[var(--sky)]` | `bg-accent-sky` |
| `text-[var(--sky)]` | `text-accent-sky` |
| `border-[var(--sky)]` | `border-accent-sky` |
| `ring-[var(--sky)]` | `ring-accent-sky` |
| `bg-[var(--emerald)]` | `bg-accent-emerald` |
| `border-[var(--emerald)]` | `border-accent-emerald` |
| `ring-[var(--emerald)]` | `ring-accent-emerald` |
| `bg-[var(--rose)]` | `bg-accent-rose` |
| `text-[var(--rose)]` | `text-accent-rose` |
| `ring-[var(--rose)]` | `ring-accent-rose` |
| `bg-[var(--amber)]` | `bg-accent-amber` |
| `ring-[var(--amber)]` | `ring-accent-amber` |
| `text-[var(--sky-hover)]` | `text-accent-sky-hover` |

## Steps

1. Find-and-replace in all 9 UI component files using the mapping above.
2. Run: `cd frontend && pnpm build`
3. Commit: `git commit -m "refactor: replace arbitrary CSS vars with theme utilities in UI components"`

## Summary of Changes

- Button.svelte: replaced rounded-[var(--radius)] with rounded-brutalist
- Card.svelte: replaced bg-[var(--bg-*)], border-[var(--border-*)], shadow-[var(--shadow)], rounded-[var(--radius-lg)] with theme utilities
- Dialog.svelte: replaced rounded-[var(--radius-lg)], rounded-[var(--radius-sm)] with theme utilities
- Input.svelte: replaced all text-[var(--text-*)], bg-[var(--bg-*)], border-[var(--border-*)], ring-[var(--sky)] with theme utilities
- Popover.svelte: replaced rounded-[var(--radius-md)] with rounded-brutalist-md
- Select.svelte: replaced all arbitrary value patterns with theme utilities
- Skeleton.svelte: replaced rounded-[var(--radius-sm)] with rounded-brutalist-sm
- Textarea.svelte: replaced all arbitrary value patterns with theme utilities
- Toast.svelte: replaced rounded-[var(--radius-md)], shadow-[var(--shadow-md)] with theme utilities
