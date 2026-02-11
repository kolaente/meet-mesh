---
# meet-mesh-us01
title: User Settings Page
status: completed
type: epic
priority: normal
created_at: 2026-02-11T20:00:00Z
updated_at: 2026-02-11T19:25:46Z
---

# User Settings Page Implementation Plan

**Goal:** Add a dedicated user settings page to the dashboard where the organizer can manage their profile (display name) and date/time preferences, separate from the existing calendar connection settings.

**Architecture:** Add a new `PUT /auth/me` API endpoint to the Go backend for updating user profile data. Restructure the frontend settings routes: current `/settings` becomes the calendar-focused settings page, new `/settings/account` handles user profile and date/time preferences. Update sidebar navigation to show a "Settings" parent with "Account" and "Calendar" sub-items.

**Tech Stack:** Go (ogen codegen), SvelteKit (Svelte 5 runes), openapi-fetch, Tailwind CSS v4

---

## Current State

- **User model** (`api/models.go`): Has `ID`, `OIDCSub`, `Email`, `Name`, `CreatedAt`. The `Name` field exists but is only populated from OIDC claims during initial login (`handler_auth.go` line 79). There is no way to update it after creation.
- **API User schema** (`api/openapi.yaml`): Has `id`, `email`, `name` (optional string). No update endpoint exists.
- **Auth store** (`frontend/src/lib/stores/auth.svelte.ts`): Stores `{ id, email, name? }`. No refresh/refetch mechanism after updates.
- **Current settings page** (`frontend/src/routes/(dashboard)/settings/+page.svelte`): Contains two sections:
  1. Calendar Connections (add/remove CalDAV connections)
  2. Date & Time Format (time format 12h/24h, week starts on Sunday/Monday) -- stored in localStorage via `dateFormat.svelte.ts`
- **Sidebar** (`frontend/src/lib/components/dashboard/Sidebar.svelte`): Has `settingsItems` array with one entry: `{ href: '/settings', label: 'Calendar', icon: 'calendar' }`.
- **MobileMenu** (`frontend/src/lib/components/dashboard/MobileMenu.svelte`): Same nav items hardcoded.

## Scope

### In scope
- New `PUT /auth/me` API endpoint for updating user name
- New `/settings/account` page with display name editing and date/time preferences
- Move date/time format section from calendar settings to account settings
- Restructure settings navigation (sidebar + mobile menu)
- Settings layout with sub-navigation tabs
- Auth store refresh mechanism after profile update

### Out of scope (separate epic)
- Profile image upload (see bean: meet-mesh-us07)

## Child Tasks

1. **meet-mesh-us02** - Add `PUT /auth/me` endpoint to OpenAPI spec and regenerate
2. **meet-mesh-us03** - Implement `UpdateCurrentUser` handler in Go backend
3. **meet-mesh-us04** - Restructure settings navigation in sidebar and mobile menu
4. **meet-mesh-us05** - Create account settings page with display name and date/time preferences
5. **meet-mesh-us06** - Remove date/time section from calendar settings page and add auth store refresh

## Key Files

| File | Role |
|------|------|
| `api/openapi.yaml` | API spec (source of truth) |
| `api/gen/` | Generated code (DO NOT EDIT, run `make generate`) |
| `api/handler_auth.go` | Auth handlers including GetCurrentUser |
| `api/models.go` | GORM data models (User already has Name field) |
| `frontend/src/routes/(dashboard)/settings/+page.svelte` | Current settings page (calendar + date/time) |
| `frontend/src/routes/(dashboard)/settings/` | Settings route directory |
| `frontend/src/lib/components/dashboard/Sidebar.svelte` | Desktop sidebar nav |
| `frontend/src/lib/components/dashboard/MobileMenu.svelte` | Mobile nav |
| `frontend/src/lib/stores/auth.svelte.ts` | Auth store with user data |
| `frontend/src/lib/stores/dateFormat.svelte.ts` | Date format preferences (localStorage) |
| `frontend/src/lib/api/client.ts` | openapi-fetch client |
| `frontend/src/lib/api/types.ts` | Re-exports generated OpenAPI types |

## Summary

All child tasks completed:
- meet-mesh-us02: Added PUT /auth/me endpoint to OpenAPI spec
- meet-mesh-us03: Implemented UpdateCurrentUser handler in Go backend
- meet-mesh-us04: Restructured settings navigation (sidebar, mobile menu, tab layout)
- meet-mesh-us05: Created account settings page with profile and date/time preferences
- meet-mesh-us06: Cleaned up calendar settings page, added auth refresh, updated sidebar to use display name

The user settings page epic is now complete. Users can:
- Edit their display name on /settings/account
- See their name reflected immediately in the sidebar
- Configure date/time format preferences
- Access calendar settings via /settings (tab navigation)
