---
# meet-mesh-gpjw
title: Fix non-idiomatic Tailwind CSS variable usage
status: completed
type: epic
priority: normal
created_at: 2026-02-11T19:30:00Z
updated_at: 2026-02-11T19:16:15Z
---

# Fix Non-Idiomatic Tailwind CSS Variable Usage

**Goal:** Replace all `[var(--...)]` arbitrary value patterns in Tailwind classes with proper first-class Tailwind v4 theme utilities using the `@theme` directive.

**Architecture:** Tailwind CSS v4 uses a CSS-based configuration model where custom design tokens are defined via the `@theme` directive. Variables defined with correct namespace prefix (e.g. `--color-*`, `--radius-*`, `--shadow-*`) automatically generate utility classes. The current codebase defines CSS custom properties in `:root` but uses them via verbose arbitrary value syntax like `bg-[var(--bg-secondary)]` instead of idiomatic utilities like `bg-bg-secondary`.

**Scope:** ~35 affected Svelte files across UI components, booking components, poll components, dashboard components, and route pages.

## Child Tasks

1. **Rewrite app.css** - Add `@theme` block with design tokens and dark mode overrides (foundational change, all other tasks depend on this)
2. **Update style blocks** - Fix raw CSS `var()` references in Sidebar.svelte, layouts, etc. to use new `@theme` variable names
3. **UI components** - Replace arbitrary values in Button, Card, Dialog, Input, Select, Skeleton, Textarea, Toast, Popover
4. **Booking components** - Replace arbitrary values in BookingPage, TimeSlotList, AvailabilityEditor, etc.
5. **Poll components** - Replace arbitrary values in VoteCard, VoteButtons, VoteSummary, VoterForm, PollPage
6. **Dashboard & routes** - Replace arbitrary values in dashboard components and all route pages
7. **Radius, shadow, border-width** - Replace `rounded-[var(--...)]`, `shadow-[var(--...)]`, and `border-[length:var(--...)]` patterns across all files
8. **Final verification** - Grep for remaining `[var(--` patterns, run build check, visual test in both themes

## Summary of Changes

All child tasks completed:
- meet-mesh-4rqq: Rewrite app.css with @theme block and dark mode overrides
- meet-mesh-audr: Update style blocks to use @theme variable names
- meet-mesh-sjc2: Replace arbitrary CSS vars with theme utilities in UI components
- meet-mesh-agk9: Replace arbitrary CSS vars with theme utilities in booking components
- meet-mesh-w1j1: Replace arbitrary CSS vars with theme utilities in poll components
- meet-mesh-nrp6: Replace arbitrary CSS vars with theme utilities in dashboard and routes
- meet-mesh-vq49: Replace radius, shadow, and border-width arbitrary values
- meet-mesh-0xyr: Final verification and cleanup

All arbitrary CSS variable patterns have been migrated to idiomatic Tailwind CSS v4 theme utilities. Build passes with 0 errors.
