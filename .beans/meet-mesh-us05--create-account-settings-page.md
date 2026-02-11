---
# meet-mesh-us05
title: Create account settings page with display name and date/time preferences
status: completed
type: task
priority: normal
created_at: 2026-02-11T20:04:00Z
updated_at: 2026-02-11T19:22:27Z
parent: meet-mesh-us01
blocked_by:
    - meet-mesh-us02
    - meet-mesh-us03
    - meet-mesh-us04
---

# Create Account Settings Page

**Goal:** Create the `/settings/account` page where the organizer can edit their display name and configure date/time format preferences. The display name is persisted via the `PUT /auth/me` API endpoint. The date/time preferences continue to use localStorage (same as before, just moved to this page).

**Architecture:** Create a new SvelteKit page at `frontend/src/routes/(dashboard)/settings/account/+page.svelte`. The page has two sections:
1. **Profile** section with a display name input field and save button (calls `PUT /auth/me`)
2. **Date & Time Format** section (moved from the current calendar settings page) with time format and week start day selects

The page follows the same patterns as the existing settings page: uses `DashboardHeader`, `Card`, `Input`, `Select`, `Button`, and `Spinner` components from the UI library.

**Prerequisites:** meet-mesh-us02 (API spec), meet-mesh-us03 (handler), meet-mesh-us04 (navigation + layout)

---

## Files

- Create: `frontend/src/routes/(dashboard)/settings/account/+page.svelte`

---

## Step 1: Create the account settings page

Create `frontend/src/routes/(dashboard)/settings/account/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { getAuth } from '$lib/stores/auth.svelte';
  import { getDateFormat } from '$lib/stores/dateFormat.svelte';
  import { DashboardHeader } from '$lib/components/dashboard';
  import { Card, Button, Input, Select, Spinner } from '$lib/components/ui';
  import { addToast } from '$lib/stores/toast.svelte';

  const auth = getAuth();
  const dateFormat = getDateFormat();

  let displayName = $state('');
  let saving = $state(false);
  let loading = $state(true);

  // Reactive bindings for date format selects
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

  onMount(async () => {
    // Load current user data
    const { data } = await api.GET('/auth/me');
    if (data) {
      displayName = data.name ?? '';
    }
    loading = false;
  });

  async function saveProfile() {
    saving = true;
    try {
      const { data, error } = await api.PUT('/auth/me', {
        body: { name: displayName }
      });

      if (error) {
        addToast({ type: 'error', message: 'Failed to save profile' });
        return;
      }

      if (data) {
        // Refresh the auth store so sidebar/header reflects the new name
        auth.refreshUser(data);
        addToast({ type: 'success', message: 'Profile updated' });
      }
    } finally {
      saving = false;
    }
  }
</script>

<DashboardHeader title="Account Settings" />

{#if loading}
  <div class="flex items-center justify-center py-12">
    <Spinner size="lg" />
  </div>
{:else}
  <div class="space-y-6">
    <!-- Profile Section -->
    <section>
      <div class="mb-4">
        <h2 class="text-lg font-medium text-[var(--text-primary)]">Profile</h2>
        <p class="text-sm text-[var(--text-secondary)]">Manage your display name and profile information</p>
      </div>

      <Card>
        {#snippet children()}
          <div class="space-y-6">
            <div class="flex flex-col sm:flex-row sm:items-start gap-4">
              <div class="flex-1">
                <Input
                  name="displayName"
                  label="Display Name"
                  description="This name is shown to guests on your booking pages and polls."
                  bind:value={displayName}
                  placeholder="Enter your name"
                />
              </div>
            </div>

            <div>
              <Input
                name="email"
                label="Email"
                description="Your email address is managed by your identity provider and cannot be changed here."
                value={auth.user?.email ?? ''}
                disabled={true}
              />
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
    </section>

    <!-- Date & Time Format Section (moved from calendar settings) -->
    <section>
      <div class="mb-4">
        <h2 class="text-lg font-medium text-[var(--text-primary)]">Date & Time Format</h2>
        <p class="text-sm text-[var(--text-secondary)]">Customize how dates and times are displayed</p>
      </div>

      <Card>
        {#snippet children()}
          <div class="space-y-6">
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              <div>
                <p class="font-medium text-[var(--text-primary)]">Time Format</p>
                <p class="text-sm text-[var(--text-secondary)]">Choose 12-hour or 24-hour time display</p>
              </div>
              <div class="w-full sm:w-48">
                <Select
                  name="timeFormat"
                  options={timeFormatOptions}
                  value={timeFormatValue}
                  onchange={(value) => dateFormat.setTimeFormat(value as '12h' | '24h')}
                />
              </div>
            </div>

            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              <div>
                <p class="font-medium text-[var(--text-primary)]">Week Starts On</p>
                <p class="text-sm text-[var(--text-secondary)]">Choose which day starts your week</p>
              </div>
              <div class="w-full sm:w-48">
                <Select
                  name="weekStartDay"
                  options={weekStartOptions}
                  value={weekStartDayValue}
                  onchange={(value) => dateFormat.setWeekStartDay(value as 'sunday' | 'monday')}
                />
              </div>
            </div>

            <div class="pt-2 border-t border-[var(--border-color)]">
              <Button variant="secondary" onclick={() => dateFormat.reset()}>
                {#snippet children()}Reset to Browser Defaults{/snippet}
              </Button>
            </div>
          </div>
        {/snippet}
      </Card>
    </section>
  </div>
{/if}
```

