---
# meet-mesh-us17
title: Add organizer info to public booking/poll API responses
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:41:15Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us09
---

# Add Organizer Info to Public Booking/Poll API Responses

**Goal:** Include the organizer's name and avatar URL in the public booking link and poll API responses, so the frontend can display the organizer's identity on public pages.

**Architecture:** Modify the OpenAPI spec to add `organizer_name` and `organizer_avatar_url` optional fields to the `getPublicBookingLink` and `getPublicPoll` response schemas. Then update the Go handlers to load the organizer (User) from the database via the link/poll's `UserID` and include the fields in the response.

---

## Files

- Modify: `api/openapi.yaml` (add organizer fields to public responses)
- Regenerate: `api/gen/` (run `make generate`)
- Regenerate: `frontend/src/lib/api/generated/` (run `pnpm generate:api`)
- Modify: `api/handler_public_booking.go` (include organizer info)
- Modify: `api/handler_public_poll.go` (include organizer info)

---

## Step 1: Update OpenAPI spec

Open `api/openapi.yaml`. Find the `getPublicBookingLink` response schema (around line 1044). Currently it is an inline object with `name`, `description`, `custom_fields`, `require_email`, `slot_durations_minutes`. Add two new optional fields:

```yaml
                  organizer_name:
                    type: string
                    description: Display name of the organizer
                  organizer_avatar_url:
                    type: string
                    description: URL to the organizer's avatar image
```

Do the same for the `getPublicPoll` response schema. Find it in the spec (search for `operationId: getPublicPoll`). Add the same two fields to its response object.

---

## Step 2: Regenerate code

```bash
make generate
cd frontend && pnpm generate:api
```

Verify: `grep -n "OrganizerName" api/gen/oas_schemas_gen.go` should show the new fields.

---

## Step 3: Update handler_public_booking.go

Open `api/handler_public_booking.go`. In `GetPublicBookingLink`, after fetching the `BookingLink`, also fetch the organizer:

```go
// Fetch organizer for public display
var organizer User
h.db.First(&organizer, link.UserID)
```

Then add the organizer fields to the response:

```go
return &gen.GetPublicBookingLinkOK{
    Name:                 link.Name,
    Description:          gen.NewOptString(link.Description),
    CustomFields:         mapCustomFieldsToGen(link.CustomFields),
    RequireEmail:         gen.NewOptBool(link.RequireEmail),
    SlotDurationsMinutes: durations,
    OrganizerName:        gen.NewOptString(organizer.Name),
    OrganizerAvatarURL:   gen.NewOptString(avatarURL(organizer.AvatarFilename)),
}, nil
```

Note: `avatarURL()` is the helper function added in meet-mesh-us09 in `handler_auth.go`. If it is not accessible from this file (it should be, since both are in package `api`), make sure it is exported or in a shared file.

---

## Step 4: Update handler_public_poll.go

Open `api/handler_public_poll.go`. In `GetPublicPoll`, after fetching the `Poll`, also fetch the organizer:

```go
// Fetch organizer for public display
var organizer User
h.db.First(&organizer, poll.UserID)
```

Then add the organizer fields to the response:

```go
return &gen.GetPublicPollOK{
    Name:                poll.Name,
    Description:         gen.NewOptString(poll.Description),
    CustomFields:        mapCustomFieldsToGen(poll.CustomFields),
    Options:             mapPollOptionsToGen(poll.PollOptions),
    ShowResults:         gen.NewOptBool(poll.ShowResults),
    RequireEmail:        gen.NewOptBool(poll.RequireEmail),
    OrganizerName:       gen.NewOptString(organizer.Name),
    OrganizerAvatarURL:  gen.NewOptString(avatarURL(organizer.AvatarFilename)),
}, nil
```

---

## Step 5: Verify

```bash
cd api && go build ./cmd
cd frontend && pnpm check
```

Expected: Both compile without errors.

---

## Step 6: Commit

```bash
git add api/openapi.yaml api/gen/ api/handler_public_booking.go api/handler_public_poll.go frontend/src/lib/api/generated/
git commit -m "feat(api): add organizer name and avatar to public booking/poll API responses"
```

## Summary of Changes

- Added organizer_name and organizer_avatar_url to getPublicBookingLink response in OpenAPI spec
- Added organizer_name and organizer_avatar_url to getPublicPoll response in OpenAPI spec
- Regenerated Go types with make generate
- Regenerated TypeScript types with pnpm generate:api
- Updated handler_public_booking.go to fetch organizer and include fields
- Updated handler_public_poll.go to fetch organizer and include fields

Both api and frontend compile without errors
