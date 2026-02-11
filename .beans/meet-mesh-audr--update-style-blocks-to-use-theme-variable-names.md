---
# meet-mesh-audr
title: Update style blocks to use @theme variable names
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:37:40Z
updated_at: 2026-02-11T18:58:06Z
parent: meet-mesh-gpjw
blocked_by:
    - meet-mesh-4rqq
---

# Update style blocks that reference old CSS variable names

These files use `var(--bg-primary)`, `var(--sidebar-width)`, `var(--border)`, etc. in `<style>` blocks (raw CSS, not Tailwind classes). After the `@theme` rewrite, the old variable names no longer exist.

## Files to Modify

- `frontend/src/lib/components/dashboard/Sidebar.svelte` (lines 116-299)
- `frontend/src/routes/(dashboard)/+layout.svelte` (lines 44-59)
- `frontend/src/routes/(dashboard)/+page.svelte` (inline `<style>` block)

## Search-and-Replace Mapping

| Old | New |
|-----|-----|
| `var(--sidebar-width)` | `var(--spacing-sidebar)` |
| `var(--bg-primary)` | `var(--color-bg-primary)` |
| `var(--bg-secondary)` | `var(--color-bg-secondary)` |
| `var(--bg-tertiary)` | `var(--color-bg-tertiary)` |
| `var(--text-primary)` | `var(--color-text-primary)` |
| `var(--text-secondary)` | `var(--color-text-secondary)` |
| `var(--text-muted)` | `var(--color-text-muted)` |
| `var(--border-color)` | `var(--color-border)` |
| `var(--border)` | `2px solid var(--color-border)` |
| `var(--border-light)` | `1px solid var(--color-border)` |
| `var(--cyan)` | `var(--color-accent-sky)` |
| `var(--orange)` | `var(--color-accent-amber)` |
| `var(--pink)` | `var(--color-accent-rose)` |
| `var(--shadow-sm)` | `var(--shadow-brutalist-sm)` |
| `var(--shadow)` | `var(--shadow-brutalist)` |
| `var(--radius)` | `var(--radius-brutalist)` |
| `var(--transition)` | `0.15s ease` |
| `var(--amber)` | `var(--color-accent-amber)` |
| `var(--sky)` | `var(--color-accent-sky)` |

**Important:** Be careful with `var(--border)` vs `var(--border-color)` - they are different variables. `--border` was a compound shorthand (e.g., `2px solid #1a1a1a`), while `--border-color` was just the color.

## Steps

1. Update all three files' `<style>` blocks with the mapping above.
2. Run: `cd frontend && pnpm build`
3. Commit: `git commit -m "refactor: update style blocks to use @theme variable names"`

## Summary of Changes

- Updated Sidebar.svelte style block to use @theme variable names
- Updated +layout.svelte style block to use @theme variable names  
- Updated +page.svelte style block to use @theme variable names
- All var(--border) replaced with 2px solid var(--color-border)
- All var(--border-light) replaced with 1px solid var(--color-border)
- All var(--transition) replaced with 0.15s ease
