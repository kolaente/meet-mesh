# Meet Mesh - Design Document

A self-hosted scheduling app combining Calendly-style booking with Doodle-style group coordination.

## Core Concepts

### Two Modes

**Booking Links** (1:1 scheduling)
- Create a link with available time windows
- Real-time CalDAV availability filters out conflicts
- Optional: add/override with specific free slots per link
- Guest picks a slot, provides email + custom fields
- Books instantly OR waits for approval (per-link setting)

**Group Polls** (many:many coordination)
- Create a poll with specific date/time options (manual selection)
- Share link with participants
- Participants vote on which times work (email optional)
- Organizer sees votes and picks the winner
- Confirmation sent to everyone who provided email

### Slot Types

Each link is configured with a slot type:
- **Time slots** - specific start/end times (e.g., "30 min at 2pm")
- **Full days** - entire days as single units (e.g., "Monday")
- **Multi-day ranges** - consecutive day spans (e.g., "Mon-Wed")

### Shared Foundation

- CalDAV calendar(s) as source of truth for availability
- Custom fields per link (optional)
- Configurable event templates (what lands on calendar)
- Simple shareable URLs

## Users

**Organizer**
- Authenticates via OIDC
- Connects CalDAV calendars
- Creates and manages booking links and polls
- Single-user deployment model

**Guests**
- No account required
- Interact via shared links only
- Email required for bookings, optional for polls

## User Flows

### Organizer

**Setup (one-time)**
1. Deploy app, configure OIDC and SMTP via config file
2. Log in, connect CalDAV calendar(s)
3. Availability pulled automatically

**Creating a booking link**
1. Name the link (e.g., "Coffee chat", "Consulting call")
2. Choose slot type: time slots / full days / multi-day
3. Set availability rules OR add specific free slots
4. Configure: auto-confirm or approval required
5. Add custom fields (optional)
6. Define event template
7. Get shareable URL

**Creating a group poll**
1. Name the poll (e.g., "Team offsite dates")
2. Choose slot type
3. Add date/time options manually
4. Configure: require email or not, show results or not
5. Add custom fields (optional)
6. Define event template
7. Get shareable URL

**Managing bookings/polls**
- Dashboard shows all active links and polls
- See incoming bookings (approve/decline if manual mode)
- See poll votes, pick winner when ready

### Guest

**Booking a slot**
1. Open shared link (no login)
2. See available slots (real-time availability)
3. Pick based on slot type:
   - Time slot: select date, then time
   - Full day: select a date
   - Multi-day: select start and end date
4. Fill in email + custom fields
5. Submit
6. Receive confirmation email when booked

**Voting on a group poll**
1. Open shared link (no login)
2. See proposed options
3. Vote: yes / no / maybe for each option
4. Optionally provide name and/or email
5. Submit, see current tally (if enabled)
6. Receive email when organizer picks winner (if email provided)

## Data Model

### User (Organizer)
- OIDC subject identifier
- Connected calendars (1 or more)

### Calendar Connection
- CalDAV server URL
- Credentials (stored encrypted)
- Which calendars to read for availability
- Which calendar to write events to

### Link
- Type: `booking` or `poll`
- Name / description
- Slot type: `time` / `full_day` / `multi_day`
- Status: `active` / `closed`
- Settings:
  - For bookings: availability rules, auto-confirm flag
  - For polls: show results flag, require email flag
- Custom fields definition
- Event template
- Unique URL slug

### Slot / Option
- Belongs to a Link
- Start datetime, end datetime
- For bookings: auto-generated from availability OR manually added

### Booking
- Selected slot
- Guest email + custom field responses
- Status: `pending` / `confirmed` / `declined`

### Vote
- Guest identifier (email or anonymous name)
- Map of option ID to response (yes/no/maybe)

## Key Interactions

### Calendar Availability (booking links)
1. Guest opens booking link
2. App fetches organizer's CalDAV calendar(s) in real-time
3. Merges with link's availability rules / manual free slots
4. Returns available slots (busy times excluded)
5. Guest selects, app double-checks availability before confirming

### Booking Flow
1. Guest submits slot + email + custom fields
2. App verifies slot still available (re-check CalDAV)
3. If auto-confirm: create calendar event, send confirmation
4. If manual: store as pending, email organizer
5. Organizer approves via email link or dashboard
6. Create event, send confirmation to guest

### Poll Flow
1. Organizer creates poll with manual options
2. Guests vote (stored per-guest)
3. Organizer views votes, picks winner
4. App creates calendar event using template
5. Sends confirmation to all guests with email

## Notifications

### Emails Sent

| Event | Recipient | Content |
|-------|-----------|---------|
| Booking confirmed | Guest | Time, date, event details |
| Booking confirmed | Organizer | Guest info, custom fields, time |
| Booking pending | Organizer | Guest info, requested time, approve/decline link |
| Poll winner picked | Guests with email | Selected date/time, event details |

### Not Sent
- New votes on polls
- Reminders (calendar apps handle this)
- Digest summaries

### Configuration
- SMTP settings in config file
- Calendar event serves as the reminder mechanism

## UI Structure

### Organizer Views (authenticated)

| Page | Purpose |
|------|---------|
| Dashboard | Overview of all links, recent activity |
| Link list | All booking links and polls, filter by status |
| Create/edit link | Form to configure booking link or poll |
| Link detail | View bookings/votes, approve/decline, pick winner |
| Calendar settings | Connect/manage CalDAV calendars |

### Guest Views (public)

| Page | Purpose |
|------|---------|
| Booking page | View available slots, select, submit |
| Poll page | View options, vote, see results |
| Confirmation | Success message after submission |

### Email Action Pages

| Page | Purpose |
|------|---------|
| Approve/decline | One-click action from email (token auth) |

## Configuration

All infrastructure config via config file (not web UI):
- OIDC provider settings
- SMTP settings
- Database connection
- CalDAV defaults (optional)

## Scope

### MVP Includes
- Self-hosted deployment
- Two modes: booking links + group polls
- CalDAV calendar integration
- OIDC authentication
- Custom fields and event templates
- Configurable slot types (time/day/multi-day)
- Auto/manual confirmation
- Minimal email notifications

### Out of Scope (MVP)
- Google/Outlook calendar support
- Hosted/SaaS version
- Reminders and digest emails
- Built-in invite sending
- Smart slot suggestions

### Future Enhancements
- Auto-vote on polls based on calendar availability
- Smart slot suggestions when creating polls
- Additional calendar providers
- Hosted version
