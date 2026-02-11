---
# meet-mesh-us07
title: Profile Image Upload
status: completed
type: epic
priority: normal
created_at: 2026-02-11T20:06:00Z
updated_at: 2026-02-11T20:14:18Z
---

# Profile Image Upload Implementation Plan

**Goal:** Allow the organizer to upload a profile image/avatar that is displayed in the sidebar, on public booking pages, and in email notifications.

**Architecture:** Add avatar upload/delete as plain HTTP handlers outside of ogen (since ogen does not support multipart file uploads natively). Store images on the local filesystem at a configurable path (`storage.avatars_path`). Serve avatars at `/api/avatars/{filename}` via a separate mux route. Add `AvatarFilename` to the User GORM model and `avatar_url` to the OpenAPI User schema. Frontend gets a drag-and-drop upload component on the account settings page, with avatar display in sidebar, mobile menu, public pages, and email templates. Filenames use content hashes for cache busting.

**Tech Stack:** Go 1.25 (stdlib `image`, `golang.org/x/image/draw` for resizing), ogen (OpenAPI codegen for all non-upload endpoints), SvelteKit with Svelte 5 runes, Tailwind CSS v4, openapi-fetch

---

## Dependencies

- Requires meet-mesh-us01 (User Settings Page) to be completed first, since the upload UI lives on the account settings page.

## Child Tasks

1. **meet-mesh-us08** - Add storage config and avatar directory initialization
2. **meet-mesh-us09** - Add AvatarFilename to User model and avatar_url to OpenAPI User schema
3. **meet-mesh-us10** - Implement avatar upload HTTP handler (POST /api/avatars)
4. **meet-mesh-us11** - Implement avatar delete HTTP handler (DELETE /api/avatars)
5. **meet-mesh-us12** - Implement avatar file serving (GET /api/avatars/{filename})
6. **meet-mesh-us13** - Register avatar routes in main.go mux
7. **meet-mesh-us14** - Create frontend AvatarUpload component
8. **meet-mesh-us15** - Integrate avatar into account settings page
9. **meet-mesh-us16** - Update sidebar and mobile menu to display avatar
10. **meet-mesh-us17** - Add organizer info to public booking/poll API responses
11. **meet-mesh-us18** - Display organizer avatar on public booking/poll pages
12. **meet-mesh-us19** - Add organizer avatar to email notification templates

## Key Files

| File | Role |
|------|------|
| `api/config.go` | Add StorageConfig with AvatarsPath |
| `config.example.yaml` | Add storage.avatars_path example |
| `api/models.go` | Add AvatarFilename to User struct |
| `api/openapi.yaml` | Add avatar_url to User schema, organizer fields to public responses |
| `api/handler_avatar.go` | New file: upload, delete, serve handlers |
| `api/cmd/main.go` | Register avatar mux routes |
| `frontend/src/lib/components/ui/AvatarUpload.svelte` | New file: drag-and-drop upload component |
| `frontend/src/lib/components/dashboard/Sidebar.svelte` | Show avatar image with initials fallback |
| `frontend/src/lib/components/dashboard/MobileMenu.svelte` | Show avatar image with initials fallback |
| `frontend/src/lib/stores/auth.svelte.ts` | Add avatar_url to User type |
| `frontend/src/routes/(dashboard)/settings/account/+page.svelte` | Add avatar upload section |
| `frontend/src/routes/(public)/p/booking/[slug]/+page.svelte` | Show organizer avatar |
| `frontend/src/routes/(public)/p/poll/[slug]/+page.svelte` | Show organizer avatar |
| `api/mailer.go` | Add avatar to email templates |

## Summary

All 12 child tasks have been completed:

1. **meet-mesh-us08** - Storage config and avatar directory initialization
2. **meet-mesh-us09** - AvatarFilename in User model and avatar_url in OpenAPI
3. **meet-mesh-us10** - Avatar upload HTTP handler (POST /api/avatars)
4. **meet-mesh-us11** - Avatar delete HTTP handler (DELETE /api/avatars)
5. **meet-mesh-us12** - Avatar file serving (GET /api/avatars/{filename})
6. **meet-mesh-us13** - Avatar routes registered in main.go mux
7. **meet-mesh-us14** - Frontend AvatarUpload component
8. **meet-mesh-us15** - Avatar upload integrated into account settings
9. **meet-mesh-us16** - Avatar displayed in sidebar and mobile menu
10. **meet-mesh-us17** - Organizer info added to public API responses
11. **meet-mesh-us18** - Organizer avatar displayed on public booking/poll pages
12. **meet-mesh-us19** - Organizer avatar added to email notification templates

The profile image upload feature is now fully implemented!
