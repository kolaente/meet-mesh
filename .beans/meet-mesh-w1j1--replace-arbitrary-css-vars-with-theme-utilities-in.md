---
# meet-mesh-w1j1
title: Replace arbitrary CSS vars with theme utilities in poll components
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:37:40Z
updated_at: 2026-02-11T19:08:31Z
parent: meet-mesh-gpjw
blocked_by:
    - meet-mesh-4rqq
---

# Replace arbitrary values in poll components

Replace `[var(--...)]` arbitrary value patterns with proper Tailwind v4 theme utilities in all poll components.

## Files to Modify (5 files)

- `frontend/src/lib/components/poll/VoteCard.svelte`
- `frontend/src/lib/components/poll/VoteButtons.svelte`
- `frontend/src/lib/components/poll/VoteSummary.svelte`
- `frontend/src/lib/components/poll/VoterForm.svelte`
- `frontend/src/lib/components/poll/PollPage.svelte`

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

1. Find-and-replace in all 5 poll component files using the mapping above.
2. Run: `cd frontend && pnpm build`
3. Commit: `git commit -m "refactor: replace arbitrary CSS vars with theme utilities in poll components"`

## Summary of Changes

Replaced arbitrary CSS var patterns with proper Tailwind theme utilities in 5 poll components:

- **VoteCard.svelte**: Updated text-muted, text-primary, text-secondary
- **VoteButtons.svelte**: Updated accent colors (emerald, rose, amber), focus rings, bg/text colors, borders, radius
- **VoteSummary.svelte**: Updated radius values
- **VoterForm.svelte**: Updated border and text-primary
- **PollPage.svelte**: Updated error message radius

Build passes successfully.
