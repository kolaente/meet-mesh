---
# meet-mesh-us16
title: Update sidebar and mobile menu to display avatar
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:39:28Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us09
---

# Update Sidebar and Mobile Menu to Display Avatar

**Goal:** Update the sidebar and mobile menu to show the organizer's avatar image when one is set, falling back to the existing initials display when no avatar exists.

**Architecture:** Both `Sidebar.svelte` and `MobileMenu.svelte` already read from the auth store. Since meet-mesh-us09 adds `avatar_url` to the auth store's User type, we just need to conditionally render an `<img>` tag when `auth.user?.avatar_url` is set, otherwise show the initials as before.

---

## Files

- Modify: `frontend/src/lib/components/dashboard/Sidebar.svelte`
- Modify: `frontend/src/lib/components/dashboard/MobileMenu.svelte`

---

## Step 1: Update Sidebar.svelte

Open `frontend/src/lib/components/dashboard/Sidebar.svelte`.

Find the user avatar section (around line 107):

```svelte
<div class="user-avatar">{getUserInitials()}</div>
```

Replace it with:

```svelte
<div class="user-avatar">
  {#if auth.user?.avatar_url}
    <img src={auth.user.avatar_url} alt="Avatar" class="user-avatar-img" />
  {:else}
    {getUserInitials()}
  {/if}
</div>
```

Add the CSS for the avatar image in the `<style>` block:

```css
.user-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: var(--radius);
}
```

---

## Step 2: Update MobileMenu.svelte

Open `frontend/src/lib/components/dashboard/MobileMenu.svelte`.

Find the user avatar section (around line 148):

```svelte
<div class="user-avatar">
    {auth.user?.email?.charAt(0).toUpperCase() ?? 'U'}
</div>
```

Replace it with:

```svelte
<div class="user-avatar">
  {#if auth.user?.avatar_url}
    <img src={auth.user.avatar_url} alt="Avatar" class="user-avatar-img" />
  {:else}
    {auth.user?.email?.charAt(0).toUpperCase() ?? 'U'}
  {/if}
</div>
```

Add the same CSS for the avatar image in the MobileMenu's `<style>` block:

```css
.user-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: var(--radius);
}
```

---

## Step 3: Verify

Run:

```bash
cd frontend && pnpm check
```

Expected: No type errors.

Manual test:
1. Log in as the organizer
2. If an avatar is set, it should appear in the sidebar user card (bottom-left) and in the mobile menu
3. If no avatar is set, initials should appear as before
4. Upload an avatar in `/settings/account` -- the sidebar should update immediately (if the auth store is reactive)

---

## Step 4: Commit

```bash
git add frontend/src/lib/components/dashboard/Sidebar.svelte frontend/src/lib/components/dashboard/MobileMenu.svelte
git commit -m "feat(frontend): display avatar in sidebar and mobile menu with initials fallback"
```

## Summary of Changes

- Updated Sidebar.svelte to display avatar_url when set, fallback to initials
- Updated MobileMenu.svelte to display avatar_url when set, fallback to initials
- Added CSS for user-avatar-img in both components

pnpm check passes with 0 errors
