---
# meet-mesh-ixmo
title: Remove non-functional calendar and activity widgets from dashboard
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:16:57Z
updated_at: 2026-02-11T18:52:08Z
---

# Remove Non-Functional Dashboard Widgets Implementation Plan

**Goal:** Remove the non-functional MiniCalendar and ActivityFeed widgets from the dashboard start page, delete orphaned component files, and adjust the layout so the page still looks good with just the stats grid, recent bookings, and quick actions.

**Architecture:** The dashboard page (`+page.svelte`) uses a two-column grid layout. The right column currently stacks three cards (Quick Actions, MiniCalendar, ActivityFeed). After removing the two non-functional widgets, the right column will only contain Quick Actions, so the layout should be simplified. The two-column layout should be kept (it separates the primary content -- bookings -- from secondary content -- quick actions), but the right column will be a single card instead of a stack.

**Tech Stack:** SvelteKit (Svelte 5 runes), Tailwind CSS v4, TypeScript

---

## Task 1: Remove MiniCalendar and ActivityFeed card blocks from the dashboard page

**Files:**
- Modify: `frontend/src/routes/(dashboard)/+page.svelte` (lines 286-316)

**Step 1: Remove the Mini Calendar card block**

Delete the entire Mini Calendar `<Card>` block (lines 287-300 in the current file):

```svelte
<!-- DELETE THIS ENTIRE BLOCK -->
<!-- Mini Calendar -->
<Card>
  {#snippet header()}
    <div class="flex items-center justify-between w-full">
      <span class="card-title">Calendar</span>
      <a href="/calendar" class="card-link">View full</a>
    </div>
  {/snippet}

  {#snippet children()}
  <div class="component-fullbleed">
    <MiniCalendar events={eventDates} />
  </div>
{/snippet}
</Card>
```

**Step 2: Remove the Activity Feed card block**

Delete the entire Activity Feed `<Card>` block (lines 302-316 in the current file):

```svelte
<!-- DELETE THIS ENTIRE BLOCK -->
<!-- Activity Feed -->
<Card>
  {#snippet header()}
    <div class="flex items-center justify-between w-full">
      <span class="card-title">Recent Activity</span>
      <a href="#" class="card-link">View all</a>
    </div>
  {/snippet}

  {#snippet children()}
  <div class="component-fullbleed">
    <ActivityFeed activities={sampleActivities} />
  </div>
{/snippet}
</Card>
```

**Step 3: Verify the right column now only contains the Quick Actions card**

After removal, the right-column `<div>` (lines 272-317) should contain only:

```svelte
<div class="right-column">
  <!-- Quick Actions -->
  <Card>
    {#snippet header()}
      <span class="card-title">Quick Actions</span>
    {/snippet}

    {#snippet children()}
    <div class="component-fullbleed">
      <QuickActions actions={quickActions} />
    </div>
  {/snippet}
  </Card>
</div>
```

---

## Task 2: Remove dead code from the dashboard page script block

**Files:**
- Modify: `frontend/src/routes/(dashboard)/+page.svelte` (script block, lines 1-127)

**Step 1: Remove MiniCalendar and ActivityFeed from the import statement**

Change the import (lines 6-13) from:

```typescript
import {
  DashboardHeader,
  StatsCard,
  BookingRow,
  QuickActions,
  MiniCalendar,
  ActivityFeed
} from '$lib/components/dashboard';
```

To:

```typescript
import {
  DashboardHeader,
  StatsCard,
  BookingRow,
  QuickActions
} from '$lib/components/dashboard';
```

**Step 2: Remove the `eventDates` derived state**

Delete lines 28-31:

```typescript
// Extract event dates from bookings for the mini calendar
const eventDates = $derived(
  recentBookings.map((b) => new Date(b.slot.start_time))
);
```

**Step 3: Remove the `sampleActivities` constant**

Delete lines 33-59 (the entire `sampleActivities` array with its hardcoded fake data):

```typescript
// Sample activity data (would come from API in a real implementation)
const sampleActivities = [
  {
    type: 'booking' as const,
    text: 'Sarah Johnson requested a booking',
    time: '2 minutes ago',
    boldParts: ['Sarah Johnson']
  },
  // ... all entries
];
```

**Step 4: Verify the build**

Run: `cd frontend && pnpm check`
Expected: No type errors

---

## Task 3: Remove the `.card-link` CSS rule if no longer used

