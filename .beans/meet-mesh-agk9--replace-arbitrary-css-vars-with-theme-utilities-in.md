---
# meet-mesh-agk9
title: Replace arbitrary CSS vars with theme utilities in booking components
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:37:40Z
updated_at: 2026-02-11T19:06:45Z
parent: meet-mesh-gpjw
blocked_by:
    - meet-mesh-4rqq
---

# Replace arbitrary values in booking components

Replace `[var(--...)]` arbitrary value patterns with proper Tailwind v4 theme utilities in all booking components.

## Files to Modify (8 files)

- `frontend/src/lib/components/booking/BookingPage.svelte`
- `frontend/src/lib/components/booking/TimeSlotList.svelte`
- `frontend/src/lib/components/booking/AvailabilityEditor.svelte`
- `frontend/src/lib/components/booking/AvailabilityRuleForm.svelte`
- `frontend/src/lib/components/booking/DurationSelector.svelte`
- `frontend/src/lib/components/booking/DateCalendar.svelte`
- `frontend/src/lib/components/booking/DateRangePicker.svelte`
- `frontend/src/lib/components/booking/DayPicker.svelte`

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

1. Find-and-replace in all 8 booking component files using the mapping above.
2. Run: `cd frontend && pnpm build`
3. Commit: `git commit -m "refactor: replace arbitrary CSS vars with theme utilities in booking components"`

## Summary of Changes

Replaced arbitrary CSS var patterns with proper Tailwind theme utilities in 8 booking components:

- **AvailabilityEditor.svelte**: Updated text and background colors
- **AvailabilityRuleForm.svelte**: Updated radius, colors, focus states
- **BookingPage.svelte**: Updated header, step indicator, error message, time slot flow, full day flow, multi-day flow sections
- **DateCalendar.svelte**: Updated 3 rounded radius values
- **DateRangePicker.svelte**: Updated 3 rounded radius values
- **DayPicker.svelte**: Updated trigger and content radius values
- **DurationSelector.svelte**: Updated button colors and borders
- **TimeSlotList.svelte**: Updated radius, colors, and tertiary text

All arbitrary patterns like \ → \, \ → \, etc.

Build passes successfully.
