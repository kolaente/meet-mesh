---
# meet-mesh-us04
title: Restructure settings navigation in sidebar and mobile menu
status: completed
type: task
priority: normal
created_at: 2026-02-11T20:03:00Z
updated_at: 2026-02-11T19:20:59Z
parent: meet-mesh-us01
---

# Restructure Settings Navigation

**Goal:** Update the sidebar and mobile menu to show two settings sub-items -- "Account" (`/settings/account`) and "Calendar" (`/settings`) -- instead of the single "Calendar" entry. Add a settings layout with tab-style sub-navigation shared between settings pages.

**Architecture:**
1. Update `settingsItems` in `Sidebar.svelte` and the corresponding array in `MobileMenu.svelte` to have two entries under a "Settings" label.
2. Create a settings layout (`frontend/src/routes/(dashboard)/settings/+layout.svelte`) with tab-style sub-navigation that renders on all `/settings/*` pages.

---

## Files

- Modify: `frontend/src/lib/components/dashboard/Sidebar.svelte`
- Modify: `frontend/src/lib/components/dashboard/MobileMenu.svelte`
- Create: `frontend/src/routes/(dashboard)/settings/+layout.svelte`

---

## Step 1: Update Sidebar.svelte navigation items

Open `frontend/src/lib/components/dashboard/Sidebar.svelte`.

Change the `settingsItems` array (currently line 21-23) from:

```typescript
const settingsItems: NavItem[] = [
  { href: '/settings', label: 'Calendar', icon: 'calendar' }
]
```

To:

```typescript
const settingsItems: NavItem[] = [
  { href: '/settings/account', label: 'Account', icon: 'user' },
  { href: '/settings', label: 'Calendar', icon: 'calendar' }
]
```

Then in the template section where `settingsItems` are rendered (around line 92-101), add a "Settings" section label before the items, and add the SVG icon for the new `user` icon type.

Find the block:

```svelte
{#each settingsItems as item}
  <a href={item.href} class="nav-item" class:active={isActive(item.href)}>
    {#if item.icon === 'calendar'}
      <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
      </svg>
    {/if}
    {item.label}
  </a>
{/each}
```

Replace with:

```svelte
<div class="nav-section-label">Settings</div>
{#each settingsItems as item}
  <a href={item.href} class="nav-item" class:active={isActive(item.href)}>
    {#if item.icon === 'user'}
      <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
      </svg>
    {:else if item.icon === 'calendar'}
      <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
      </svg>
    {/if}
    {item.label}
  </a>
{/each}
```

Add CSS for the section label in the `<style>` block:

```css
.nav-section-label {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-muted);
  padding: 1rem 0.75rem 0.35rem;
}
```

Also update the `isActive` function to handle the `/settings` vs `/settings/account` ambiguity. The current implementation uses `startsWith`, which means `/settings` would match `/settings/account`. Fix by making the `/settings` check exact when there is no sub-path:

```typescript
function isActive(href: string): boolean {
  if (href === '/') {
    return page.url.pathname === '/'
  }
  if (href === '/settings') {
    return page.url.pathname === '/settings'
  }
  return page.url.pathname.startsWith(href)
}
```

---

## Step 2: Update MobileMenu.svelte navigation items

Open `frontend/src/lib/components/dashboard/MobileMenu.svelte`.

Change the `navItems` array (currently line 15-20) from:

```typescript
const navItems = [
  { href: '/', label: 'Dashboard', icon: 'dashboard' },
  { href: '/booking-links', label: 'Booking Links', icon: 'link' },
  { href: '/polls', label: 'Polls', icon: 'poll' },
  { href: '/settings', label: 'Calendar', icon: 'calendar' }
] as const;
```

To:

```typescript
const navItems = [
  { href: '/', label: 'Dashboard', icon: 'dashboard' },
  { href: '/booking-links', label: 'Booking Links', icon: 'link' },
  { href: '/polls', label: 'Polls', icon: 'poll' },
  { href: '/settings/account', label: 'Account', icon: 'user' },
  { href: '/settings', label: 'Calendar', icon: 'calendar' }
] as const;
```

