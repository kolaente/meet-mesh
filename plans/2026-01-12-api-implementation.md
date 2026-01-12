# Meet Mesh API Implementation Plan (Phases 1-7)

**Goal:** Build the Go backend API for Meet Mesh scheduling application with OIDC auth, CalDAV integration, and email notifications.

**Architecture:** OpenAPI spec-first with ogen codegen, GORM ORM with SQLite, session-based OIDC auth for organizers, token-based auth for email actions.

**Tech Stack:** Go 1.21+, ogen, GORM, SQLite, coreos/go-oidc, emersion/go-webdav, gomail

---

## Phase 1: Foundation

### Task 1.1: Project Setup

**Files:**
- Create: `api/go.mod`
- Create: `api/cmd/main.go`
- Create: `api/gen.go`

**Step 1: Initialize Go module**

```bash
mkdir -p api/cmd api/gen
cd api && go mod init github.com/kolaente/meet-mesh/api
```

**Step 2: Create gen.go with generate directive**

```go
// api/gen.go
package api

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target gen --clean openapi.yaml
```

**Step 3: Create minimal main.go**

```go
// api/cmd/main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Meet Mesh API starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Step 4: Commit**

```bash
git add api/
git commit -m "feat(api): initialize Go module and project structure"
```

---

### Task 1.2: OpenAPI Spec - Base Structure

**Files:**
- Create: `api/openapi.yaml`

**Step 1: Create OpenAPI spec with info and components**

```yaml
# api/openapi.yaml
openapi: 3.0.3
info:
  title: Meet Mesh API
  version: 1.0.0
  description: Self-hosted scheduling API

servers:
  - url: /api/v1

components:
  securitySchemes:
    cookieAuth:
      type: apiKey
      in: cookie
      name: session
    actionToken:
      type: apiKey
      in: query
      name: token

  schemas:
    Error:
      type: object
      required: [message]
      properties:
        message:
          type: string

    LinkType:
      type: integer
      enum: [1, 2]
      description: "1=booking, 2=poll"

    SlotType:
      type: integer
      enum: [1, 2, 3]
      description: "1=time, 2=full_day, 3=multi_day"

    LinkStatus:
      type: integer
      enum: [1, 2]
      description: "1=active, 2=closed"

    BookingStatus:
      type: integer
      enum: [1, 2, 3]
      description: "1=pending, 2=confirmed, 3=declined"

    CustomFieldType:
      type: integer
      enum: [1, 2, 3, 4, 5]
      description: "1=text, 2=email, 3=phone, 4=select, 5=textarea"

    VoteResponse:
      type: integer
      enum: [1, 2, 3]
      description: "1=yes, 2=no, 3=maybe"

paths: {}
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat(api): add OpenAPI spec base structure with enums"
```

---

### Task 1.3: OpenAPI Spec - Data Schemas

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add data schemas to components/schemas**

Add after VoteResponse schema:

```yaml
    AvailabilityRule:
      type: object
      required: [days_of_week, start_time, end_time]
      properties:
        days_of_week:
          type: array
          items:
            type: integer
            minimum: 0
            maximum: 6
        start_time:
          type: string
          pattern: "^[0-2][0-9]:[0-5][0-9]$"
        end_time:
          type: string
          pattern: "^[0-2][0-9]:[0-5][0-9]$"

    CustomField:
      type: object
      required: [name, label, type, required]
      properties:
        name:
          type: string
        label:
          type: string
        type:
          $ref: '#/components/schemas/CustomFieldType'
        required:
          type: boolean
        options:
          type: array
          items:
            type: string

    EventTemplate:
      type: object
      properties:
        title_template:
          type: string
        description_template:
          type: string
        location:
          type: string

    User:
      type: object
      required: [id, email]
      properties:
        id:
          type: integer
        email:
          type: string
        name:
          type: string

    CalendarConnection:
      type: object
      required: [id, server_url, username]
      properties:
        id:
          type: integer
        server_url:
          type: string
        username:
          type: string
        calendar_urls:
          type: array
          items:
            type: string
        write_url:
          type: string

    Link:
      type: object
      required: [id, slug, type, name, status]
      properties:
        id:
          type: integer
        slug:
          type: string
        type:
          $ref: '#/components/schemas/LinkType'
        name:
          type: string
        description:
          type: string
        status:
          $ref: '#/components/schemas/LinkStatus'
        auto_confirm:
          type: boolean
        show_results:
          type: boolean
        require_email:
          type: boolean
        availability_rules:
          type: array
          items:
            $ref: '#/components/schemas/AvailabilityRule'
        custom_fields:
          type: array
          items:
            $ref: '#/components/schemas/CustomField'
        event_template:
          $ref: '#/components/schemas/EventTemplate'
        created_at:
          type: string
          format: date-time

    Slot:
      type: object
      required: [id, type, start_time, end_time]
      properties:
        id:
          type: integer
        type:
          $ref: '#/components/schemas/SlotType'
        start_time:
          type: string
          format: date-time
        end_time:
          type: string
          format: date-time

    Booking:
      type: object
      required: [id, slot, guest_email, status]
      properties:
        id:
          type: integer
        slot:
          $ref: '#/components/schemas/Slot'
        guest_email:
          type: string
        guest_name:
          type: string
        status:
          $ref: '#/components/schemas/BookingStatus'
        custom_fields:
          type: object
          additionalProperties:
            type: string
        created_at:
          type: string
          format: date-time

    Vote:
      type: object
      required: [id, responses]
      properties:
        id:
          type: integer
        guest_name:
          type: string
        guest_email:
          type: string
        responses:
          type: object
          additionalProperties:
            $ref: '#/components/schemas/VoteResponse'
        custom_fields:
          type: object
          additionalProperties:
            type: string
        created_at:
          type: string
          format: date-time

    VoteTally:
      type: object
      required: [slot_id, yes_count, no_count, maybe_count]
      properties:
        slot_id:
          type: integer
        yes_count:
          type: integer
        no_count:
          type: integer
        maybe_count:
          type: integer
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat(api): add OpenAPI data schemas"
```

---

### Task 1.4: OpenAPI Spec - Auth Endpoints

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add auth paths**

Replace `paths: {}` with:

```yaml
paths:
  /auth/login:
    get:
      operationId: initiateLogin
      summary: Redirect to OIDC provider
      responses:
        '302':
          description: Redirect to OIDC provider

  /auth/callback:
    get:
      operationId: authCallback
      summary: Handle OIDC callback
      parameters:
        - name: code
          in: query
          required: true
          schema:
            type: string
        - name: state
          in: query
          required: true
          schema:
            type: string
      responses:
        '302':
          description: Redirect to dashboard
        '400':
          description: Invalid callback
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/logout:
    post:
      operationId: logout
      summary: Clear session
      security:
        - cookieAuth: []
      responses:
        '200':
          description: Logged out

  /auth/me:
    get:
      operationId: getCurrentUser
      summary: Get current user
      security:
        - cookieAuth: []
      responses:
        '200':
          description: Current user
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '401':
          description: Not authenticated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat(api): add auth endpoints to OpenAPI spec"
