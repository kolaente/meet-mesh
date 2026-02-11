---
# meet-mesh-us06
title: Remove date/time section from calendar settings and add auth store refresh
status: completed
type: task
priority: normal
created_at: 2026-02-11T20:05:00Z
updated_at: 2026-02-11T19:25:28Z
parent: meet-mesh-us01
blocked_by:
    - meet-mesh-us05
---

# Cleanup Calendar Settings Page and Add Auth Store Refresh

**Goal:** Complete the settings restructuring by (1) removing the Date & Time Format section from the calendar settings page (it now lives on the account settings page), (2) adding a `refreshUser` method to the auth store so the sidebar/header immediately reflects name changes after saving.

**Architecture:** Two small changes in parallel:
- Remove the Date & Time section from `frontend/src/routes/(dashboard)/settings/+page.svelte` (the calendar settings page)
- Add a `refreshUser(user)` method to the auth store at `frontend/src/lib/stores/auth.svelte.ts`
- Update the sidebar to use the user's name from the auth store (if set) instead of parsing it from email

---

## Files

- Modify: `frontend/src/routes/(dashboard)/settings/+page.svelte` (remove date/time section and related imports)
- Modify: `frontend/src/lib/stores/auth.svelte.ts` (add refreshUser method)
- Modify: `frontend/src/lib/components/dashboard/Sidebar.svelte` (use user.name if available)

---

## Step 1: Remove Date & Time Format section from calendar settings page

Open `frontend/src/routes/(dashboard)/settings/+page.svelte`.

**1a. Remove the dateFormat import and related code from the script block.**

Remove these lines from the `<script>` block:

```typescript
import { getDateFormat } from '$lib/stores/dateFormat.svelte';
```

```typescript
const dateFormat = getDateFormat();

// Reactive bindings for select components
let timeFormatValue = $derived(dateFormat.timeFormat);
let weekStartDayValue = $derived(dateFormat.weekStartDay);

const timeFormatOptions = [
  { value: '12h', label: '12-hour (2:00 PM)' },
  { value: '24h', label: '24-hour (14:00)' }
];

const weekStartOptions = [
  { value: 'sunday', label: 'Sunday' },
  { value: 'monday', label: 'Monday' }
];
```

Also remove `Select` from the component imports since it is no longer used on this page:

Change:
```typescript
import { Card, Button, Spinner, Select } from '$lib/components/ui';
```
To:
```typescript
import { Card, Button, Spinner } from '$lib/components/ui';
```

**1b. Remove the Date & Time Format `<section>` from the template.**

Remove the entire section (lines 124-172 in the current file):

```svelte
<!-- Date & Time Format Section -->
<section>
  <div class="mb-4">
    <h2 class="text-lg font-medium text-[var(--text-primary)]">Date & Time Format</h2>
    <p class="text-sm text-[var(--text-secondary)]">Customize how dates and times are displayed</p>
  </div>

  <Card>
    {#snippet children()}
      ... (entire card contents) ...
    {/snippet}
  </Card>
</section>
```

**1c. Update the page header title.**

Change the header from:
```svelte
<DashboardHeader title="Settings">
```
To:
```svelte
<DashboardHeader title="Calendar Settings">
```

This makes it clear that this page is specifically for calendar connection settings, since the account/profile settings now live at `/settings/account`.

After these changes, the calendar settings page should only contain:
1. DashboardHeader with "Calendar Settings" title and "Add Calendar" button
2. Calendar Connections section (the CalDAV connection cards)
3. AddCalendarDialog

---

## Step 2: Add refreshUser method to auth store

Open `frontend/src/lib/stores/auth.svelte.ts`.

The current store looks like:

```typescript
type User = { id: number; email: string; name?: string }

let user = $state<User | null>(null)
let loading = $state(true)
let checked = $state(false)

export function getAuth() {
  return {
    get user() { return user },
    get loading() { return loading },
    get isAuthenticated() { return user !== null },

    async check() {
      if (checked) return
      loading = true
      const { data } = await api.GET('/auth/me')
      user = data ?? null
      loading = false
      checked = true
    },

    async logout() {
      await api.POST('/auth/logout')
      user = null
      window.location.href = '/auth/login'
    }
  }
}
```

Add a `refreshUser` method to the returned object. This method accepts a user object (from the API response) and updates the local state without making another network request:

