---
# meet-mesh-nrp6
title: Replace arbitrary CSS vars with theme utilities in dashboard and routes
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:37:40Z
updated_at: 2026-02-11T19:14:26Z
parent: meet-mesh-gpjw
blocked_by:
    - meet-mesh-4rqq
---

# Replace arbitrary values in dashboard components and route pages

Replace `[var(--...)]` arbitrary value patterns with proper Tailwind v4 theme utilities in dashboard components and all route pages.

## Files to Modify (13 files)

**Dashboard Components:**
- `frontend/src/lib/components/dashboard/CalendarConnectionCard.svelte`
- `frontend/src/lib/components/dashboard/AddCalendarDialog.svelte`

**Route Pages:**
- `frontend/src/routes/(public)/+layout.svelte`
- `frontend/src/routes/(dashboard)/+page.svelte`
- `frontend/src/routes/(dashboard)/settings/+page.svelte`
- `frontend/src/routes/(dashboard)/polls/+page.svelte`
- `frontend/src/routes/(dashboard)/polls/new/+page.svelte`
- `frontend/src/routes/(dashboard)/polls/[id]/+page.svelte`
- `frontend/src/routes/(dashboard)/polls/[id]/edit/+page.svelte`
- `frontend/src/routes/(dashboard)/booking-links/+page.svelte`
- `frontend/src/routes/(dashboard)/booking-links/[id]/+page.svelte`
- `frontend/src/routes/(dashboard)/booking-links/[id]/edit/+page.svelte`
- `frontend/src/routes/(dashboard)/booking-links/new/+page.svelte`

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

1. Find-and-replace in all 13 files using the mapping above.
2. Run: `cd frontend && pnpm build`
3. Commit: `git commit -m "refactor: replace arbitrary CSS vars with theme utilities in dashboard and routes"`

## Summary of Changes

Replaced all arbitrary CSS variable patterns with theme utilities in dashboard components and route pages:

**Dashboard Components:**
- CalendarConnectionCard.svelte: text-text-primary, text-text-secondary
- AddCalendarDialog.svelte: text-text-secondary, text-text-muted, text-text-primary, border-border, bg-bg-tertiary, bg-bg-secondary, bg-accent-sky

**Route Pages:**
- (public)/+layout.svelte: bg-bg-primary
- (dashboard)/+page.svelte: text-text-muted
- settings/+page.svelte: text-text-primary, text-text-secondary, text-text-muted, border-border
- polls/+page.svelte: text-text-muted, text-text-primary, text-text-secondary
- polls/new/+page.svelte: rounded-brutalist-md
- polls/[id]/+page.svelte: text-text-secondary, text-accent-sky, text-text-primary, divide-border
- polls/[id]/edit/+page.svelte: rounded-brutalist-md
- booking-links/+page.svelte: text-text-muted, text-text-primary, text-text-secondary
- booking-links/[id]/+page.svelte: text-text-secondary, text-text-primary, text-accent-sky, divide-border
- booking-links/[id]/edit/+page.svelte: rounded-brutalist-md
- booking-links/new/+page.svelte: rounded-brutalist-md

Build passes with pnpm check (0 errors, 4 existing warnings)