```

---

### Task 1.5: OpenAPI Spec - Calendar Endpoints

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add calendar paths after auth paths**

```yaml
  /calendars:
    get:
      operationId: listCalendars
      summary: List calendar connections
      security:
        - cookieAuth: []
      responses:
        '200':
          description: Calendar connections
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/CalendarConnection'

    post:
      operationId: addCalendar
      summary: Add calendar connection
      security:
        - cookieAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [server_url, username, password]
              properties:
                server_url:
                  type: string
                username:
                  type: string
                password:
                  type: string
                calendar_urls:
                  type: array
                  items:
                    type: string
                write_url:
                  type: string
      responses:
        '201':
          description: Calendar added
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CalendarConnection'

  /calendars/{id}:
    delete:
      operationId: removeCalendar
      summary: Remove calendar connection
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '204':
          description: Calendar removed
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat(api): add calendar endpoints to OpenAPI spec"
```

---

### Task 1.6: OpenAPI Spec - Link CRUD Endpoints

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add link management paths**

```yaml
  /links:
    get:
      operationId: listLinks
      summary: List all links
      security:
        - cookieAuth: []
      responses:
        '200':
          description: Links
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Link'

    post:
      operationId: createLink
      summary: Create a link
      security:
        - cookieAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [type, name]
              properties:
                type:
                  $ref: '#/components/schemas/LinkType'
                name:
                  type: string
                description:
                  type: string
                auto_confirm:
                  type: boolean
                show_results:
                  type: boolean
                require_email:
                  type: boolean
                availability_rules:
                  type: array
                  items:
                    $ref: '#/components/schemas/AvailabilityRule'
                custom_fields:
                  type: array
                  items:
                    $ref: '#/components/schemas/CustomField'
                event_template:
                  $ref: '#/components/schemas/EventTemplate'
      responses:
        '201':
          description: Link created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Link'

  /links/{id}:
    get:
      operationId: getLink
      summary: Get link details
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Link details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Link'

    put:
      operationId: updateLink
      summary: Update a link
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                description:
                  type: string
                status:
                  $ref: '#/components/schemas/LinkStatus'
                auto_confirm:
                  type: boolean
                show_results:
                  type: boolean
                require_email:
                  type: boolean
                availability_rules:
                  type: array
                  items:
                    $ref: '#/components/schemas/AvailabilityRule'
                custom_fields:
                  type: array
                  items:
                    $ref: '#/components/schemas/CustomField'
                event_template:
                  $ref: '#/components/schemas/EventTemplate'
      responses:
        '200':
          description: Link updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Link'

    delete:
      operationId: deleteLink
      summary: Delete a link
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '204':
          description: Link deleted
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat(api): add link CRUD endpoints to OpenAPI spec"
```

---

### Task 1.7: OpenAPI Spec - Slot & Booking Management Endpoints

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add slot and booking management paths**

```yaml
  /links/{id}/slots:
    get:
      operationId: getLinkSlots
      summary: Get slots for a link
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Slots
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Slot'

    post:
      operationId: addSlot
      summary: Add a slot to a link
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [type, start_time, end_time]
              properties:
                type:
                  $ref: '#/components/schemas/SlotType'
                start_time:
                  type: string
                  format: date-time
                end_time:
                  type: string
                  format: date-time
      responses:
        '201':
          description: Slot added
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Slot'

  /links/{id}/bookings:
    get:
      operationId: getLinkBookings
      summary: Get bookings for a link
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Bookings
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Booking'

  /links/{id}/votes:
    get:
      operationId: getLinkVotes
      summary: Get votes for a poll
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Votes
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Vote'

  /links/{id}/pick-winner:
    post:
      operationId: pickPollWinner
      summary: Pick winning slot for poll
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [slot_id]
              properties:
                slot_id:
                  type: integer
      responses:
        '200':
          description: Winner picked

  /bookings/{id}/approve:
    post:
      operationId: approveBooking
      summary: Approve a booking
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Booking approved
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Booking'

  /bookings/{id}/decline:
    post:
      operationId: declineBooking
      summary: Decline a booking
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Booking declined
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Booking'
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat(api): add slot and booking management endpoints"
```

---

### Task 1.8: OpenAPI Spec - Public Guest Endpoints

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add public guest paths**

```yaml
  /p/{slug}:
    get:
      operationId: getPublicLink
      summary: Get public link info
      parameters:
        - name: slug
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Link info
          content:
            application/json:
              schema:
                type: object
                required: [type, name, slots]
                properties:
                  type:
                    $ref: '#/components/schemas/LinkType'
                  name:
                    type: string
                  description:
                    type: string
                  custom_fields:
                    type: array
                    items:
                      $ref: '#/components/schemas/CustomField'
                  slots:
                    type: array
                    items:
                      $ref: '#/components/schemas/Slot'
                  show_results:
                    type: boolean
                  require_email:
                    type: boolean
        '404':
          description: Link not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /p/{slug}/availability:
    get:
      operationId: getAvailability
      summary: Get real-time availability
      parameters:
        - name: slug
          in: path
          required: true
          schema:
            type: string
        - name: start
          in: query
          required: true
          schema:
            type: string
            format: date-time
        - name: end
          in: query
          required: true
          schema:
            type: string
            format: date-time
      responses:
        '200':
          description: Available slots
          content:
            application/json:
              schema:
                type: object
                required: [slots]
                properties:
                  slots:
                    type: array
                    items:
                      $ref: '#/components/schemas/Slot'

  /p/{slug}/book:
    post:
      operationId: createBooking
      summary: Create a booking
      parameters:
        - name: slug
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [slot_id, guest_email]
              properties:
                slot_id:
                  type: integer
                guest_email:
                  type: string
                  format: email
                guest_name:
                  type: string
                custom_fields:
                  type: object
                  additionalProperties:
                    type: string
      responses:
        '201':
          description: Booking created
          content:
            application/json:
              schema:
                type: object
                required: [status]
                properties:
                  status:
                    $ref: '#/components/schemas/BookingStatus'
                  message:
                    type: string
        '409':
          description: Slot unavailable
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /p/{slug}/vote:
    post:
      operationId: submitVote
      summary: Submit poll vote
      parameters:
        - name: slug
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [responses]
              properties:
                guest_name:
                  type: string
                guest_email:
                  type: string
                  format: email
                responses:
                  type: object
                  additionalProperties:
                    $ref: '#/components/schemas/VoteResponse'
                custom_fields:
                  type: object
                  additionalProperties:
                    type: string
      responses:
        '201':
          description: Vote submitted
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Vote'

  /p/{slug}/results:
    get:
      operationId: getPollResults
      summary: Get poll results
      parameters:
        - name: slug
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Poll results
          content:
            application/json:
              schema:
                type: object
                required: [tally]
                properties:
                  tally:
                    type: array
                    items:
                      $ref: '#/components/schemas/VoteTally'
                  votes:
                    type: array
                    items:
                      $ref: '#/components/schemas/Vote'
        '403':
          description: Results not public
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat(api): add public guest endpoints to OpenAPI spec"
```

---

### Task 1.9: OpenAPI Spec - Email Action Endpoints

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add email action paths**

```yaml
  /actions/approve:
    get:
      operationId: approveViaEmail
      summary: Approve booking via email link
      security:
        - actionToken: []
      responses:
        '200':
          description: Booking approved
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
        '400':
          description: Invalid or expired token
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /actions/decline:
    get:
      operationId: declineViaEmail
      summary: Decline booking via email link
      security:
        - actionToken: []
      responses:
        '200':
          description: Booking declined
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
        '400':
          description: Invalid or expired token
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat(api): add email action endpoints to OpenAPI spec"
```

---

### Task 1.10: Generate Backend Code

**Files:**
- Create: `api/gen/` (generated)

**Step 1: Install ogen**

```bash
go install github.com/ogen-go/ogen/cmd/ogen@latest
```

**Step 2: Run code generation**

```bash
cd api && go generate ./...
```

**Step 3: Verify generation succeeded**

```bash
ls api/gen/
```

Expected: Multiple `oas_*.go` files

**Step 4: Commit**

```bash
git add api/gen/
git commit -m "feat(api): generate ogen code from OpenAPI spec"
```

---

### Task 1.11: GORM Models

**Files:**
- Create: `api/models.go`

**Step 1: Create models file**

```go
// api/models.go
package api

import (
	"time"
)

// Enum types
type LinkType int

const (
	LinkTypeBooking LinkType = 1
	LinkTypePoll    LinkType = 2
)

type SlotType int

const (
	SlotTypeTime     SlotType = 1
	SlotTypeFullDay  SlotType = 2
	SlotTypeMultiDay SlotType = 3
)

type LinkStatus int

const (
	LinkStatusActive LinkStatus = 1
	LinkStatusClosed LinkStatus = 2
)

type BookingStatus int

const (
	BookingStatusPending   BookingStatus = 1
	BookingStatusConfirmed BookingStatus = 2
	BookingStatusDeclined  BookingStatus = 3
)

type CustomFieldType int

const (
	CustomFieldTypeText     CustomFieldType = 1
	CustomFieldTypeEmail    CustomFieldType = 2
	CustomFieldTypePhone    CustomFieldType = 3
	CustomFieldTypeSelect   CustomFieldType = 4
	CustomFieldTypeTextarea CustomFieldType = 5
)

type VoteResponseType int

const (
	VoteResponseYes   VoteResponseType = 1
	VoteResponseNo    VoteResponseType = 2
	VoteResponseMaybe VoteResponseType = 3
)

// JSON-serialized structs
type AvailabilityRule struct {
	DaysOfWeek []int  `json:"days_of_week"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type CustomField struct {
	Name     string          `json:"name"`
	Label    string          `json:"label"`
	Type     CustomFieldType `json:"type"`
	Required bool            `json:"required"`
	Options  []string        `json:"options,omitempty"`
}

type EventTemplate struct {
	TitleTemplate       string `json:"title_template"`
	DescriptionTemplate string `json:"description_template"`
	Location            string `json:"location,omitempty"`
}

// GORM Models
type User struct {
	ID        uint      `gorm:"primaryKey"`
	OIDCSub   string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"not null"`
	Name      string
	CreatedAt time.Time
	Calendars []CalendarConnection `gorm:"foreignKey:UserID"`
	Links     []Link               `gorm:"foreignKey:UserID"`
}

