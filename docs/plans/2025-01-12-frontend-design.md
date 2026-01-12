# Meet Mesh Frontend Design

## Overview

SvelteKit SPA (adapter-static) embedded in Go binary. Svelte 5 with runes for state management, Tailwind CSS for styling, Bits UI for accessible primitives, openapi-fetch for typed API calls.

**Visual direction:** Warm and friendly but professional. Soft rounded corners, warm neutral palette, subtle shadows, smooth transitions.

---

## Tech Stack

- **Framework:** SvelteKit with adapter-static
- **UI Library:** Svelte 5 (runes: `$state`, `$derived`, `$effect`)
- **Styling:** Tailwind CSS v3
- **Components:** Bits UI (headless primitives)
- **API Client:** openapi-fetch with generated types from openapi.yaml
- **Package Manager:** pnpm

---

## SPA Configuration

Since the frontend is embedded in a Go binary and served statically, we use pure SPA mode with no server-side rendering.

**`svelte.config.js`:**
```js
import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      fallback: '200.html'  // SPA fallback for client-side routing
    }),
    alias: {
      $components: 'src/lib/components'
    }
  }
};

export default config;
```

**`src/routes/+layout.js`:**
```js
// Disable SSR globally - pure client-side SPA
export const ssr = false;

// Disable prerendering since Go backend serves all routes
export const prerender = false;
```

The Go backend serves `200.html` for any route that doesn't match a static file, enabling client-side routing.

---

## Project Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts          # openapi-fetch instance
│   │   │   ├── types.ts           # Re-exports from generated types
│   │   │   └── generated/         # Output from openapi-typescript
│   │   │       └── schema.d.ts
│   │   ├── components/
│   │   │   ├── ui/                # Button, Card, Input, Select, Dialog, etc.
│   │   │   ├── booking/           # DateCalendar, TimeSlotList, BookingForm
│   │   │   ├── poll/              # VoteCard, VoteButtons, VoteSummary
│   │   │   ├── dashboard/         # Sidebar, LinkCard, BookingRow, StatsCard
│   │   │   └── shared/            # CustomFieldForm, Spinner, EmptyState, Toast
│   │   ├── stores/
│   │   │   └── auth.svelte.ts     # Auth state using Svelte 5 runes
│   │   └── utils/
│   │       ├── dates.ts           # Date formatting helpers
│   │       └── transitions.ts     # Shared transition configs
│   ├── routes/
│   │   ├── +layout.svelte
│   │   ├── (dashboard)/           # Organizer authenticated area
│   │   ├── (public)/              # Guest-facing pages
│   │   ├── (actions)/             # Email action handlers
│   │   └── auth/                  # Login/callback
│   └── app.css                    # Tailwind + CSS variables
├── static/
├── svelte.config.js
├── tailwind.config.js
├── vite.config.ts
└── package.json
```

---

## Routing

```
src/routes/
├── +layout.svelte                    # Root: global CSS, minimal wrapper
│
├── (dashboard)/                      # Organizer area (auth required)
│   ├── +layout.svelte                # Sidebar + auth guard
│   ├── +page.svelte                  # / - dashboard overview
│   ├── links/
│   │   ├── +page.svelte              # /links - list all links
│   │   ├── new/+page.svelte          # /links/new - create wizard
│   │   └── [id]/
│   │       ├── +page.svelte          # /links/:id - detail view
│   │       └── edit/+page.svelte     # /links/:id/edit
│   └── settings/+page.svelte         # /settings - calendar connections
│
├── (public)/                         # Guest pages (no auth)
│   ├── +layout.svelte                # Minimal centered layout
│   └── p/[slug]/
│       ├── +page.svelte              # /p/:slug - booking or poll
│       └── confirmed/+page.svelte    # /p/:slug/confirmed
│
├── (actions)/actions/                # Email action pages
│   ├── approve/+page.svelte          # /actions/approve?token=xxx
│   └── decline/+page.svelte          # /actions/decline?token=xxx
│
└── auth/
    ├── login/+page.svelte            # Redirects to OIDC
    └── callback/+page.svelte         # Handles OIDC callback