**Files:**
- Modify: `frontend/src/routes/(dashboard)/+page.svelte` (style block)

**Step 1: Check if `.card-link` is still used in the template**

After removing the two card blocks, search the template for `card-link`. Both "View full" and "View all" links were inside the removed blocks. If no other element uses `card-link`, remove the CSS rules (lines 361-370):

```css
/* DELETE IF NO LONGER USED */
.card-link {
  font-size: 0.8rem;
  color: var(--sky);
  text-decoration: none;
  font-weight: 600;
}

.card-link:hover {
  text-decoration: underline;
}
```

**Step 2: Commit**

```bash
git add frontend/src/routes/\(dashboard\)/+page.svelte
git commit -m "refactor: remove non-functional calendar and activity widgets from dashboard"
```

---

## Task 4: Delete orphaned component files

**Files:**
- Delete: `frontend/src/lib/components/dashboard/MiniCalendar.svelte`
- Delete: `frontend/src/lib/components/dashboard/ActivityFeed.svelte`

**Step 1: Verify no other files import these components**

Run: `grep -r "MiniCalendar\|ActivityFeed" frontend/src/`

After the previous task's changes, only `index.ts` should still reference them.

**Step 2: Delete the files**

```bash
rm frontend/src/lib/components/dashboard/MiniCalendar.svelte
rm frontend/src/lib/components/dashboard/ActivityFeed.svelte
```

---

## Task 5: Clean up barrel exports in index.ts

**Files:**
- Modify: `frontend/src/lib/components/dashboard/index.ts`

**Step 1: Remove the two export lines**

Remove these two lines from `frontend/src/lib/components/dashboard/index.ts`:

```typescript
export { default as ActivityFeed } from './ActivityFeed.svelte';
export { default as MiniCalendar } from './MiniCalendar.svelte';
```

The file should go from 12 exports to 10 exports. The remaining exports are all still used:
- `Sidebar` - used in dashboard layout
- `MobileMenu` - used in dashboard layout
- `DashboardHeader` - used in dashboard and other pages
- `BookingRow` - used in dashboard page
- `VoteRow` - used in poll detail pages
- `StatsCard` - used in dashboard page
- `CalendarConnectionCard` - used in settings page
- `AddCalendarDialog` - used in settings page
- `PollOptionManager` - used in poll create/edit pages
- `QuickActions` - used in dashboard page

**Step 2: Commit**

```bash
git add frontend/src/lib/components/dashboard/MiniCalendar.svelte
git add frontend/src/lib/components/dashboard/ActivityFeed.svelte
git add frontend/src/lib/components/dashboard/index.ts
git commit -m "refactor: delete orphaned MiniCalendar and ActivityFeed components"
```

---

## Task 6: Verify the build passes

**Step 1: Run the type checker**

```bash
cd frontend && pnpm check
```

Expected: No errors

**Step 2: Run the full frontend build**

```bash
cd frontend && pnpm build
```

Expected: Build succeeds with no errors

---

## Summary of changes

| Action | File |
|--------|------|
| Modify | `frontend/src/routes/(dashboard)/+page.svelte` - remove widget cards, dead code, unused imports, unused CSS |
| Delete | `frontend/src/lib/components/dashboard/MiniCalendar.svelte` |
| Delete | `frontend/src/lib/components/dashboard/ActivityFeed.svelte` |
| Modify | `frontend/src/lib/components/dashboard/index.ts` - remove 2 export lines |

## Layout after removal

The dashboard page will consist of:
1. **DashboardHeader** - title "Dashboard" with "New Link" and "New Poll" buttons (unchanged)
2. **Stats grid** - 4 stat cards in a row (unchanged)
3. **Two-column content grid** (1.6fr left, 1fr right):
   - **Left column:** Recent Bookings card (unchanged)
   - **Right column:** Quick Actions card only (the `.right-column` flex container stays, but now just has one child; this is fine as it is self-sizing)

The right column with only Quick Actions will be shorter than the left column, but this is acceptable and looks natural -- the Quick Actions card will align to the top of the grid cell. No additional layout changes are needed.

## Summary of Changes

- Removed MiniCalendar and ActivityFeed card blocks from dashboard page
- Removed dead code (eventDates, sampleActivities) from script block
- Removed unused imports (MiniCalendar, ActivityFeed)
- Removed unused .card-link CSS rules
- Deleted orphaned component files (MiniCalendar.svelte, ActivityFeed.svelte)
- Updated index.ts barrel exports