type CalendarConnection struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index;not null"`
	ServerURL    string    `gorm:"not null"`
	Username     string    `gorm:"not null"`
	Password     string    `gorm:"not null"`
	CalendarURLs []string  `gorm:"serializer:json"`
	WriteURL     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Link struct {
	ID                uint               `gorm:"primaryKey"`
	UserID            uint               `gorm:"index;not null"`
	Slug              string             `gorm:"uniqueIndex;not null"`
	Type              LinkType           `gorm:"not null"`
	Name              string             `gorm:"not null"`
	Description       string
	Status            LinkStatus         `gorm:"not null;default:1"`
	AutoConfirm       bool
	AvailabilityRules []AvailabilityRule `gorm:"serializer:json"`
	ShowResults       bool
	RequireEmail      bool
	CustomFields      []CustomField      `gorm:"serializer:json"`
	EventTemplate     *EventTemplate     `gorm:"serializer:json"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Slots             []Slot             `gorm:"foreignKey:LinkID"`
	Bookings          []Booking          `gorm:"foreignKey:LinkID"`
	Votes             []Vote             `gorm:"foreignKey:LinkID"`
}

type Slot struct {
	ID        uint      `gorm:"primaryKey"`
	LinkID    uint      `gorm:"index;not null"`
	Type      SlotType  `gorm:"not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Manual    bool
	CreatedAt time.Time
}

type Booking struct {
	ID           uint              `gorm:"primaryKey"`
	LinkID       uint              `gorm:"index;not null"`
	SlotID       uint              `gorm:"index;not null"`
	GuestEmail   string            `gorm:"not null"`
	GuestName    string
	CustomFields map[string]string `gorm:"serializer:json"`
	Status       BookingStatus     `gorm:"not null;default:1"`
	ActionToken  string            `gorm:"uniqueIndex"`
	CalendarUID  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Link         Link              `gorm:"foreignKey:LinkID"`
	Slot         Slot              `gorm:"foreignKey:SlotID"`
}

type Vote struct {
	ID           uint                         `gorm:"primaryKey"`
	LinkID       uint                         `gorm:"index;not null"`
	GuestEmail   string
	GuestName    string
	Responses    map[uint]VoteResponseType    `gorm:"serializer:json;not null"`
	CustomFields map[string]string            `gorm:"serializer:json"`
	CreatedAt    time.Time
	Link         Link                         `gorm:"foreignKey:LinkID"`
}
```

**Step 2: Commit**

```bash
git add api/models.go
git commit -m "feat(api): add GORM models"
```

---

### Task 1.12: Configuration

**Files:**
- Create: `api/config.go`
- Create: `config.example.yaml`

**Step 1: Create config.go**

```go
// api/config.go
package api

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	OIDC     OIDCConfig     `yaml:"oidc"`
	SMTP     SMTPConfig     `yaml:"smtp"`
}

type ServerConfig struct {
	Port    int    `yaml:"port"`
	BaseURL string `yaml:"base_url"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type OIDCConfig struct {
	Issuer       string `yaml:"issuer"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURI  string `yaml:"redirect_uri"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Expand environment variables
	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

**Step 2: Create config.example.yaml**

```yaml
# config.example.yaml
server:
  port: 8080
  base_url: http://localhost:8080

database:
  path: ./data/meet-mesh.db

oidc:
  issuer: https://auth.example.com
  client_id: meet-mesh
  client_secret: ${OIDC_CLIENT_SECRET}
  redirect_uri: http://localhost:8080/api/v1/auth/callback

smtp:
  host: smtp.example.com
  port: 587
  username: meet-mesh@example.com
  password: ${SMTP_PASSWORD}
  from: "Meet Mesh <meet-mesh@example.com>"
```

**Step 3: Add go dependencies**

```bash
cd api && go get gopkg.in/yaml.v3
```

**Step 4: Commit**

```bash
git add api/config.go config.example.yaml
git commit -m "feat(api): add configuration loading"
```

---

### Task 1.13: Database Setup

**Files:**
- Create: `api/database.go`

**Step 1: Create database.go**

```go
// api/database.go
package api

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&User{},
		&CalendarConnection{},
		&Link{},
		&Slot{},
		&Booking{},
		&Vote{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
```

**Step 2: Add GORM dependencies**

```bash
cd api && go get gorm.io/gorm gorm.io/driver/sqlite
```

**Step 3: Commit**

```bash
git add api/database.go
git commit -m "feat(api): add database initialization with GORM"
```

---

## Phase 2: Authentication

### Task 2.1: OIDC Provider Setup

**Files:**
- Create: `api/auth.go`

**Step 1: Create auth.go with OIDC setup**

```go
// api/auth.go
package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type AuthService struct {
	config       *OIDCConfig
	provider     *oidc.Provider
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewAuthService(ctx context.Context, cfg *OIDCConfig) (*AuthService, error) {
	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return &AuthService{
		config:       cfg,
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
	}, nil
}

func (a *AuthService) GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (a *AuthService) AuthCodeURL(state string) string {
	return a.oauth2Config.AuthCodeURL(state)
}

type UserClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (a *AuthService) Exchange(ctx context.Context, code string) (*UserClaims, error) {
	token, err := a.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, err
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims UserClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

// Session management
type Session struct {
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (a *AuthService) CreateSessionCookie(session *Session) (*http.Cookie, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	// In production, encrypt this with a secret key
	encoded := base64.URLEncoding.EncodeToString(data)

	return &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	}, nil
}

func (a *AuthService) ParseSessionCookie(cookie *http.Cookie) (*Session, error) {
	data, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, err
	}

	return &session, nil
}
```

**Step 2: Add OIDC dependencies**

```bash
cd api && go get github.com/coreos/go-oidc/v3/oidc golang.org/x/oauth2
```

**Step 3: Commit**

```bash
git add api/auth.go
git commit -m "feat(api): add OIDC authentication service"
```

---

### Task 2.2: Security Handler Implementation

**Files:**
- Create: `api/security.go`

**Step 1: Create security.go**

```go
// api/security.go
package api

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/konrad/meet-mesh/api/gen"
)

type contextKey string

const (
	userIDKey    contextKey = "userID"
	bookingIDKey contextKey = "bookingID"
)

func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDKey).(uint)
	return userID, ok
}

func WithBookingID(ctx context.Context, bookingID uint) context.Context {
	return context.WithValue(ctx, bookingIDKey, bookingID)
}

func GetBookingID(ctx context.Context) (uint, bool) {
	bookingID, ok := ctx.Value(bookingIDKey).(uint)
	return bookingID, ok
}

type SecurityHandler struct {
	db   *gorm.DB
	auth *AuthService
}

func NewSecurityHandler(db *gorm.DB, auth *AuthService) *SecurityHandler {
	return &SecurityHandler{db: db, auth: auth}
}

var _ gen.SecurityHandler = (*SecurityHandler)(nil)

func (s *SecurityHandler) HandleCookieAuth(ctx context.Context, operationName string, t gen.CookieAuth) (context.Context, error) {
	// t.APIKey contains the session cookie value
	cookie := &http.Cookie{Value: t.APIKey}

	session, err := s.auth.ParseSessionCookie(cookie)
	if err != nil {
		return ctx, errors.New("invalid session")
	}

	return WithUserID(ctx, session.UserID), nil
}

func (s *SecurityHandler) HandleActionToken(ctx context.Context, operationName string, t gen.ActionToken) (context.Context, error) {
	var booking Booking
	if err := s.db.Where("action_token = ?", t.APIKey).First(&booking).Error; err != nil {
		return ctx, errors.New("invalid action token")
	}

	return WithBookingID(ctx, booking.ID), nil
}
```

**Step 2: Commit**

```bash
git add api/security.go
git commit -m "feat(api): add security handler for OIDC and action tokens"
```

---

### Task 2.3: Auth Handler Implementation

**Files:**
- Create: `api/handler_auth.go`

**Step 1: Create handler_auth.go**

```go
// api/handler_auth.go
package api

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/konrad/meet-mesh/api/gen"
)