Then in the template where icons are rendered (around line 105-139), add the `user` icon case. After the `{:else if item.icon === 'poll'}` block and before `{:else if item.icon === 'calendar'}`, add:

```svelte
{:else if item.icon === 'user'}
  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
  </svg>
```

Also update the `isActive` function with the same fix as the Sidebar:

```typescript
function isActive(href: string): boolean {
  if (href === '/') {
    return page.url.pathname === '/';
  }
  if (href === '/settings') {
    return page.url.pathname === '/settings';
  }
  return page.url.pathname.startsWith(href);
}
```

---

## Step 3: Create the settings layout with tab sub-navigation

Create `frontend/src/routes/(dashboard)/settings/+layout.svelte`:

```svelte
<script lang="ts">
  import { page } from '$app/state'

  let { children } = $props()

  const tabs = [
    { href: '/settings/account', label: 'Account' },
    { href: '/settings', label: 'Calendar' }
  ]

  function isActive(href: string): boolean {
    if (href === '/settings') {
      return page.url.pathname === '/settings'
    }
    return page.url.pathname.startsWith(href)
  }
</script>

<div class="settings-tabs">
  {#each tabs as tab}
    <a href={tab.href} class="settings-tab" class:active={isActive(tab.href)}>
      {tab.label}
    </a>
  {/each}
</div>

{@render children()}

<style>
  .settings-tabs {
    display: flex;
    gap: 0.25rem;
    margin-bottom: 1.5rem;
    border-bottom: 2px solid var(--border-color);
    padding-bottom: 0;
  }

  .settings-tab {
    padding: 0.6rem 1.25rem;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--text-secondary);
    text-decoration: none;
    border-bottom: 2px solid transparent;
    margin-bottom: -2px;
    transition: all var(--transition);
  }

  .settings-tab:hover {
    color: var(--text-primary);
  }

  .settings-tab.active {
    color: var(--text-primary);
    border-bottom-color: var(--cyan);
  }
</style>
```

This layout renders tab navigation above the page content. Both `/settings` and `/settings/account` will show these tabs, making it easy to switch between the two settings pages.

---

## Step 4: Verify the build compiles

Run:

```bash
cd frontend && pnpm check
```

Expected: No type errors.

---

## Step 5: Manual test

Run: `cd frontend && pnpm dev`

Verify:
1. The sidebar now shows a "Settings" section label with "Account" and "Calendar" items below it
2. Clicking "Account" navigates to `/settings/account` (will show a 404 or empty page for now -- this is expected, the page is created in meet-mesh-us05)
3. Clicking "Calendar" navigates to `/settings` and shows the existing calendar settings page
4. Tab navigation appears above the settings page content with "Account" and "Calendar" tabs
5. The correct tab is highlighted based on the current URL
6. Mobile menu shows both "Account" and "Calendar" nav items

---

## Step 6: Commit

```bash
git add frontend/src/lib/components/dashboard/Sidebar.svelte \
       frontend/src/lib/components/dashboard/MobileMenu.svelte \
       frontend/src/routes/\(dashboard\)/settings/+layout.svelte
git commit -m "feat(frontend): restructure settings navigation with account and calendar sub-items"
```

## Summary of Changes

- Updated Sidebar.svelte:
  - Added 'Account' nav item with user icon
  - Updated isActive function to handle /settings vs /settings/account
  - Added 'Settings' section label
  - Added nav-section-label CSS

- Updated MobileMenu.svelte:
  - Added 'Account' nav item with user icon
  - Updated isActive function

- Created settings/+layout.svelte:
  - Tab navigation with Account and Calendar tabs
  - Active state based on current URL
  - Shared across all /settings/* pages

Build passes with pnpm check (0 errors, 4 pre-existing warnings)
