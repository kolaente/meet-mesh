---
# meet-mesh-us15
title: Integrate avatar upload into account settings page
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:37:19Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us14
---

# Integrate Avatar Upload into Account Settings Page

**Goal:** Add the AvatarUpload component to the account settings page (`/settings/account`) so the organizer can upload, preview, and delete their profile avatar.

**Architecture:** Add an avatar section to the Profile card on the account settings page. The avatar is displayed alongside the display name field. When the avatar changes, the auth store is updated so the sidebar reflects the new avatar immediately.

**Prerequisite:** This task assumes meet-mesh-us05 (Create account settings page) is completed. The page exists at `frontend/src/routes/(dashboard)/settings/account/+page.svelte`.

---

## Files

- Modify: `frontend/src/routes/(dashboard)/settings/account/+page.svelte`

---

## Step 1: Add avatar section to the account settings page

Open `frontend/src/routes/(dashboard)/settings/account/+page.svelte`.

Add the import for `AvatarUpload`:

```typescript
import { Card, Button, Input, Select, Spinner, AvatarUpload } from '$lib/components/ui';
```

Add a state variable for the avatar URL (alongside existing state variables like `displayName`, `saving`, `loading`):

```typescript
let avatarUrl = $state('');
```

In the `onMount` callback where user data is loaded, also set the avatar URL:

```typescript
onMount(async () => {
  const { data } = await api.GET('/auth/me');
  if (data) {
    displayName = data.name ?? '';
    avatarUrl = data.avatar_url ?? '';
  }
  loading = false;
});
```

Add a helper function to get user initials (same logic as Sidebar.svelte):

```typescript
function getUserInitials(): string {
  const name = auth.user?.name || auth.user?.email || '';
  if (!name) return '?';
  if (name.includes('@')) {
    const local = name.split('@')[0];
    const parts = local.split(/[._-]/);
    if (parts.length >= 2) {
      return (parts[0][0] + parts[1][0]).toUpperCase();
    }
    return local.substring(0, 2).toUpperCase();
  }
  const parts = name.split(/\s+/);
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  return name.substring(0, 2).toUpperCase();
}
```

Add a handler for avatar changes:

```typescript
function handleAvatarChange(newUrl: string) {
  avatarUrl = newUrl;
  // Update auth store so sidebar reflects the change immediately
  if (auth.user) {
    auth.user = { ...auth.user, avatar_url: newUrl || undefined };
  }
}
```

Note: The auth store may not support direct mutation of `auth.user` depending on how it was implemented. If `auth.user` is read-only, you may need to add a method like `auth.updateAvatarUrl(url)` or `auth.refreshUser({...data, avatar_url: newUrl})` to the auth store. Check the auth store implementation (meet-mesh-us06 may have added a `refreshUser` method).

---

## Step 2: Add AvatarUpload to the Profile section template

In the Profile section (inside the Card), add the AvatarUpload component above or alongside the display name input. The layout should show the avatar on the left and the form fields on the right.

Update the Profile section to look like:

```svelte
<Card>
  {#snippet children()}
    <div class="space-y-6">
      <!-- Avatar -->
      <div class="flex flex-col sm:flex-row sm:items-start gap-6">
        <AvatarUpload
          avatarUrl={avatarUrl}
          initials={getUserInitials()}
          onchange={handleAvatarChange}
        />
        <div class="flex-1 space-y-4">
          <Input
            name="displayName"
            label="Display Name"
            description="This name is shown to guests on your booking pages and polls."
            bind:value={displayName}
            placeholder="Enter your name"
          />
          <Input
            name="email"
            label="Email"
            description="Your email address is managed by your identity provider and cannot be changed here."
            value={auth.user?.email ?? ''}
            disabled={true}
          />
        </div>
      </div>

      <div class="flex justify-end pt-2 border-t border-[var(--border-color)]">
        <Button variant="primary" onclick={saveProfile} disabled={saving}>
          {#snippet children()}
            {#if saving}
              <Spinner size="sm" />
              Saving...
            {:else}
              Save Profile
            {/if}
          {/snippet}
        </Button>
      </div>
    </div>
  {/snippet}
</Card>
```

---

## Step 3: Verify

Run:

```bash
cd frontend && pnpm check
```

Expected: No type errors.

Manual test:
1. Navigate to `/settings/account`
2. The avatar upload circle should appear with initials
3. Click to open file picker, or drag a file
4. After upload, the circle should show the uploaded image
5. The "Remove avatar" link should appear
6. Clicking "Remove avatar" should revert to initials

---

## Step 4: Commit

```bash
git add frontend/src/routes/\(dashboard\)/settings/account/+page.svelte
git commit -m "feat(frontend): add avatar upload to account settings page"
```

## Summary of Changes

- Added AvatarUpload import
- Added avatarUrl state variable
- Added getUserInitials() helper function
- Added handleAvatarChange() to update auth store on avatar change
- Updated onMount to load avatar_url
- Integrated AvatarUpload component into Profile section layout

pnpm check passes with 0 errors