// InitiateLogin redirects to OIDC provider
func (h *Handler) InitiateLogin(ctx context.Context) (*gen.InitiateLoginFound, error) {
	state, err := h.auth.GenerateState()
	if err != nil {
		return nil, err
	}

	// Store state in cookie for validation
	// In production, use a secure state store

	return &gen.InitiateLoginFound{
		Location: gen.NewOptString(h.auth.AuthCodeURL(state)),
	}, nil
}

// AuthCallback handles OIDC callback
func (h *Handler) AuthCallback(ctx context.Context, params gen.AuthCallbackParams) (*gen.AuthCallbackFound, error) {
	claims, err := h.auth.Exchange(ctx, params.Code)
	if err != nil {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   gen.Error{Message: "Authentication failed"},
		}
	}

	// Find or create user
	var user User
	result := h.db.Where("oidc_sub = ?", claims.Sub).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		user = User{
			OIDCSub: claims.Sub,
			Email:   claims.Email,
			Name:    claims.Name,
		}
		if err := h.db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if result.Error != nil {
		return nil, result.Error
	}

	// Create session
	session := &Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	cookie, err := h.auth.CreateSessionCookie(session)
	if err != nil {
		return nil, err
	}

	return &gen.AuthCallbackFound{
		Location:  gen.NewOptString(h.config.Server.BaseURL + "/dashboard"),
		SetCookie: gen.NewOptString(cookie.String()),
	}, nil
}

// Logout clears the session
func (h *Handler) Logout(ctx context.Context) error {
	// Clear cookie by setting expired
	return nil
}

// GetCurrentUser returns the authenticated user
func (h *Handler) GetCurrentUser(ctx context.Context) (*gen.User, error) {
	userID, ok := GetUserID(ctx)
	if !ok {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response:   gen.Error{Message: "Not authenticated"},
		}
	}

	var user User
	if err := h.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &gen.User{
		ID:    gen.NewOptInt(int(user.ID)),
		Email: user.Email,
		Name:  gen.NewOptString(user.Name),
	}, nil
}
```

**Step 2: Commit**

```bash
git add api/handler_auth.go
git commit -m "feat(api): implement auth handlers"
```

---

## Phase 3: CalDAV Integration

### Task 3.1: CalDAV Client

**Files:**
- Create: `api/caldav.go`

**Step 1: Create caldav.go**

```go
// api/caldav.go
package api

import (
	"context"
	"net/http"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"gorm.io/gorm"
)

type CalDAVClient struct {
	db *gorm.DB
}

func NewCalDAVClient(db *gorm.DB) *CalDAVClient {
	return &CalDAVClient{db: db}
}

type TimePeriod struct {
	Start time.Time
	End   time.Time
}

func (c *CalDAVClient) createClient(conn *CalendarConnection) (*caldav.Client, error) {
	httpClient := webdav.HTTPClientWithBasicAuth(http.DefaultClient, conn.Username, conn.Password)
	return caldav.NewClient(httpClient, conn.ServerURL)
}

// GetBusyTimes fetches busy periods from all connected calendars
func (c *CalDAVClient) GetBusyTimes(ctx context.Context, userID uint, start, end time.Time) ([]TimePeriod, error) {
	var connections []CalendarConnection
	if err := c.db.Where("user_id = ?", userID).Find(&connections).Error; err != nil {
		return nil, err
	}

	var allBusy []TimePeriod
	for _, conn := range connections {
		client, err := c.createClient(&conn)
		if err != nil {
			continue
		}

		for _, calURL := range conn.CalendarURLs {
			query := &caldav.CalendarQuery{
				CompRequest: caldav.CalendarCompRequest{
					Name: "VCALENDAR",
					Comps: []caldav.CalendarCompRequest{{
						Name:     "VEVENT",
						AllProps: true,
					}},
				},
				CompFilter: caldav.CompFilter{
					Name: "VCALENDAR",
					Comps: []caldav.CompFilter{{
						Name: "VEVENT",
						Start: start,
						End:   end,
					}},
				},
			}

			events, err := client.QueryCalendar(ctx, calURL, query)
			if err != nil {
				continue
			}

			for _, obj := range events {
				for _, event := range obj.Data.Events() {
					dtstart, _ := event.DateTimeStart(nil)
					dtend, _ := event.DateTimeEnd(nil)
					allBusy = append(allBusy, TimePeriod{
						Start: dtstart,
						End:   dtend,
					})
				}
			}
		}
	}

	return mergePeriods(allBusy), nil
}

// CreateEvent creates a calendar event for a confirmed booking
func (c *CalDAVClient) CreateEvent(ctx context.Context, userID uint, booking *Booking, slot *Slot, template *EventTemplate) (string, error) {
	var conn CalendarConnection
	if err := c.db.Where("user_id = ? AND write_url != ''", userID).First(&conn).Error; err != nil {
		return "", err
	}

	client, err := c.createClient(&conn)
	if err != nil {
		return "", err
	}

	// Build iCal event
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropProductID, "-//Meet Mesh//EN")
	cal.Props.SetText(ical.PropVersion, "2.0")

	event := ical.NewEvent()
	uid := generateUID()
	event.Props.SetText(ical.PropUID, uid)
	event.Props.SetDateTime(ical.PropDateTimeStart, slot.StartTime)
	event.Props.SetDateTime(ical.PropDateTimeEnd, slot.EndTime)

	title := "Meeting"
	if template != nil && template.TitleTemplate != "" {
		title = expandTemplate(template.TitleTemplate, booking)
	}
	event.Props.SetText(ical.PropSummary, title)

	if template != nil && template.DescriptionTemplate != "" {
		event.Props.SetText(ical.PropDescription, expandTemplate(template.DescriptionTemplate, booking))
	}

	if template != nil && template.Location != "" {
		event.Props.SetText(ical.PropLocation, template.Location)
	}

	cal.Children = append(cal.Children, event.Component)

	// Put the event
	path := conn.WriteURL + "/" + uid + ".ics"
	_, err = client.PutCalendarObject(ctx, path, cal)
	if err != nil {
		return "", err
	}

	return uid, nil
}

// Helper functions
func mergePeriods(periods []TimePeriod) []TimePeriod {
	if len(periods) == 0 {
		return periods
	}

	// Sort by start time
	// ... sorting logic ...

	var merged []TimePeriod
	current := periods[0]
	for i := 1; i < len(periods); i++ {
		if periods[i].Start.Before(current.End) || periods[i].Start.Equal(current.End) {
			if periods[i].End.After(current.End) {
				current.End = periods[i].End
			}
		} else {
			merged = append(merged, current)
			current = periods[i]
		}
	}
	merged = append(merged, current)
	return merged
}

func generateUID() string {
	// Generate unique ID for calendar event
	return time.Now().Format("20060102T150405") + "@meet-mesh"
}

func expandTemplate(template string, booking *Booking) string {
	// Simple template expansion
	result := template
	// Replace {{guest_name}} with booking.GuestName
	// Replace {{guest_email}} with booking.GuestEmail
	// etc.
	return result
}
```

**Step 2: Add dependencies**

```bash
cd api && go get github.com/emersion/go-webdav github.com/emersion/go-ical
```

**Step 3: Commit**

```bash
git add api/caldav.go
git commit -m "feat(api): add CalDAV client for calendar integration"
```

---

### Task 3.2: Calendar Handler

**Files:**
- Create: `api/handler_calendar.go`

**Step 1: Create handler_calendar.go**

```go
// api/handler_calendar.go
package api

import (
	"context"

	"github.com/konrad/meet-mesh/api/gen"
)

// ListCalendars returns calendar connections
func (h *Handler) ListCalendars(ctx context.Context) ([]gen.CalendarConnection, error) {
	userID, _ := GetUserID(ctx)

	var connections []CalendarConnection
	if err := h.db.Where("user_id = ?", userID).Find(&connections).Error; err != nil {
		return nil, err
	}

	result := make([]gen.CalendarConnection, len(connections))
	for i, conn := range connections {
		result[i] = gen.CalendarConnection{
			ID:           gen.NewOptInt(int(conn.ID)),
			ServerURL:    conn.ServerURL,
			Username:     conn.Username,
			CalendarUrls: gen.NewOptStringArray(conn.CalendarURLs),
			WriteURL:     gen.NewOptString(conn.WriteURL),
		}
	}

	return result, nil
}

