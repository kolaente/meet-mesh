---
# meet-mesh-us19
title: Add organizer avatar to email notification templates
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:44:59Z
parent: meet-mesh-us07
blocked_by:
    - meet-mesh-us09
    - meet-mesh-us13
---

# Add Organizer Avatar to Email Notification Templates

**Goal:** Include the organizer's avatar in email notification templates sent to guests, giving a more personal and branded appearance.

**Architecture:** Update the email templates in `api/mailer.go` to include the organizer avatar as a linked image (not inline/embedded, to keep email size small and avoid spam filters). The avatar URL is constructed as an absolute URL using the configured `base_url`. If no avatar is set, the avatar section is omitted from the email.

The mailer methods that send emails to guests (`SendBookingConfirmation`, `SendBookingApproved`, `SendBookingDeclined`, `SendBookingConfirmationWithICS`, `SendBookingApprovedWithICS`, `SendPollWinner`) need to receive the organizer's avatar URL as template data.

---

## Files

- Modify: `api/mailer.go` (update templates and send methods)
- Modify: `api/handler_public_booking.go` (pass organizer avatar to mailer calls, if not already)
- Modify: `api/handler_bookings.go` (pass organizer avatar to mailer calls)
- Modify: `api/handler_actions.go` (pass organizer avatar to mailer calls)

---

## Step 1: Update email template data to include avatar

Open `api/mailer.go`.

First, update each mailer method to accept or construct the organizer avatar URL. The simplest approach is to add `OrganizerAvatarURL` and `OrganizerName` to each template data map.

For example, in `SendBookingConfirmation`:

```go
func (m *Mailer) SendBookingConfirmation(booking *Booking, link *BookingLink, organizer *User) error {
    body := m.renderTemplate("booking_confirmed_guest", map[string]any{
        "LinkName":            link.Name,
        "GuestName":           booking.GuestName,
        "Time":                booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
        "OrganizerName":       organizer.Name,
        "OrganizerAvatarURL":  m.organizerAvatarURL(organizer),
    })
    return m.send(booking.GuestEmail, "Booking Confirmed: "+link.Name, body)
}
```

Add a helper method to the Mailer:

```go
// organizerAvatarURL returns the absolute URL for the organizer's avatar,
// or empty string if no avatar is set.
func (m *Mailer) organizerAvatarURL(organizer *User) string {
    if organizer.AvatarFilename == "" {
        return ""
    }
    return m.baseURL + "/api/avatars/" + organizer.AvatarFilename
}
```

Apply the same pattern to all mailer methods that send emails to guests. The exact set of methods to update:

1. `SendBookingConfirmation` -- may need signature change to accept `*User`
2. `SendBookingConfirmationWithICS` -- already receives `organizerEmail`, change to receive `*User`
3. `SendBookingPending` -- already receives `*User`
4. `SendBookingApproved` -- may need signature change
5. `SendBookingApprovedWithICS` -- may need signature change
6. `SendBookingDeclined` -- may need signature change
7. `SendPollWinner` -- may need to load organizer

**Important:** Changing method signatures will require updating all call sites. Check `handler_public_booking.go`, `handler_bookings.go`, and `handler_actions.go` for calls to these methods and update them to pass the organizer `*User`.

---

## Step 2: Update email templates

In the `emailTemplates` constant at the bottom of `api/mailer.go`, update each template to include an avatar header. Add this block at the top of each email body (inside `<body>`, before the `<h1>`):

```html
{{if .OrganizerAvatarURL}}
<div style="margin-bottom: 16px;">
  <img src="{{.OrganizerAvatarURL}}" alt="{{.OrganizerName}}" width="48" height="48" style="border-radius: 50%; width: 48px; height: 48px; object-fit: cover;" />
</div>
{{end}}
```

For example, the `booking_confirmed_guest` template becomes:

```html
{{define "booking_confirmed_guest"}}
<html>
<body>
{{if .OrganizerAvatarURL}}
<div style="margin-bottom: 16px;">
  <img src="{{.OrganizerAvatarURL}}" alt="{{.OrganizerName}}" width="48" height="48" style="border-radius: 50%; width: 48px; height: 48px; object-fit: cover;" />
</div>
{{end}}
<h1>Booking Confirmed!</h1>
<p>Hi {{.GuestName}},</p>
<p>Your booking for <strong>{{.LinkName}}</strong> has been confirmed.</p>
<p><strong>When:</strong> {{.Time}}</p>
{{if .MeetingLink}}
<p><strong>Meeting Link:</strong> <a href="{{.MeetingLink}}">{{.MeetingLink}}</a></p>
{{end}}
<p style="margin-top: 20px; padding: 15px; background: #f0f9ff; border-radius: 8px;">
&#128197; <strong>Add to your calendar:</strong> Open the attached <code>invite.ics</code> file to add this event to your calendar.
</p>
</body>
</html>
{{end}}
```

Apply the same avatar block to all templates: `booking_pending`, `booking_approved`, `booking_declined`, `poll_winner`.

---

## Step 3: Update call sites

Search for all calls to the mailer methods and update them to pass the organizer `*User`. The key files to check:

- `api/handler_public_booking.go` -- `CreateBooking` calls `SendBookingConfirmation` or `SendBookingPending`
- `api/handler_bookings.go` -- `ApproveBooking`/`DeclineBooking` calls `SendBookingApproved`/`SendBookingDeclined`
- `api/handler_actions.go` -- action handlers call mailer methods
- `api/handler_public_poll.go` -- may call `SendPollWinner`

In each case, load the organizer via the link/poll's `UserID`:

```go
var organizer User
if err := h.db.First(&organizer, link.UserID).Error; err != nil {
    // handle error
}
```

Then pass `&organizer` to the mailer method.

---

## Step 4: Verify

```bash
cd api && go build ./cmd
```

Expected: Compiles without errors.

Test by triggering a booking notification and checking the email includes the avatar image (if set).

---

## Step 5: Commit

```bash
git add api/mailer.go api/handler_public_booking.go api/handler_bookings.go api/handler_actions.go api/handler_public_poll.go
git commit -m "feat(api): add organizer avatar to email notification templates"
```

## Summary of Changes

- Added `organizerAvatarURL()` helper method to Mailer struct
- Updated all guest-facing email methods to accept `*User` instead of just email string:
  - `SendBookingConfirmation`
  - `SendBookingConfirmationWithICS`
  - `SendBookingApproved`
  - `SendBookingApprovedWithICS`
  - `SendBookingDeclined`
  - `SendPollWinner`
- Added `OrganizerName` and `OrganizerAvatarURL` to all template data maps
- Updated email templates to conditionally display avatar image when available
- Updated all call sites in:
  - `handler_public_booking.go`
  - `handler_bookings.go`
  - `handler_actions.go`
  - `handler_polls.go`

Build verified: `go build ./...` passes without errors