```typescript
refreshUser(updatedUser: User) {
  user = updatedUser
},
```

Add it after the `isAuthenticated` getter. The full store becomes:

```typescript
import { api } from '$lib/api/client'

type User = { id: number; email: string; name?: string }

let user = $state<User | null>(null)
let loading = $state(true)
let checked = $state(false)

export function getAuth() {
  return {
    get user() { return user },
    get loading() { return loading },
    get isAuthenticated() { return user !== null },

    refreshUser(updatedUser: User) {
      user = updatedUser
    },

    async check() {
      if (checked) return
      loading = true
      const { data } = await api.GET('/auth/me')
      user = data ?? null
      loading = false
      checked = true
    },

    async logout() {
      await api.POST('/auth/logout')
      user = null
      window.location.href = '/auth/login'
    }
  }
}
```

---

## Step 3: Update Sidebar to use user.name if available

Open `frontend/src/lib/components/dashboard/Sidebar.svelte`.

The current `getUserDisplayName` function (line 49-54) always parses the name from the email:

```typescript
function getUserDisplayName(): string {
  const email = auth.user?.email ?? ''
  if (!email) return 'User'
  const name = email.split('@')[0]
  return name.split(/[._-]/).map(p => p.charAt(0).toUpperCase() + p.slice(1)).join(' ')
}
```

Update it to prefer the user's actual name if set:

```typescript
function getUserDisplayName(): string {
  if (auth.user?.name) return auth.user.name
  const email = auth.user?.email ?? ''
  if (!email) return 'User'
  const name = email.split('@')[0]
  return name.split(/[._-]/).map(p => p.charAt(0).toUpperCase() + p.slice(1)).join(' ')
}
```

Similarly, update `getUserInitials` to use the real name if available:

```typescript
function getUserInitials(): string {
  if (auth.user?.name) {
    const parts = auth.user.name.trim().split(/\s+/)
    if (parts.length >= 2) {
      return (parts[0][0] + parts[1][0]).toUpperCase()
    }
    return auth.user.name.substring(0, 2).toUpperCase()
  }
  const email = auth.user?.email ?? ''
  if (!email) return '?'
  const name = email.split('@')[0]
  const parts = name.split(/[._-]/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}
```

---

## Step 4: Uncomment refreshUser call in account settings page (if it was commented out)

If `auth.refreshUser(data)` was commented out in meet-mesh-us05 due to the method not existing yet, uncomment it now in `frontend/src/routes/(dashboard)/settings/account/+page.svelte`.

---

## Step 5: Verify the build compiles

Run:

```bash
cd frontend && pnpm check
```

Expected: No type errors.

---

## Step 6: Manual test

Run: `cd frontend && pnpm dev` (and `cd api && go run ./cmd` for the backend)

Verify:
1. Navigate to `/settings` -- the page shows "Calendar Settings" header and only the Calendar Connections section (no Date & Time section)
2. Navigate to `/settings/account` -- the page shows the Profile section and the Date & Time Format section
3. Change the display name on the account page and click "Save Profile"
4. The sidebar immediately reflects the new name (no page reload needed)
5. The initials in the avatar update based on the new name
6. Refresh the page -- the name persists (loaded from API)

---

## Step 7: Commit

```bash
git add frontend/src/routes/\(dashboard\)/settings/+page.svelte \
       frontend/src/lib/stores/auth.svelte.ts \
       frontend/src/lib/components/dashboard/Sidebar.svelte
git commit -m "feat(frontend): move date/time settings to account page, add auth refresh, use display name in sidebar"
```

## Summary of Changes

1. Removed Date & Time Format section from calendar settings page:
   - Removed dateFormat import and related code
   - Removed Select component import (no longer used)
   - Removed the Date & Time Format section from template
   - Updated header title from "Settings" to "Calendar Settings"

2. Added refreshUser method to auth store:
   - Added refreshUser(updatedUser) method to update local user state
   - Allows immediate UI updates after profile changes without page reload

3. Updated Sidebar to use user.name if available:
   - Updated getUserInitials() to prefer user.name over email parsing
   - Updated getUserDisplayName() to return user.name if set

4. Updated account settings page to call auth.refreshUser(data) after save

Build passes with pnpm check (0 errors)