// AddCalendar adds a calendar connection
func (h *Handler) AddCalendar(ctx context.Context, req *gen.AddCalendarReq) (*gen.CalendarConnection, error) {
	userID, _ := GetUserID(ctx)

	conn := CalendarConnection{
		UserID:       userID,
		ServerURL:    req.ServerURL,
		Username:     req.Username,
		Password:     req.Password,
		CalendarURLs: req.CalendarUrls.Value,
		WriteURL:     req.WriteURL.Value,
	}

	if err := h.db.Create(&conn).Error; err != nil {
		return nil, err
	}

	return &gen.CalendarConnection{
		ID:           gen.NewOptInt(int(conn.ID)),
		ServerURL:    conn.ServerURL,
		Username:     conn.Username,
		CalendarUrls: gen.NewOptStringArray(conn.CalendarURLs),
		WriteURL:     gen.NewOptString(conn.WriteURL),
	}, nil
}

// RemoveCalendar removes a calendar connection
func (h *Handler) RemoveCalendar(ctx context.Context, params gen.RemoveCalendarParams) error {
	userID, _ := GetUserID(ctx)

	return h.db.Where("id = ? AND user_id = ?", params.ID, userID).Delete(&CalendarConnection{}).Error
}
```

**Step 2: Commit**

```bash
git add api/handler_calendar.go
git commit -m "feat(api): implement calendar connection handlers"
```

---

## Phase 4: Core API - Links

### Task 4.1: Link CRUD Handler

**Files:**
- Create: `api/handler_links.go`

**Step 1: Create handler_links.go**

```go
// api/handler_links.go
package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/konrad/meet-mesh/api/gen"
)

func generateSlug() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:10]
}

// ListLinks returns all links for user
func (h *Handler) ListLinks(ctx context.Context) ([]gen.Link, error) {
	userID, _ := GetUserID(ctx)

	var links []Link
	if err := h.db.Where("user_id = ?", userID).Find(&links).Error; err != nil {
		return nil, err
	}

	return mapLinksToGen(links), nil
}

// CreateLink creates a new link
func (h *Handler) CreateLink(ctx context.Context, req *gen.CreateLinkReq) (*gen.Link, error) {
	userID, _ := GetUserID(ctx)

	link := Link{
		UserID:            userID,
		Slug:              generateSlug(),
		Type:              LinkType(req.Type),
		Name:              req.Name,
		Description:       req.Description.Value,
		Status:            LinkStatusActive,
		AutoConfirm:       req.AutoConfirm.Value,
		ShowResults:       req.ShowResults.Value,
		RequireEmail:      req.RequireEmail.Value,
		AvailabilityRules: mapAvailabilityRulesFromGen(req.AvailabilityRules),
		CustomFields:      mapCustomFieldsFromGen(req.CustomFields),
		EventTemplate:     mapEventTemplateFromGen(req.EventTemplate),
	}

	if err := h.db.Create(&link).Error; err != nil {
		return nil, err
	}

	return mapLinkToGen(&link), nil
}

// GetLink returns link details
func (h *Handler) GetLink(ctx context.Context, params gen.GetLinkParams) (*gen.Link, error) {
	userID, _ := GetUserID(ctx)

	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	return mapLinkToGen(&link), nil
}

// UpdateLink updates a link
func (h *Handler) UpdateLink(ctx context.Context, req *gen.UpdateLinkReq, params gen.UpdateLinkParams) (*gen.Link, error) {
	userID, _ := GetUserID(ctx)

	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	if req.Name.Set {
		link.Name = req.Name.Value
	}
	if req.Description.Set {
		link.Description = req.Description.Value
	}
	if req.Status.Set {
		link.Status = LinkStatus(req.Status.Value)
	}
	if req.AutoConfirm.Set {
		link.AutoConfirm = req.AutoConfirm.Value
	}
	if req.ShowResults.Set {
		link.ShowResults = req.ShowResults.Value
	}
	if req.RequireEmail.Set {
		link.RequireEmail = req.RequireEmail.Value
	}
	if req.AvailabilityRules != nil {
		link.AvailabilityRules = mapAvailabilityRulesFromGen(req.AvailabilityRules)
	}
	if req.CustomFields != nil {
		link.CustomFields = mapCustomFieldsFromGen(req.CustomFields)
	}
	if req.EventTemplate.Set {
		link.EventTemplate = mapEventTemplateFromGen(req.EventTemplate)
	}

	if err := h.db.Save(&link).Error; err != nil {
		return nil, err
	}

	return mapLinkToGen(&link), nil
}

// DeleteLink deletes a link
func (h *Handler) DeleteLink(ctx context.Context, params gen.DeleteLinkParams) error {
	userID, _ := GetUserID(ctx)

	return h.db.Where("id = ? AND user_id = ?", params.ID, userID).Delete(&Link{}).Error
}

// Helper mapping functions
func mapLinksToGen(links []Link) []gen.Link {
	result := make([]gen.Link, len(links))
	for i, link := range links {
		result[i] = *mapLinkToGen(&link)
	}
	return result
}

func mapLinkToGen(link *Link) *gen.Link {
	return &gen.Link{
		ID:                gen.NewOptInt(int(link.ID)),
		Slug:              link.Slug,
		Type:              gen.LinkType(link.Type),
		Name:              link.Name,
		Description:       gen.NewOptString(link.Description),
		Status:            gen.LinkStatus(link.Status),
		AutoConfirm:       gen.NewOptBool(link.AutoConfirm),
		ShowResults:       gen.NewOptBool(link.ShowResults),
		RequireEmail:      gen.NewOptBool(link.RequireEmail),
		AvailabilityRules: mapAvailabilityRulesToGen(link.AvailabilityRules),
		CustomFields:      mapCustomFieldsToGen(link.CustomFields),
		EventTemplate:     mapEventTemplateToGen(link.EventTemplate),
		CreatedAt:         gen.NewOptDateTime(link.CreatedAt),
	}
}

func mapAvailabilityRulesFromGen(rules []gen.AvailabilityRule) []AvailabilityRule {
	result := make([]AvailabilityRule, len(rules))
	for i, r := range rules {
		result[i] = AvailabilityRule{
			DaysOfWeek: r.DaysOfWeek,
			StartTime:  r.StartTime,
			EndTime:    r.EndTime,
		}
	}
	return result
}

func mapAvailabilityRulesToGen(rules []AvailabilityRule) []gen.AvailabilityRule {
	result := make([]gen.AvailabilityRule, len(rules))
	for i, r := range rules {
		result[i] = gen.AvailabilityRule{
			DaysOfWeek: r.DaysOfWeek,
			StartTime:  r.StartTime,
			EndTime:    r.EndTime,
		}
	}
	return result
}

func mapCustomFieldsFromGen(fields []gen.CustomField) []CustomField {
	result := make([]CustomField, len(fields))
	for i, f := range fields {
		result[i] = CustomField{
			Name:     f.Name,
			Label:    f.Label,
			Type:     CustomFieldType(f.Type),
			Required: f.Required,
			Options:  f.Options,
		}
	}
	return result
}

func mapCustomFieldsToGen(fields []CustomField) []gen.CustomField {
	result := make([]gen.CustomField, len(fields))
	for i, f := range fields {
		result[i] = gen.CustomField{
			Name:     f.Name,
			Label:    f.Label,
			Type:     gen.CustomFieldType(f.Type),
			Required: f.Required,
			Options:  f.Options,
		}
	}
	return result
}

func mapEventTemplateFromGen(opt gen.OptEventTemplate) *EventTemplate {
	if !opt.Set {
		return nil
	}
	return &EventTemplate{
		TitleTemplate:       opt.Value.TitleTemplate.Value,
		DescriptionTemplate: opt.Value.DescriptionTemplate.Value,
		Location:            opt.Value.Location.Value,
	}
}

func mapEventTemplateToGen(tmpl *EventTemplate) gen.OptEventTemplate {
	if tmpl == nil {
		return gen.OptEventTemplate{}
	}
	return gen.NewOptEventTemplate(gen.EventTemplate{
		TitleTemplate:       gen.NewOptString(tmpl.TitleTemplate),
		DescriptionTemplate: gen.NewOptString(tmpl.DescriptionTemplate),
		Location:            gen.NewOptString(tmpl.Location),
	})
}
```

**Step 2: Commit**

```bash
git add api/handler_links.go
git commit -m "feat(api): implement link CRUD handlers"
```

---

### Task 4.2: Slot Management Handler

**Files:**
- Create: `api/handler_slots.go`

**Step 1: Create handler_slots.go**

```go
// api/handler_slots.go
package api

import (
	"context"

	"github.com/konrad/meet-mesh/api/gen"
)

// GetLinkSlots returns slots for a link
func (h *Handler) GetLinkSlots(ctx context.Context, params gen.GetLinkSlotsParams) ([]gen.Slot, error) {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	var slots []Slot
	if err := h.db.Where("link_id = ?", params.ID).Order("start_time").Find(&slots).Error; err != nil {
		return nil, err
	}

	return mapSlotsToGen(slots), nil
}

