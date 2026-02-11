---
# meet-mesh-99dd
title: Add dedicated logout button to sidebar
status: completed
type: task
priority: normal
created_at: 2026-02-11T19:20:00Z
updated_at: 2026-02-11T18:54:17Z
---

# Sidebar Logout UX Improvement - Implementation Plan

**Goal:** Replace the current click-on-email-to-logout pattern with a proper dropdown menu on the user card, featuring an explicit "Log out" menu item.

**Architecture:** Use bits-ui DropdownMenu (already a project dependency at ^2.15.4) to wrap the existing user card in the sidebar footer. The user card becomes the dropdown trigger. The dropdown content contains a "Log out" item. The same pattern already exists correctly in MobileMenu.svelte (separate logout button) -- this brings the desktop Sidebar into parity.

**Tech Stack:** Svelte 5 (runes), bits-ui DropdownMenu, existing CSS variables

---

## Current State (the problem)

In `frontend/src/lib/components/dashboard/Sidebar.svelte` (lines 105-113), the entire user card is a `<button>` with `onclick={handleLogout}` and `title="Click to logout"`. This is bad UX because:
1. There is no visual indication that clicking the user info will log you out
2. Users expect clicking on their profile to show account options, not immediately log out
3. Accidental clicks cause unexpected logout
4. The mobile menu (`MobileMenu.svelte`) already has a separate, explicit logout button -- the desktop sidebar should match

## Desired State

- The user card area (avatar + name + email) should be the **trigger** for a dropdown menu
- The dropdown opens **upward** (side="top") since the user card is at the bottom of the sidebar
- The dropdown contains a single "Log out" item (with a logout icon)
- The user card itself no longer triggers logout directly
- Styling follows the project's neo-brutalist design system (CSS variables like `--bg-secondary`, `--border`, `--radius`, `--shadow-sm`, etc.)

---

### Task 1: Add DropdownMenu to Sidebar user card

**Files:**
- Modify: `frontend/src/lib/components/dashboard/Sidebar.svelte`

**Step 1: Add the bits-ui DropdownMenu import**

At the top of the `<script>` block, add:

```svelte
import { DropdownMenu } from 'bits-ui'
```

**Step 2: Replace the user card button with a DropdownMenu**

Replace the entire `<!-- User section -->` block (lines 104-113):

```svelte
<!-- User section -->
<div class="sidebar-footer">
  <DropdownMenu.Root>
    <DropdownMenu.Trigger>
      {#snippet child({ props })}
        <button {...props} class="user-card">
          <div class="user-avatar">{getUserInitials()}</div>
          <div class="user-info">
            <div class="user-name">{getUserDisplayName()}</div>
            <div class="user-email">{auth.user?.email ?? ''}</div>
          </div>
          <svg class="user-chevron" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
          </svg>
        </button>
      {/snippet}
    </DropdownMenu.Trigger>

    <DropdownMenu.Portal>
      <DropdownMenu.Content class="user-menu-content" side="top" sideOffset={8} align="start">
        <DropdownMenu.Item class="user-menu-item" onSelect={handleLogout}>
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          Log out
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Portal>
  </DropdownMenu.Root>
</div>
```

**Step 3: Add styles for the dropdown and chevron**

Add the following CSS rules inside the `<style>` block:

```css
.user-chevron {
  width: 16px;
  height: 16px;
  color: var(--text-muted);
  flex-shrink: 0;
  transition: color var(--transition);
}

.user-card:hover .user-chevron {
  color: var(--text-primary);
}

:global(.user-menu-content) {
  min-width: 200px;
  background: var(--bg-secondary);
  border: var(--border);
  border-radius: var(--radius);
  padding: 0.35rem;
  box-shadow: var(--shadow);
  z-index: 200;
}

:global(.user-menu-content[data-state="open"]) {
  animation: menuIn 0.15s ease;
}

:global(.user-menu-content[data-state="closed"]) {
  animation: menuOut 0.1s ease;
}

@keyframes menuIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes menuOut {
  from { opacity: 1; transform: translateY(0); }
  to { opacity: 0; transform: translateY(4px); }
}

:global(.user-menu-item) {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.6rem 0.75rem;
  border-radius: calc(var(--radius) - 2px);
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition);
  border: none;
  background: transparent;
  width: 100%;
}

:global(.user-menu-item:hover),
:global(.user-menu-item[data-highlighted]) {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

:global(.user-menu-item svg) {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}
```

**Step 4: Verify the build compiles**

Run: `cd frontend && pnpm check`
Expected: No type errors

**Step 5: Manual test**

Run: `cd frontend && pnpm dev`
Verify:
1. User card at bottom of sidebar shows avatar, name, email, and a chevron icon
2. Clicking the user card opens a dropdown menu ABOVE the card
3. The dropdown contains a "Log out" item with a logout icon
4. Clicking "Log out" triggers logout (redirects to /auth/login)
5. Clicking outside the dropdown closes it
6. Pressing Escape closes the dropdown
7. Keyboard navigation works (arrow keys, Enter to select)
8. Dark mode styling looks correct

**Step 6: Commit**

```bash
git add frontend/src/lib/components/dashboard/Sidebar.svelte
git commit -m "feat: replace sidebar email-click logout with dropdown menu"
```

---

## Summary of changes

Only one file is modified: `frontend/src/lib/components/dashboard/Sidebar.svelte`

The change:
1. Imports `DropdownMenu` from bits-ui
2. Replaces the `<button class="user-card" onclick={handleLogout}>` with a `DropdownMenu.Root` > `DropdownMenu.Trigger` > `DropdownMenu.Content` > `DropdownMenu.Item` structure
3. Adds a chevron icon to the user card to hint at the dropdown
4. Adds CSS for the dropdown content and items using the project's existing CSS variable design system
5. Uses `side="top"` on the content so the menu opens upward (since the user card is at the bottom of the sidebar)

No changes needed to:
- `auth.svelte.ts` (logout function stays the same)
- `MobileMenu.svelte` (already has a proper separate logout button)
- `openapi.yaml` or any backend files
- No new files or components to create

## Summary of Changes

- Replaced click-on-email-to-logout pattern with DropdownMenu from bits-ui
- User card now shows a chevron icon indicating interactive dropdown
- Dropdown opens upward (side="top") with "Log out" menu item
- Added CSS for dropdown animations and styling following project's design system