```

---

## Components

### UI Primitives (`src/lib/components/ui/`)

Built on Bits UI, styled with Tailwind:

| Component | Purpose |
|-----------|---------|
| Button | Primary, secondary, ghost, danger variants |
| Card | Elevated container with optional header/footer |
| Input | Text input with label, error state |
| Select | Dropdown using Bits UI Select |
| Dialog | Modal wrapper around Bits UI Dialog |
| Popover | Tooltips and popovers |
| Badge | Status badges (pending, confirmed, active, etc.) |
| Spinner | Loading indicator |
| EmptyState | Icon + message + optional action |
| Toast | Notifications (success, error, info) |

**Button variants:**
- Primary: Warm accent color, white text, subtle shadow
- Secondary: White background, gray border, dark text
- Ghost: No background, hover shows subtle fill
- Danger: Soft red for destructive actions

### Booking Components (`src/lib/components/booking/`)

| Component | Purpose |
|-----------|---------|
| BookingPage | Orchestrates full booking flow |
| DateCalendar | Month view, highlights available days |
| TimeSlotList | List of times for selected day |
| DayPicker | For full-day slot type |
| DateRangePicker | For multi-day slot type |
| BookingForm | Email + custom fields + submit |

**Booking flow (time slots):**
1. Guest sees calendar with dots on available days
2. Clicks day → time list slides in from right
3. Picks time → transitions to booking form
4. Submits → confirmation page

### Poll Components (`src/lib/components/poll/`)

| Component | Purpose |
|-----------|---------|
| PollPage | Orchestrates poll voting flow |
| VoteCard | Single option card with vote buttons |
| VoteButtons | Yes/No/Maybe toggle group |
| VoteSummary | Horizontal bar chart of results |
| VoterForm | Optional name/email before submit |

**VoteCard layout:**
```
┌─────────────────────────────────────────┐
│  📅  Monday, January 20                 │
│      2:00 PM - 3:00 PM                  │
│                                         │
│   [ Yes ]  [ No ]  [ Maybe ]            │
└─────────────────────────────────────────┘
```

Vote button states:
- Yes selected: Filled accent color
- No selected: Filled soft red
- Maybe selected: Filled amber
- Unselected: Ghost/outlined

### Dashboard Components (`src/lib/components/dashboard/`)

| Component | Purpose |
|-----------|---------|
| Sidebar | Persistent left nav (~240px) |
| DashboardHeader | Page title + actions |
| LinkCard | Link summary in grid view |
| BookingRow | Single booking in detail view |
| VoteRow | Single vote in poll detail |
| StatsCard | Metric card (counts, etc.) |
| CalendarConnectionCard | CalDAV setup card |

**Sidebar structure:**
```
┌──────────────────┐
│  Meet Mesh       │
├──────────────────┤
│  📊 Dashboard    │
│  🔗 Links        │
│  ⚙️ Settings     │
├──────────────────┤
│     (spacer)     │
├──────────────────┤
│  user@email      │
│  Logout          │
└──────────────────┘
```

---

## API Client

**Setup (`src/lib/api/client.ts`):**
```ts
import createClient from 'openapi-fetch'
import type { paths } from './generated/schema'

export const api = createClient<paths>({
  baseUrl: '/api',
})
```

**Build scripts (`package.json`):**
```json
{
  "scripts": {
    "generate:api": "openapi-typescript ../api/openapi.yaml -o src/lib/api/generated/schema.d.ts",
    "dev": "pnpm generate:api && vite dev",
    "build": "pnpm generate:api && vite build"
  }
}
```

**Usage:**
```ts
const { data, error } = await api.GET('/links/{id}', {
  params: { path: { id } }
})
if (error) {
  toast.error('Failed to load')
  return
}
// data is fully typed
```

---

## Auth State

**`src/lib/stores/auth.svelte.ts`:**
```ts
import { api } from '$lib/api/client'

type User = { id: number; email: string; name: string }

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

**Dashboard layout guard (`src/routes/(dashboard)/+layout.svelte`):**
```svelte
<script lang="ts">
  import { getAuth } from '$lib/stores/auth.svelte'
  import { goto } from '$app/navigation'
  import Sidebar from '$lib/components/dashboard/Sidebar.svelte'
  import LoadingScreen from '$lib/components/shared/LoadingScreen.svelte'

  let { children } = $props()
  const auth = getAuth()

  $effect(() => { auth.check() })
  $effect(() => {
    if (!auth.loading && !auth.isAuthenticated) {
      goto('/auth/login')
    }
  })
</script>

{#if auth.loading}
  <LoadingScreen />
{:else if auth.isAuthenticated}
  <div class="flex min-h-screen">
    <Sidebar />
    <main class="flex-1 p-6">
      {@render children()}
    </main>
  </div>
{/if}
```

---

## Design Tokens

**`src/app.css`:**
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --color-accent: 99 102 241;      /* Indigo - primary actions */
  --color-accent-hover: 79 70 229;
  --color-surface: 255 255 255;
  --color-muted: 248 250 252;
  --color-border: 226 232 240;

  --color-success: 34 197 94;
  --color-warning: 251 191 36;
  --color-danger: 239 68 68;

  --radius-sm: 0.375rem;
  --radius-md: 0.5rem;
  --radius-lg: 0.75rem;

  --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1);
}
```

---

## Transitions & Polish

**Page transitions:**
```svelte
<script>
  import { page } from '$app/state'
  import { fade } from 'svelte/transition'

  let { children } = $props()
</script>

{#key page.url.pathname}
  <div in:fade={{ duration: 150, delay: 75 }} out:fade={{ duration: 75 }}>
    {@render children()}
  </div>
{/key}
```

**Component transitions:**
- `fly` for sliding panels (150-200ms)
- `fade` for overlays and modals
- `scale` subtle on button press

**Loading states:**
- Skeleton placeholders with shimmer animation
- Button loading: spinner replaces text
- Full page: centered spinner with backdrop

**Empty states:**
- Icon + clear message + CTA button

**Toast notifications:**
- Bottom-right stack
- Auto-dismiss 4s
- Slide in from right

**Form validation:**
- Inline errors below fields
- Shake animation on invalid submit
- Focus first invalid field

---

## Responsive Breakpoints

| Breakpoint | Sidebar behavior |
|------------|------------------|
| Mobile (< 640px) | Hidden, hamburger menu |
| Tablet (640-1024px) | Collapsed to icons |
| Desktop (> 1024px) | Full sidebar |

---

## Implementation Order

1. Project scaffolding (SvelteKit, Tailwind, Bits UI, openapi-fetch)
2. Design tokens and base UI components
3. API client setup and type generation
4. Auth flow (login, callback, auth guard)
5. Dashboard layout with sidebar
6. Public layout (minimal, centered)
7. Booking flow components
8. Poll voting components
9. Dashboard pages (overview, links list, link detail)
10. Settings page (calendar connections)
11. Email action pages
12. Loading states, empty states, toasts
13. Responsive polish