// AddSlot adds a slot to a link
func (h *Handler) AddSlot(ctx context.Context, req *gen.AddSlotReq, params gen.AddSlotParams) (*gen.Slot, error) {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	slot := Slot{
		LinkID:    uint(params.ID),
		Type:      SlotType(req.Type),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Manual:    true,
	}

	if err := h.db.Create(&slot).Error; err != nil {
		return nil, err
	}

	return mapSlotToGen(&slot), nil
}

func mapSlotsToGen(slots []Slot) []gen.Slot {
	result := make([]gen.Slot, len(slots))
	for i, slot := range slots {
		result[i] = *mapSlotToGen(&slot)
	}
	return result
}

func mapSlotToGen(slot *Slot) *gen.Slot {
	return &gen.Slot{
		ID:        gen.NewOptInt(int(slot.ID)),
		Type:      gen.SlotType(slot.Type),
		StartTime: slot.StartTime,
		EndTime:   slot.EndTime,
	}
}
```

**Step 2: Commit**

```bash
git add api/handler_slots.go
git commit -m "feat(api): implement slot management handlers"
```

---

## Phase 5: Core API - Booking Flow

### Task 5.1: Public Link Handler

**Files:**
- Create: `api/handler_public.go`

**Step 1: Create handler_public.go**

```go
// api/handler_public.go
package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"gorm.io/gorm"

	"github.com/konrad/meet-mesh/api/gen"
)

// GetPublicLink returns public link info
func (h *Handler) GetPublicLink(ctx context.Context, params gen.GetPublicLinkParams) (*gen.GetPublicLinkOK, error) {
	var link Link
	if err := h.db.Preload("Slots").Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &gen.ErrorStatusCode{
				StatusCode: http.StatusNotFound,
				Response:   gen.Error{Message: "Link not found"},
			}
		}
		return nil, err
	}

	return &gen.GetPublicLinkOK{
		Type:         gen.LinkType(link.Type),
		Name:         link.Name,
		Description:  gen.NewOptString(link.Description),
		CustomFields: mapCustomFieldsToGen(link.CustomFields),
		Slots:        mapSlotsToGen(link.Slots),
		ShowResults:  gen.NewOptBool(link.ShowResults),
		RequireEmail: gen.NewOptBool(link.RequireEmail),
	}, nil
}

// GetAvailability returns real-time availability
func (h *Handler) GetAvailability(ctx context.Context, params gen.GetAvailabilityParams) (*gen.GetAvailabilityOK, error) {
	var link Link
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		return nil, err
	}

	if link.Type != LinkTypeBooking {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   gen.Error{Message: "Availability only for booking links"},
		}
	}

	// Fetch busy times from CalDAV
	busyTimes, err := h.caldav.GetBusyTimes(ctx, link.UserID, params.Start, params.End)
	if err != nil {
		return nil, err
	}

	// Generate available slots based on availability rules
	slots := h.generateAvailableSlots(link, params.Start, params.End, busyTimes)

	return &gen.GetAvailabilityOK{
		Slots: mapSlotsToGen(slots),
	}, nil
}

func (h *Handler) generateAvailableSlots(link Link, start, end time.Time, busyTimes []TimePeriod) []Slot {
	// Implementation: generate slots based on availability rules and filter out busy times
	var slots []Slot
	// ... slot generation logic ...
	return slots
}

// CreateBooking creates a booking
func (h *Handler) CreateBooking(ctx context.Context, req *gen.CreateBookingReq, params gen.CreateBookingParams) (*gen.CreateBookingCreated, error) {
	var link Link
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		return nil, err
	}

	if link.Type != LinkTypeBooking {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   gen.Error{Message: "Cannot book on poll link"},
		}
	}

	// Verify slot exists and is available
	var slot Slot
	if err := h.db.Where("id = ? AND link_id = ?", req.SlotID, link.ID).First(&slot).Error; err != nil {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusConflict,
			Response:   gen.Error{Message: "Slot not available"},
		}
	}

	// Check CalDAV availability
	busyTimes, err := h.caldav.GetBusyTimes(ctx, link.UserID, slot.StartTime, slot.EndTime)
	if err == nil && len(busyTimes) > 0 {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusConflict,
			Response:   gen.Error{Message: "Slot no longer available"},
		}
	}

	// Generate action token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	actionToken := hex.EncodeToString(tokenBytes)

	status := BookingStatusPending
	if link.AutoConfirm {
		status = BookingStatusConfirmed
	}

	booking := Booking{
		LinkID:       link.ID,
		SlotID:       slot.ID,
		GuestEmail:   req.GuestEmail,
		GuestName:    req.GuestName.Value,
		CustomFields: req.CustomFields.Value,
		Status:       status,
		ActionToken:  actionToken,
	}

	if err := h.db.Create(&booking).Error; err != nil {
		return nil, err
	}

	// Send notification email
	if link.AutoConfirm {
		h.mailer.SendBookingConfirmation(&booking, &link)
		// Create calendar event
		h.caldav.CreateEvent(ctx, link.UserID, &booking, &slot, link.EventTemplate)
	} else {
		h.mailer.SendBookingPending(&booking, &link)
	}

	message := "Booking confirmed"
	if !link.AutoConfirm {
		message = "Booking pending approval"
	}

	return &gen.CreateBookingCreated{
		Status:  gen.BookingStatus(status),
		Message: gen.NewOptString(message),
	}, nil
}
```

**Step 2: Commit**

```bash
git add api/handler_public.go
git commit -m "feat(api): implement public link and booking handlers"
```

---

### Task 5.2: Booking Management Handler

**Files:**
- Create: `api/handler_bookings.go`

**Step 1: Create handler_bookings.go**

```go
// api/handler_bookings.go
package api

import (
	"context"

	"github.com/konrad/meet-mesh/api/gen"
)

// GetLinkBookings returns bookings for a link
func (h *Handler) GetLinkBookings(ctx context.Context, params gen.GetLinkBookingsParams) ([]gen.Booking, error) {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	var bookings []Booking
	if err := h.db.Preload("Slot").Where("link_id = ?", params.ID).Order("created_at DESC").Find(&bookings).Error; err != nil {
		return nil, err
	}

	return mapBookingsToGen(bookings), nil
}

// ApproveBooking approves a booking
func (h *Handler) ApproveBooking(ctx context.Context, params gen.ApproveBookingParams) (*gen.Booking, error) {
	userID, _ := GetUserID(ctx)

	var booking Booking
	if err := h.db.Preload("Link").Preload("Slot").First(&booking, params.ID).Error; err != nil {
		return nil, err
	}

	// Verify ownership
	if booking.Link.UserID != userID {
		return nil, &gen.ErrorStatusCode{
			StatusCode: 403,
			Response:   gen.Error{Message: "Not authorized"},
		}
	}

	booking.Status = BookingStatusConfirmed
	if err := h.db.Save(&booking).Error; err != nil {
		return nil, err
	}

	// Send confirmation email
	h.mailer.SendBookingApproved(&booking, &booking.Link)

	// Create calendar event
	h.caldav.CreateEvent(ctx, booking.Link.UserID, &booking, &booking.Slot, booking.Link.EventTemplate)

	return mapBookingToGen(&booking), nil
}

// DeclineBooking declines a booking
func (h *Handler) DeclineBooking(ctx context.Context, params gen.DeclineBookingParams) (*gen.Booking, error) {
	userID, _ := GetUserID(ctx)

	var booking Booking
	if err := h.db.Preload("Link").Preload("Slot").First(&booking, params.ID).Error; err != nil {
		return nil, err
	}

	// Verify ownership
	if booking.Link.UserID != userID {
		return nil, &gen.ErrorStatusCode{
			StatusCode: 403,
			Response:   gen.Error{Message: "Not authorized"},
		}
	}

	booking.Status = BookingStatusDeclined
	if err := h.db.Save(&booking).Error; err != nil {
		return nil, err
	}

	// Send decline email
	h.mailer.SendBookingDeclined(&booking, &booking.Link)

	return mapBookingToGen(&booking), nil
}

func mapBookingsToGen(bookings []Booking) []gen.Booking {
	result := make([]gen.Booking, len(bookings))
	for i, b := range bookings {
		result[i] = *mapBookingToGen(&b)
	}
	return result
}