**Notes about the code above:**

- The `addToast` import assumes a toast store exists at `$lib/stores/toast.svelte`. Check the actual import path -- the project has `ToastContainer.svelte` and likely a `toast.svelte.ts` store. Adjust the import if needed.
- The `auth.refreshUser(data)` method does not exist yet -- it is added in meet-mesh-us06.
- The date/time format section is an exact copy of what is currently in the calendar settings page. After this page is created, the next task (meet-mesh-us06) removes it from the calendar settings.
- The `api.PUT('/auth/me', ...)` call will be type-safe because the OpenAPI types were regenerated in meet-mesh-us02.

---

## Step 2: Verify the build compiles

Run:

```bash
cd frontend && pnpm check
```

Expected: No type errors (there may be a warning about `auth.refreshUser` not existing yet -- this is addressed in meet-mesh-us06).

If `pnpm check` fails because `auth.refreshUser` does not exist yet, temporarily comment out the `auth.refreshUser(data)` line and add a `// TODO: uncomment after meet-mesh-us06` comment. It will be uncommented in the next task.

---

## Step 3: Manual test

Run: `cd frontend && pnpm dev`

Navigate to `/settings/account`.

Verify:
1. The page shows "Account Settings" header
2. The "Account" tab is highlighted in the tab navigation
3. The Profile section shows a display name input field (pre-filled from the API) and a disabled email field
4. Typing in the display name field and clicking "Save Profile" calls `PUT /auth/me`
5. Success/error toast appears
6. The Date & Time Format section shows time format and week start day selects
7. Changing date/time settings works immediately (persists to localStorage)
8. "Reset to Browser Defaults" button works

---

## Step 4: Commit

```bash
git add frontend/src/routes/\(dashboard\)/settings/account/+page.svelte
git commit -m "feat(frontend): create account settings page with display name and date/time preferences"
```

## Summary of Changes

Created frontend/src/routes/(dashboard)/settings/account/+page.svelte:
- Profile section with display name input and save functionality
- Email field (disabled, read-only from OIDC)
- Calls PUT /auth/me to persist name changes
- Toast notifications for success/error
- Date & Time Format section (moved from calendar settings)
- Time format (12h/24h) and week start day selects
- Reset to Browser Defaults button
- Uses existing UI components: DashboardHeader, Card, Input, Select, Button, Spinner

Build passes with pnpm check (0 errors)
