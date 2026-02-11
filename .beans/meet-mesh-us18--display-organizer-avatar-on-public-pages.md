---
# meet-mesh-us18
title: Display organizer avatar on public booking/poll pages
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T20:13:51Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us17
---

# Display Organizer Avatar on Public Booking/Poll Pages

**Goal:** Show the organizer's name and avatar on the public booking and poll pages, giving guests a sense of who they are scheduling with.

**Architecture:** The public booking and poll API responses now include `organizer_name` and `organizer_avatar_url` (from meet-mesh-us17). Update the frontend pages to display an organizer identity section at the top of the booking/poll page, with the avatar (or initials fallback) and name.

---

## Files

- Modify: `frontend/src/routes/(public)/p/booking/[slug]/+page.svelte`
- Modify: `frontend/src/routes/(public)/p/poll/[slug]/+page.svelte`

---

## Step 1: Update the booking page

Open `frontend/src/routes/(public)/p/booking/[slug]/+page.svelte`.

The page fetches data from `api.GET('/p/booking/{slug}', ...)` and stores it in `bookingLink`. After meet-mesh-us17, the response includes `organizer_name` and `organizer_avatar_url`.

Update the `PublicBookingLink` interface (around line 11) to include the new fields:

```typescript
interface PublicBookingLink {
    name: string;
    description?: string;
    custom_fields?: CustomField[];
    require_email?: boolean;
    slot_durations_minutes?: number[];
    organizer_name?: string;
    organizer_avatar_url?: string;
}
```

Add an organizer identity section to the page template, above the booking link name. Add this right after the loading state block, before the main content:

```svelte
{#if bookingLink?.organizer_name}
  <div class="organizer-identity">
    <div class="organizer-avatar">
      {#if bookingLink.organizer_avatar_url}
        <img src={bookingLink.organizer_avatar_url} alt={bookingLink.organizer_name} />
      {:else}
        <span>{bookingLink.organizer_name.charAt(0).toUpperCase()}</span>
      {/if}
    </div>
    <span class="organizer-name">{bookingLink.organizer_name}</span>
  </div>
{/if}
```

Add the CSS:

```css
.organizer-identity {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.organizer-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--orange, #f97316);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.organizer-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.organizer-avatar span {
  color: white;
  font-weight: 700;
  font-size: 0.8rem;
}

.organizer-name {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text-secondary, #6b7280);
}
```

---

## Step 2: Update the poll page

Open `frontend/src/routes/(public)/p/poll/[slug]/+page.svelte`.

Apply the same changes: update the interface/type for the poll data to include `organizer_name` and `organizer_avatar_url`, and add the same organizer identity section to the template with the same CSS.

The exact interface name and structure may differ from the booking page. Look for the type that stores the poll data fetched from `api.GET('/p/poll/{slug}', ...)` and add the two new fields. Then add the same organizer identity HTML/CSS block in an appropriate location above the poll title.

---

## Step 3: Verify

Run:

```bash
cd frontend && pnpm check
```

Expected: No type errors.

Manual test:
1. Open a public booking page (e.g., `/p/booking/some-slug`)
2. If the organizer has a name and avatar set, they should appear at the top
3. If the organizer has a name but no avatar, the first letter should appear as a circle
4. If no organizer name is set, the section should not appear
5. Repeat for a public poll page

---

## Step 4: Commit

```bash
git add frontend/src/routes/\(public\)/p/booking/\[slug\]/+page.svelte frontend/src/routes/\(public\)/p/poll/\[slug\]/+page.svelte
git commit -m "feat(frontend): display organizer avatar and name on public booking/poll pages"
```

## Summary of Changes

- Updated `PublicBookingLink` interface in booking page route to include `organizer_name` and `organizer_avatar_url`
- Updated `PublicPoll` interface in poll page route to include `organizer_name` and `organizer_avatar_url`
- Updated `PublicLink` interface in BookingPage.svelte component to accept organizer fields
- Updated `PublicLink` interface in PollPage.svelte component to accept organizer fields
- Added organizer identity section to BookingPage.svelte header (avatar with initials fallback + name)
- Added organizer identity section to PollPage.svelte header (avatar with initials fallback + name)
- Used Tailwind CSS classes for styling (consistent with project style)

Build verified: `pnpm check` passes with 0 errors