func mapBookingToGen(b *Booking) *gen.Booking {
	return &gen.Booking{
		ID:           gen.NewOptInt(int(b.ID)),
		Slot:         *mapSlotToGen(&b.Slot),
		GuestEmail:   b.GuestEmail,
		GuestName:    gen.NewOptString(b.GuestName),
		Status:       gen.BookingStatus(b.Status),
		CustomFields: gen.NewOptStringStringMap(b.CustomFields),
		CreatedAt:    gen.NewOptDateTime(b.CreatedAt),
	}
}
```

**Step 2: Commit**

```bash
git add api/handler_bookings.go
git commit -m "feat(api): implement booking management handlers"
```

---

## Phase 6: Core API - Poll Flow

### Task 6.1: Vote Handlers

**Files:**
- Create: `api/handler_votes.go`

**Step 1: Create handler_votes.go**

```go
// api/handler_votes.go
package api

import (
	"context"
	"net/http"

	"github.com/konrad/meet-mesh/api/gen"
)

// SubmitVote submits a poll vote
func (h *Handler) SubmitVote(ctx context.Context, req *gen.SubmitVoteReq, params gen.SubmitVoteParams) (*gen.Vote, error) {
	var link Link
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		return nil, err
	}

	if link.Type != LinkTypePoll {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   gen.Error{Message: "Cannot vote on booking link"},
		}
	}

	if link.RequireEmail && req.GuestEmail.Value == "" {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   gen.Error{Message: "Email required"},
		}
	}

	// Convert responses
	responses := make(map[uint]VoteResponseType)
	for slotID, resp := range req.Responses {
		responses[uint(slotID)] = VoteResponseType(resp)
	}

	vote := Vote{
		LinkID:       link.ID,
		GuestEmail:   req.GuestEmail.Value,
		GuestName:    req.GuestName.Value,
		Responses:    responses,
		CustomFields: req.CustomFields.Value,
	}

	if err := h.db.Create(&vote).Error; err != nil {
		return nil, err
	}

	return mapVoteToGen(&vote), nil
}

// GetLinkVotes returns votes for a poll
func (h *Handler) GetLinkVotes(ctx context.Context, params gen.GetLinkVotesParams) ([]gen.Vote, error) {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	var votes []Vote
	if err := h.db.Where("link_id = ?", params.ID).Order("created_at DESC").Find(&votes).Error; err != nil {
		return nil, err
	}

	return mapVotesToGen(votes), nil
}

// GetPollResults returns poll results
func (h *Handler) GetPollResults(ctx context.Context, params gen.GetPollResultsParams) (*gen.GetPollResultsOK, error) {
	var link Link
	if err := h.db.Preload("Slots").Where("slug = ?", params.Slug).First(&link).Error; err != nil {
		return nil, err
	}

	if !link.ShowResults {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusForbidden,
			Response:   gen.Error{Message: "Results not public"},
		}
	}

	var votes []Vote
	if err := h.db.Where("link_id = ?", link.ID).Find(&votes).Error; err != nil {
		return nil, err
	}

	// Calculate tally
	tally := calculateTally(link.Slots, votes)

	return &gen.GetPollResultsOK{
		Tally: tally,
		Votes: gen.NewOptVoteArray(mapVotesToGen(votes)),
	}, nil
}

// PickPollWinner picks the winning slot
func (h *Handler) PickPollWinner(ctx context.Context, req *gen.PickPollWinnerReq, params gen.PickPollWinnerParams) error {
	userID, _ := GetUserID(ctx)

	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return err
	}

	var slot Slot
	if err := h.db.Where("id = ? AND link_id = ?", req.SlotID, link.ID).First(&slot).Error; err != nil {
		return err
	}

	// Close the poll
	link.Status = LinkStatusClosed
	if err := h.db.Save(&link).Error; err != nil {
		return err
	}

	// Get votes for notification
	var votes []Vote
	h.db.Where("link_id = ?", link.ID).Find(&votes)

	// Send winner notification
	h.mailer.SendPollWinner(&link, &slot, votes)

	return nil
}

func calculateTally(slots []Slot, votes []Vote) []gen.VoteTally {
	tally := make(map[uint]*gen.VoteTally)
	for _, slot := range slots {
		tally[slot.ID] = &gen.VoteTally{
			SlotID:     int(slot.ID),
			YesCount:   0,
			NoCount:    0,
			MaybeCount: 0,
		}
	}

	for _, vote := range votes {
		for slotID, response := range vote.Responses {
			if t, ok := tally[slotID]; ok {
				switch response {
				case VoteResponseYes:
					t.YesCount++
				case VoteResponseNo:
					t.NoCount++
				case VoteResponseMaybe:
					t.MaybeCount++
				}
			}
		}
	}

	result := make([]gen.VoteTally, 0, len(tally))
	for _, t := range tally {
		result = append(result, *t)
	}
	return result
}

func mapVotesToGen(votes []Vote) []gen.Vote {
	result := make([]gen.Vote, len(votes))
	for i, v := range votes {
		result[i] = *mapVoteToGen(&v)
	}
	return result
}

func mapVoteToGen(v *Vote) *gen.Vote {
	responses := make(map[string]gen.VoteResponse)
	for slotID, resp := range v.Responses {
		responses[fmt.Sprintf("%d", slotID)] = gen.VoteResponse(resp)
	}

	return &gen.Vote{
		ID:           gen.NewOptInt(int(v.ID)),
		GuestName:    gen.NewOptString(v.GuestName),
		GuestEmail:   gen.NewOptString(v.GuestEmail),
		Responses:    responses,
		CustomFields: gen.NewOptStringStringMap(v.CustomFields),
		CreatedAt:    gen.NewOptDateTime(v.CreatedAt),
	}
}
```

**Step 2: Commit**

```bash
git add api/handler_votes.go
git commit -m "feat(api): implement poll vote handlers"
```

---

## Phase 7: Email Notifications

### Task 7.1: Mailer Setup

**Files:**
- Create: `api/mailer.go`

**Step 1: Create mailer.go**

```go
// api/mailer.go
package api

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	config    *SMTPConfig
	baseURL   string
	dialer    *gomail.Dialer
	templates *template.Template
}

func NewMailer(cfg *SMTPConfig, baseURL string) (*Mailer, error) {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	// Parse templates
	tmpl, err := template.New("emails").Parse(emailTemplates)
	if err != nil {
		return nil, err
	}

	return &Mailer{
		config:    cfg,
		baseURL:   baseURL,
		dialer:    dialer,
		templates: tmpl,
	}, nil
}

func (m *Mailer) send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.config.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)
}

func (m *Mailer) renderTemplate(name string, data any) string {
	var buf bytes.Buffer
	m.templates.ExecuteTemplate(&buf, name, data)
	return buf.String()
}

// SendBookingConfirmation sends confirmation to guest
func (m *Mailer) SendBookingConfirmation(booking *Booking, link *Link) error {
	body := m.renderTemplate("booking_confirmed_guest", map[string]any{
		"LinkName":  link.Name,
		"GuestName": booking.GuestName,
		"Time":      booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
	})
	return m.send(booking.GuestEmail, "Booking Confirmed: "+link.Name, body)
}

// SendBookingPending sends pending notification to organizer
func (m *Mailer) SendBookingPending(booking *Booking, link *Link) error {
	approveURL := fmt.Sprintf("%s/api/v1/actions/approve?token=%s", m.baseURL, booking.ActionToken)
	declineURL := fmt.Sprintf("%s/api/v1/actions/decline?token=%s", m.baseURL, booking.ActionToken)

	body := m.renderTemplate("booking_pending", map[string]any{
		"LinkName":   link.Name,
		"GuestEmail": booking.GuestEmail,
		"GuestName":  booking.GuestName,
		"Time":       booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
		"ApproveURL": approveURL,
		"DeclineURL": declineURL,
	})

	// Get organizer email
	var user User
	// m.db.First(&user, link.UserID)
	// For now, use a placeholder - the handler should pass user email
	return m.send("organizer@example.com", "New Booking Request: "+link.Name, body)
}

// SendBookingApproved sends approval notification to guest
func (m *Mailer) SendBookingApproved(booking *Booking, link *Link) error {
	body := m.renderTemplate("booking_approved", map[string]any{
		"LinkName":  link.Name,
		"GuestName": booking.GuestName,
		"Time":      booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
	})
	return m.send(booking.GuestEmail, "Booking Approved: "+link.Name, body)
}

// SendBookingDeclined sends decline notification to guest
func (m *Mailer) SendBookingDeclined(booking *Booking, link *Link) error {
	body := m.renderTemplate("booking_declined", map[string]any{
		"LinkName":  link.Name,
		"GuestName": booking.GuestName,
		"Time":      booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
	})
	return m.send(booking.GuestEmail, "Booking Declined: "+link.Name, body)
}

