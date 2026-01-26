# Meet Mesh

A self-hosted scheduling application combining **Calendly-style booking links** (1:1 scheduling with real-time calendar availability) and **Doodle-style group polls** (many-to-many coordination with voting).

Single binary deployment. Single organizer per instance.

## Features

### Booking Links
- Create shareable links with available time windows
- Real-time CalDAV availability filters out conflicts
- Guests pick a slot and provide email + custom fields
- Instant booking or manual approval (configurable per link)
- Automatic calendar event creation

### Group Polls
- Create polls with specific date/time options
- Participants vote yes/no/maybe on each option
- Optional: require email, show live results
- Organizer picks winner, confirmation sent to all participants

### Slot Types
- **Time slots** - specific start/end times (e.g., "30 min at 2pm")
- **Full days** - entire days as single units
- **Multi-day ranges** - consecutive day spans

### Additional Features
- CalDAV calendar integration for availability and event creation
- OIDC authentication for the organizer
- Custom fields per link
- Configurable event templates
- Email notifications for bookings and confirmations
- Responsive dashboard UI

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.25, [ogen](https://github.com/ogen-go/ogen) (OpenAPI codegen), GORM, SQLite |
| Frontend | SvelteKit, Svelte 5 (runes), Tailwind CSS v4, bits-ui |
| Auth | OIDC (organizer), JWT tokens (email actions) |
| Calendar | CalDAV via [go-webdav](https://github.com/emersion/go-webdav) |
| Email | SMTP via gomail |
| API | OpenAPI 3.0 spec-first design |

## Prerequisites

- Go 1.25+
- Node.js 20+ and pnpm
- An OIDC provider (Keycloak, Authentik, Auth0, etc.)
- SMTP server for email notifications
- CalDAV server (Nextcloud, Radicale, etc.)

## Quick Start

### 1. Clone and build

```bash
git clone https://github.com/kolaente/meet-mesh.git
cd meet-mesh
make build
```

This builds the frontend, generates API code, and compiles everything into a single `meet-mesh` binary.

### 2. Configure

```bash
cp config.example.yaml config.yaml
```

Edit `config.yaml` with your settings:

```yaml
server:
  port: 8080
  base_url: http://localhost:8080

database:
  path: ./data/meet-mesh.db

oidc:
  issuer: https://auth.example.com
  client_id: meet-mesh
  client_secret: ${OIDC_CLIENT_SECRET}
  redirect_uri: http://localhost:8080/api/auth/callback

smtp:
  host: smtp.example.com
  port: 587
  username: meet-mesh@example.com
  password: ${SMTP_PASSWORD}
  from: "Meet Mesh <meet-mesh@example.com>"
```

Secrets support environment variable interpolation with `${VAR_NAME}` syntax.

### 3. Run

```bash
./meet-mesh -config config.yaml
```

Open http://localhost:8080 and log in with your OIDC provider.

### 4. Initial setup

1. Go to **Settings** and connect your CalDAV calendar(s)
2. Create your first booking link or poll
3. Share the generated URL with guests

## Development

Run the frontend and backend in separate terminals:

```bash
# Terminal 1: Frontend dev server (hot reload)
cd frontend && pnpm install && pnpm dev

# Terminal 2: Backend
cd api && go run ./cmd
```

The frontend dev server proxies API requests to the backend.

### Local Email Testing with Mailhog

For local development, [Mailhog](https://github.com/mailhog/MailHog) is included to capture outgoing emails:

```bash
# Start Mailhog alongside other services
devenv up
```

Then use the development config:

```bash
cp config.dev.yaml config.yaml
# Edit OIDC settings as needed
```

- **SMTP**: localhost:1025 (captured by Mailhog)
- **Web UI**: http://localhost:8025 (view captured emails)

All emails sent by the application will appear in the Mailhog web interface instead of being delivered.

### Build Commands

```bash
make build          # Full build: frontend + codegen + Go binary
make build-api      # Build Go binary only (assumes frontend built)
make frontend-dist  # Build frontend and copy to api/embedded/dist
make generate       # Regenerate Go code from OpenAPI spec
make clean          # Remove all build artifacts
```

### Frontend Commands

```bash
cd frontend
pnpm dev              # Development server
pnpm build            # Production build
pnpm check            # TypeScript/Svelte type checking
pnpm generate:api     # Regenerate TypeScript types from OpenAPI spec
```

### OpenAPI-First Development

1. Modify `api/openapi.yaml`
2. Run `make generate` to regenerate Go types/interfaces
3. Implement handler methods in `api/handler_*.go`
4. Run `cd frontend && pnpm generate:api` for TypeScript types

## Project Structure

```
meet-mesh/
├── api/
│   ├── cmd/main.go           # Entry point
│   ├── openapi.yaml          # API specification (source of truth)
│   ├── gen/                  # Generated ogen code (DO NOT EDIT)
│   ├── handler_*.go          # Request handlers by domain
│   ├── models.go             # GORM data models
│   ├── security.go           # OIDC & token validation
│   ├── caldav.go             # CalDAV client
│   ├── mailer.go             # Email sending
│   └── embedded/dist/        # Built frontend (generated)
│
├── frontend/
│   ├── src/lib/
│   │   ├── components/       # Svelte components
│   │   │   ├── ui/           # Button, Card, Input, etc.
│   │   │   ├── booking/      # DateCalendar, TimeSlotList, BookingForm
│   │   │   ├── poll/         # VoteCard, VoteButtons, VoteSummary
│   │   │   └── dashboard/    # Sidebar, LinkCard, StatsCard
│   │   ├── api/
│   │   │   ├── client.ts     # API client wrapper
│   │   │   └── generated/    # Generated from OpenAPI (DO NOT EDIT)
│   │   └── stores/           # Auth, toast stores
│   └── src/routes/
│       ├── (dashboard)/      # Organizer views (auth required)
│       ├── (public)/p/[slug] # Guest booking/poll pages
│       └── (actions)/        # Email action pages
│
├── docs/                     # Documentation
├── Makefile
└── config.example.yaml
```

## Authentication Model

| User Type | Method | Use Case |
|-----------|--------|----------|
| Organizer | OIDC session cookie | Dashboard, link management |
| Guest | None | Booking/voting (anonymous) |
| Email Actions | JWT in query param | Approve/decline via email links |

## User Flows

### Organizer

1. **Setup**: Deploy, configure OIDC/SMTP, connect CalDAV calendars
2. **Create booking link**: Name, slot type, availability rules, custom fields
3. **Create poll**: Name, add date/time options, configure visibility
4. **Manage**: View bookings, approve/decline, pick poll winners

### Guest

1. **Booking**: Open link, view available slots, pick one, submit email + fields
2. **Voting**: Open poll, vote on options, optionally provide email

## Documentation

- `docs/spec.md` - Feature specification and user flows
- `docs/technical-overview.md` - Architecture and data models

## License

MIT