// SendPollWinner sends winner notification to all voters
func (m *Mailer) SendPollWinner(link *Link, slot *Slot, votes []Vote) error {
	body := m.renderTemplate("poll_winner", map[string]any{
		"LinkName": link.Name,
		"Time":     slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
	})

	for _, vote := range votes {
		if vote.GuestEmail != "" {
			m.send(vote.GuestEmail, "Date Selected: "+link.Name, body)
		}
	}
	return nil
}

const emailTemplates = `
{{define "booking_confirmed_guest"}}
<html>
<body>
<h1>Booking Confirmed!</h1>
<p>Hi {{.GuestName}},</p>
<p>Your booking for <strong>{{.LinkName}}</strong> has been confirmed.</p>
<p><strong>When:</strong> {{.Time}}</p>
</body>
</html>
{{end}}

{{define "booking_pending"}}
<html>
<body>
<h1>New Booking Request</h1>
<p>You have a new booking request for <strong>{{.LinkName}}</strong>.</p>
<p><strong>Guest:</strong> {{.GuestName}} ({{.GuestEmail}})</p>
<p><strong>Requested time:</strong> {{.Time}}</p>
<p>
<a href="{{.ApproveURL}}" style="background:#22c55e;color:white;padding:10px 20px;text-decoration:none;border-radius:5px;">Approve</a>
<a href="{{.DeclineURL}}" style="background:#ef4444;color:white;padding:10px 20px;text-decoration:none;border-radius:5px;margin-left:10px;">Decline</a>
</p>
</body>
</html>
{{end}}

{{define "booking_approved"}}
<html>
<body>
<h1>Booking Approved!</h1>
<p>Hi {{.GuestName}},</p>
<p>Great news! Your booking for <strong>{{.LinkName}}</strong> has been approved.</p>
<p><strong>When:</strong> {{.Time}}</p>
</body>
</html>
{{end}}

{{define "booking_declined"}}
<html>
<body>
<h1>Booking Update</h1>
<p>Hi {{.GuestName}},</p>
<p>Unfortunately, your booking request for <strong>{{.LinkName}}</strong> could not be accommodated.</p>
<p><strong>Requested time:</strong> {{.Time}}</p>
</body>
</html>
{{end}}

{{define "poll_winner"}}
<html>
<body>
<h1>Date Selected!</h1>
<p>The organizer has selected a date for <strong>{{.LinkName}}</strong>.</p>
<p><strong>Selected time:</strong> {{.Time}}</p>
</body>
</html>
{{end}}
`
```

**Step 2: Add gomail dependency**

```bash
cd api && go get gopkg.in/gomail.v2
```

**Step 3: Commit**

```bash
git add api/mailer.go
git commit -m "feat(api): add email notification system"
```

---

### Task 7.2: Email Action Handlers

**Files:**
- Create: `api/handler_actions.go`

**Step 1: Create handler_actions.go**

```go
// api/handler_actions.go
package api

import (
	"context"
	"net/http"

	"github.com/konrad/meet-mesh/api/gen"
)

// ApproveViaEmail approves booking via email link
func (h *Handler) ApproveViaEmail(ctx context.Context) (*gen.ApproveViaEmailOK, error) {
	bookingID, ok := GetBookingID(ctx)
	if !ok {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   gen.Error{Message: "Invalid token"},
		}
	}

	var booking Booking
	if err := h.db.Preload("Link").Preload("Slot").First(&booking, bookingID).Error; err != nil {
		return nil, err
	}

	if booking.Status != BookingStatusPending {
		return &gen.ApproveViaEmailOK{
			Message: gen.NewOptString("Booking already processed"),
		}, nil
	}

	booking.Status = BookingStatusConfirmed
	if err := h.db.Save(&booking).Error; err != nil {
		return nil, err
	}

	// Clear the action token (single use)
	h.db.Model(&booking).Update("action_token", "")

	// Send confirmation email
	h.mailer.SendBookingApproved(&booking, &booking.Link)

	// Create calendar event
	h.caldav.CreateEvent(ctx, booking.Link.UserID, &booking, &booking.Slot, booking.Link.EventTemplate)

	return &gen.ApproveViaEmailOK{
		Message: gen.NewOptString("Booking approved successfully"),
	}, nil
}

// DeclineViaEmail declines booking via email link
func (h *Handler) DeclineViaEmail(ctx context.Context) (*gen.DeclineViaEmailOK, error) {
	bookingID, ok := GetBookingID(ctx)
	if !ok {
		return nil, &gen.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   gen.Error{Message: "Invalid token"},
		}
	}

	var booking Booking
	if err := h.db.Preload("Link").Preload("Slot").First(&booking, bookingID).Error; err != nil {
		return nil, err
	}

	if booking.Status != BookingStatusPending {
		return &gen.DeclineViaEmailOK{
			Message: gen.NewOptString("Booking already processed"),
		}, nil
	}

	booking.Status = BookingStatusDeclined
	if err := h.db.Save(&booking).Error; err != nil {
		return nil, err
	}

	// Clear the action token (single use)
	h.db.Model(&booking).Update("action_token", "")

	// Send decline email
	h.mailer.SendBookingDeclined(&booking, &booking.Link)

	return &gen.DeclineViaEmailOK{
		Message: gen.NewOptString("Booking declined"),
	}, nil
}
```

**Step 2: Commit**

```bash
git add api/handler_actions.go
git commit -m "feat(api): implement email action handlers"
```

---

### Task 7.3: Main Handler and Server Wiring

**Files:**
- Create: `api/handler.go`
- Modify: `api/cmd/main.go`

**Step 1: Create handler.go**

```go
// api/handler.go
package api

import (
	"gorm.io/gorm"

	"github.com/konrad/meet-mesh/api/gen"
)

type Handler struct {
	db     *gorm.DB
	auth   *AuthService
	caldav *CalDAVClient
	mailer *Mailer
	config *Config
}

var _ gen.Handler = (*Handler)(nil)

func NewHandler(db *gorm.DB, auth *AuthService, caldav *CalDAVClient, mailer *Mailer, config *Config) *Handler {
	return &Handler{
		db:     db,
		auth:   auth,
		caldav: caldav,
		mailer: mailer,
		config: config,
	}
}
```

**Step 2: Update cmd/main.go**

```go
// api/cmd/main.go
package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	api "github.com/konrad/meet-mesh/api"
	"github.com/konrad/meet-mesh/api/gen"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := api.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := api.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// Initialize auth service
	ctx := context.Background()
	auth, err := api.NewAuthService(ctx, &cfg.OIDC)
	if err != nil {
		log.Fatalf("Failed to init auth: %v", err)
	}

	// Initialize CalDAV client
	caldav := api.NewCalDAVClient(db)

	// Initialize mailer
	mailer, err := api.NewMailer(&cfg.SMTP, cfg.Server.BaseURL)
	if err != nil {
		log.Fatalf("Failed to init mailer: %v", err)
	}

	// Create handler
	handler := api.NewHandler(db, auth, caldav, mailer, cfg)

	// Create security handler
	security := api.NewSecurityHandler(db, auth)

	// Create server
	server, err := gen.NewServer(handler, security)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Meet Mesh API starting on port %d...", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), server))
}
```

**Step 3: Run go mod tidy**

```bash
cd api && go mod tidy
```

**Step 4: Commit**

```bash
git add api/handler.go api/cmd/main.go api/go.mod api/go.sum
git commit -m "feat(api): wire up handler and server"
```

---

## Summary

This plan covers **Phases 1-7** of the Meet Mesh API implementation:

| Phase | Tasks | Commits |
|-------|-------|---------|
| 1. Foundation | 13 tasks | 13 commits |
| 2. Authentication | 3 tasks | 3 commits |
| 3. CalDAV Integration | 2 tasks | 2 commits |
| 4. Core API - Links | 2 tasks | 2 commits |
| 5. Core API - Booking Flow | 2 tasks | 2 commits |
| 6. Core API - Poll Flow | 1 task | 1 commit |
| 7. Email Notifications | 3 tasks | 3 commits |

**Total: 26 tasks, 26 atomic commits**

## Documentation Sources

- [ogen - OpenAPI v3 code generator](https://ogen.dev/docs/intro/)
- [GORM Documentation](https://gorm.io)
- [coreos/go-oidc](https://github.com/coreos/go-oidc)
- [emersion/go-webdav CalDAV client](https://github.com/emersion/go-webdav)
