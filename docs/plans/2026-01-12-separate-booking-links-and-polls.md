# Separate Booking Links and Polls Implementation Plan

**Goal:** Split the unified `Link` model into separate `BookingLink` and `Poll` entities with distinct database tables, API endpoints, handlers, and frontend components.

**Architecture:** Replace single-table inheritance with two independent entity types. BookingLinks handle 1:1 scheduling with CalDAV availability. Polls handle many-to-many voting on predefined options. Each gets dedicated API endpoints (`/booking-links/*` and `/polls/*`) and type-safe schemas.

**Tech Stack:** Go/GORM (backend), OpenAPI/ogen (API codegen), SvelteKit (frontend), SQLite (database)

---

## Phase 1: Database Schema Changes

### Task 1.1: Create BookingLink Model

**Files:**
- Modify: `api/models.go:101-120`

**Step 1: Add BookingLink struct after Link struct**

```go
type BookingLink struct {
	ID                uint               `gorm:"primaryKey"`
	UserID            uint               `gorm:"index;not null"`
	Slug              string             `gorm:"uniqueIndex;not null"`
	Name              string             `gorm:"not null"`
	Description       string
	Status            LinkStatus         `gorm:"not null;default:1"`
	AutoConfirm       bool
	AvailabilityRules []AvailabilityRule `gorm:"serializer:json"`
	RequireEmail      bool
	CustomFields      []CustomField      `gorm:"serializer:json"`
	EventTemplate     *EventTemplate     `gorm:"serializer:json"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Bookings          []Booking          `gorm:"foreignKey:BookingLinkID"`
}
```

**Step 2: Commit**

```bash
git add api/models.go
git commit -m "feat: add BookingLink model"
```

---

### Task 1.2: Create Poll Model

**Files:**
- Modify: `api/models.go`

**Step 1: Add Poll struct after BookingLink**

```go
type Poll struct {
	ID           uint          `gorm:"primaryKey"`
	UserID       uint          `gorm:"index;not null"`
	Slug         string        `gorm:"uniqueIndex;not null"`
	Name         string        `gorm:"not null"`
	Description  string
	Status       LinkStatus    `gorm:"not null;default:1"`
	ShowResults  bool
	RequireEmail bool
	CustomFields []CustomField `gorm:"serializer:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PollOptions  []PollOption  `gorm:"foreignKey:PollID"`
	Votes        []Vote        `gorm:"foreignKey:PollID"`
}
```

**Step 2: Commit**

```bash
git add api/models.go
git commit -m "feat: add Poll model"
```

---

### Task 1.3: Create PollOption Model (Replace Slot for Polls)

**Files:**
- Modify: `api/models.go`

**Step 1: Add PollOption struct**

```go
type PollOption struct {
	ID        uint      `gorm:"primaryKey"`
	PollID    uint      `gorm:"index;not null"`
	Type      SlotType  `gorm:"not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	CreatedAt time.Time
}
```

**Step 2: Commit**

```bash
git add api/models.go
git commit -m "feat: add PollOption model for poll time options"
```

---

### Task 1.4: Update Booking Model

**Files:**
- Modify: `api/models.go:132-146`

**Step 1: Update Booking to reference BookingLink instead of Link**

```go
type Booking struct {
	ID            uint              `gorm:"primaryKey"`
	BookingLinkID uint              `gorm:"index;not null"`
	SlotID        uint              `gorm:"index;not null"`
	GuestEmail    string            `gorm:"not null"`
	GuestName     string
	CustomFields  map[string]string `gorm:"serializer:json"`
	Status        BookingStatus     `gorm:"not null;default:1"`
	ActionToken   string            `gorm:"uniqueIndex"`
	CalendarUID   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	BookingLink   BookingLink       `gorm:"foreignKey:BookingLinkID"`
	Slot          Slot              `gorm:"foreignKey:SlotID"`
}
```

**Step 2: Commit**

```bash
git add api/models.go
git commit -m "refactor: update Booking to reference BookingLink"
```

---

### Task 1.5: Update Vote Model

**Files:**
- Modify: `api/models.go:148-157`

**Step 1: Update Vote to reference Poll and use PollOption IDs**

```go
type Vote struct {
	ID           uint                      `gorm:"primaryKey"`
	PollID       uint                      `gorm:"index;not null"`
	GuestEmail   string
	GuestName    string
	Responses    map[uint]VoteResponseType `gorm:"serializer:json;not null"`
	CustomFields map[string]string         `gorm:"serializer:json"`
	CreatedAt    time.Time
	Poll         Poll                      `gorm:"foreignKey:PollID"`
}
```

**Step 2: Commit**

```bash
git add api/models.go
git commit -m "refactor: update Vote to reference Poll"
```

---

### Task 1.6: Update AutoMigrate in main.go

**Files:**
- Modify: `api/cmd/main.go`

**Step 1: Find the AutoMigrate call and add new models**

```go
db.AutoMigrate(&User{}, &CalendarConnection{}, &BookingLink{}, &Poll{}, &PollOption{}, &Slot{}, &Booking{}, &Vote{})
```

**Step 2: Commit**

```bash
git add api/cmd/main.go
git commit -m "feat: add new models to AutoMigrate"
```

---

## Phase 2: OpenAPI Schema Updates

### Task 2.1: Add BookingLink Schema

**Files:**
- Modify: `api/openapi.yaml:131-166`

**Step 1: Add BookingLink schema after Link schema**

```yaml
    BookingLink:
      type: object
      required: [id, slug, name, status]
      properties:
        id:
          type: integer
        slug:
          type: string
        name:
          type: string
        description:
          type: string
        status:
          $ref: '#/components/schemas/LinkStatus'
        auto_confirm:
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
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat: add BookingLink schema to OpenAPI spec"
```

---

### Task 2.2: Add Poll Schema

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add Poll schema**

```yaml
    Poll:
      type: object
      required: [id, slug, name, status]
      properties:
        id:
          type: integer
        slug:
          type: string
        name:
          type: string
        description:
          type: string
        status:
          $ref: '#/components/schemas/LinkStatus'
        show_results:
          type: boolean
        require_email:
          type: boolean
        custom_fields:
          type: array
          items:
            $ref: '#/components/schemas/CustomField'
        created_at:
          type: string
          format: date-time
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat: add Poll schema to OpenAPI spec"
```

---

### Task 2.3: Add PollOption Schema

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add PollOption schema**

```yaml
    PollOption:
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
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat: add PollOption schema to OpenAPI spec"
```

---

### Task 2.4: Add BookingLink API Endpoints

**Files:**
- Modify: `api/openapi.yaml:386-527`

**Step 1: Add /booking-links endpoints (copy and adapt from /links)**

```yaml
  /booking-links:
    get:
      operationId: listBookingLinks
      summary: List all booking links
      security:
        - cookieAuth: []
      responses:
        '200':
          description: Booking links
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/BookingLink'

    post:
      operationId: createBookingLink
      summary: Create a booking link
      security:
        - cookieAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [name]
              properties:
                name:
                  type: string
                description:
                  type: string
                auto_confirm:
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
          description: Booking link created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BookingLink'

  /booking-links/{id}:
    get:
      operationId: getBookingLink
      summary: Get booking link details
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
          description: Booking link details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BookingLink'

    put:
      operationId: updateBookingLink
      summary: Update a booking link
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
          description: Booking link updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BookingLink'

    delete:
      operationId: deleteBookingLink
      summary: Delete a booking link
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
          description: Booking link deleted

  /booking-links/{id}/bookings:
    get:
      operationId: getBookingLinkBookings
      summary: Get bookings for a booking link
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
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat: add BookingLink API endpoints to OpenAPI spec"
```

---

### Task 2.5: Add Poll API Endpoints

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Add /polls endpoints**

```yaml
  /polls:
    get:
      operationId: listPolls
      summary: List all polls
      security:
        - cookieAuth: []
      responses:
        '200':
          description: Polls
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Poll'

    post:
      operationId: createPoll
      summary: Create a poll
      security:
        - cookieAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [name]
              properties:
                name:
                  type: string
                description:
                  type: string
                show_results:
                  type: boolean
                require_email:
                  type: boolean
                custom_fields:
                  type: array
                  items:
                    $ref: '#/components/schemas/CustomField'
      responses:
        '201':
          description: Poll created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Poll'

  /polls/{id}:
    get:
      operationId: getPoll
      summary: Get poll details
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
          description: Poll details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Poll'

    put:
      operationId: updatePoll
      summary: Update a poll
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
                show_results:
                  type: boolean
                require_email:
                  type: boolean
                custom_fields:
                  type: array
                  items:
                    $ref: '#/components/schemas/CustomField'
      responses:
        '200':
          description: Poll updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Poll'

    delete:
      operationId: deletePoll
      summary: Delete a poll
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
          description: Poll deleted

  /polls/{id}/options:
    get:
      operationId: getPollOptions
      summary: Get options for a poll
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
          description: Poll options
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/PollOption'

    post:
      operationId: addPollOption
      summary: Add an option to a poll
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
          description: Option added
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PollOption'

  /polls/{id}/options/{optionId}:
    delete:
      operationId: deletePollOption
      summary: Delete an option from a poll
      security:
        - cookieAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
        - name: optionId
          in: path
          required: true
          schema:
            type: integer
      responses:
        '204':
          description: Option deleted

  /polls/{id}/votes:
    get:
      operationId: getPollVotes
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

  /polls/{id}/pick-winner:
    post:
      operationId: pickPollWinner
      summary: Pick winning option for poll
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
              required: [option_id]
              properties:
                option_id:
                  type: integer
      responses:
        '200':
          description: Winner picked
```

**Step 2: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat: add Poll API endpoints to OpenAPI spec"
```

---

### Task 2.6: Update Public Endpoints

**Files:**
- Modify: `api/openapi.yaml:716-916`

**Step 1: Add separate public endpoints for booking links**

```yaml
  /p/booking/{slug}:
    get:
      operationId: getPublicBookingLink
      summary: Get public booking link info
      parameters:
        - name: slug
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Booking link info
          content:
            application/json:
              schema:
                type: object
                required: [name]
                properties:
                  name:
                    type: string
                  description:
                    type: string
                  custom_fields:
                    type: array
                    items:
                      $ref: '#/components/schemas/CustomField'
                  require_email:
                    type: boolean
        '404':
          description: Not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /p/booking/{slug}/availability:
    get:
      operationId: getBookingAvailability
      summary: Get real-time availability for booking link
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

  /p/booking/{slug}/book:
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
```

**Step 2: Add separate public endpoints for polls**

```yaml
  /p/poll/{slug}:
    get:
      operationId: getPublicPoll
      summary: Get public poll info
      parameters:
        - name: slug
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Poll info
          content:
            application/json:
              schema:
                type: object
                required: [name, options]
                properties:
                  name:
                    type: string
                  description:
                    type: string
                  custom_fields:
                    type: array
                    items:
                      $ref: '#/components/schemas/CustomField'
                  options:
                    type: array
                    items:
                      $ref: '#/components/schemas/PollOption'
                  show_results:
                    type: boolean
                  require_email:
                    type: boolean
        '404':
          description: Not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /p/poll/{slug}/vote:
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

  /p/poll/{slug}/results:
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

**Step 3: Commit**

```bash
git add api/openapi.yaml
git commit -m "feat: add separate public endpoints for booking links and polls"
```

---

### Task 2.7: Remove Old Link Schema and Endpoints

**Files:**
- Modify: `api/openapi.yaml`

**Step 1: Remove LinkType enum (no longer needed)**

Delete:
```yaml
    LinkType:
      type: integer
      enum: [1, 2]
      description: "1=booking, 2=poll"
```

**Step 2: Remove old Link schema (lines 131-166)**

**Step 3: Remove old /links endpoints (lines 386-649)**

**Step 4: Remove old /p/{slug} endpoints (lines 716-916)**

**Step 5: Run code generation**

```bash
make generate
```

**Step 6: Commit**

```bash
git add api/openapi.yaml api/gen/
git commit -m "refactor: remove unified Link schema and endpoints"
```

---

## Phase 3: Backend Handler Implementation

### Task 3.1: Create handler_booking_links.go

**Files:**
- Create: `api/handler_booking_links.go`

**Step 1: Create the file with CRUD handlers**

```go
package api

import (
	"context"
	"strings"

	"github.com/gosimple/slug"
	"meet-mesh/api/gen"
)

func (s *Server) ListBookingLinks(ctx context.Context) ([]gen.BookingLink, error) {
	user := ctx.Value(userContextKey).(*User)
	var links []BookingLink
	if err := s.db.Where("user_id = ?", user.ID).Find(&links).Error; err != nil {
		return nil, err
	}
	result := make([]gen.BookingLink, len(links))
	for i, link := range links {
		result[i] = bookingLinkToGen(link)
	}
	return result, nil
}

func (s *Server) CreateBookingLink(ctx context.Context, req *gen.CreateBookingLinkReq) (*gen.BookingLink, error) {
	user := ctx.Value(userContextKey).(*User)

	linkSlug := slug.Make(req.Name)
	// Ensure unique slug
	var count int64
	s.db.Model(&BookingLink{}).Where("slug LIKE ?", linkSlug+"%").Count(&count)
	if count > 0 {
		linkSlug = fmt.Sprintf("%s-%d", linkSlug, count+1)
	}

	link := BookingLink{
		UserID:       user.ID,
		Slug:         linkSlug,
		Name:         req.Name,
		Description:  req.Description.Or(""),
		Status:       LinkStatusActive,
		AutoConfirm:  req.AutoConfirm.Or(false),
		RequireEmail: req.RequireEmail.Or(true),
	}

	if rules, ok := req.AvailabilityRules.Get(); ok {
		link.AvailabilityRules = make([]AvailabilityRule, len(rules))
		for i, r := range rules {
			link.AvailabilityRules[i] = AvailabilityRule{
				DaysOfWeek: r.DaysOfWeek,
				StartTime:  r.StartTime,
				EndTime:    r.EndTime,
			}
		}
	}

	if fields, ok := req.CustomFields.Get(); ok {
		link.CustomFields = customFieldsFromGen(fields)
	}

	if tmpl, ok := req.EventTemplate.Get(); ok {
		link.EventTemplate = eventTemplateFromGen(tmpl)
	}

	if err := s.db.Create(&link).Error; err != nil {
		return nil, err
	}

	result := bookingLinkToGen(link)
	return &result, nil
}

func (s *Server) GetBookingLink(ctx context.Context, params gen.GetBookingLinkParams) (*gen.BookingLink, error) {
	user := ctx.Value(userContextKey).(*User)
	var link BookingLink
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&link).Error; err != nil {
		return nil, err
	}
	result := bookingLinkToGen(link)
	return &result, nil
}

func (s *Server) UpdateBookingLink(ctx context.Context, req *gen.UpdateBookingLinkReq, params gen.UpdateBookingLinkParams) (*gen.BookingLink, error) {
	user := ctx.Value(userContextKey).(*User)
	var link BookingLink
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&link).Error; err != nil {
		return nil, err
	}

	if name, ok := req.Name.Get(); ok {
		link.Name = name
	}
	if desc, ok := req.Description.Get(); ok {
		link.Description = desc
	}
	if status, ok := req.Status.Get(); ok {
		link.Status = LinkStatus(status)
	}
	if autoConfirm, ok := req.AutoConfirm.Get(); ok {
		link.AutoConfirm = autoConfirm
	}
	if requireEmail, ok := req.RequireEmail.Get(); ok {
		link.RequireEmail = requireEmail
	}
	if rules, ok := req.AvailabilityRules.Get(); ok {
		link.AvailabilityRules = make([]AvailabilityRule, len(rules))
		for i, r := range rules {
			link.AvailabilityRules[i] = AvailabilityRule{
				DaysOfWeek: r.DaysOfWeek,
				StartTime:  r.StartTime,
				EndTime:    r.EndTime,
			}
		}
	}
	if fields, ok := req.CustomFields.Get(); ok {
		link.CustomFields = customFieldsFromGen(fields)
	}
	if tmpl, ok := req.EventTemplate.Get(); ok {
		link.EventTemplate = eventTemplateFromGen(tmpl)
	}

	if err := s.db.Save(&link).Error; err != nil {
		return nil, err
	}

	result := bookingLinkToGen(link)
	return &result, nil
}

func (s *Server) DeleteBookingLink(ctx context.Context, params gen.DeleteBookingLinkParams) error {
	user := ctx.Value(userContextKey).(*User)
	return s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).Delete(&BookingLink{}).Error
}

func (s *Server) GetBookingLinkBookings(ctx context.Context, params gen.GetBookingLinkBookingsParams) ([]gen.Booking, error) {
	user := ctx.Value(userContextKey).(*User)
	var link BookingLink
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&link).Error; err != nil {
		return nil, err
	}

	var bookings []Booking
	if err := s.db.Preload("Slot").Where("booking_link_id = ?", link.ID).Find(&bookings).Error; err != nil {
		return nil, err
	}

	result := make([]gen.Booking, len(bookings))
	for i, b := range bookings {
		result[i] = bookingToGen(b)
	}
	return result, nil
}

func bookingLinkToGen(link BookingLink) gen.BookingLink {
	result := gen.BookingLink{
		ID:          int(link.ID),
		Slug:        link.Slug,
		Name:        link.Name,
		Status:      gen.LinkStatus(link.Status),
		AutoConfirm: gen.NewOptBool(link.AutoConfirm),
		RequireEmail: gen.NewOptBool(link.RequireEmail),
		CreatedAt:   gen.NewOptDateTime(link.CreatedAt),
	}
	if link.Description != "" {
		result.Description = gen.NewOptString(link.Description)
	}
	if len(link.AvailabilityRules) > 0 {
		rules := make([]gen.AvailabilityRule, len(link.AvailabilityRules))
		for i, r := range link.AvailabilityRules {
			rules[i] = gen.AvailabilityRule{
				DaysOfWeek: r.DaysOfWeek,
				StartTime:  r.StartTime,
				EndTime:    r.EndTime,
			}
		}
		result.AvailabilityRules = gen.NewOptNilAvailabilityRuleArray(rules)
	}
	if len(link.CustomFields) > 0 {
		result.CustomFields = gen.NewOptNilCustomFieldArray(customFieldsToGen(link.CustomFields))
	}
	if link.EventTemplate != nil {
		result.EventTemplate = gen.NewOptEventTemplate(eventTemplateToGen(*link.EventTemplate))
	}
	return result
}
```

**Step 2: Commit**

```bash
git add api/handler_booking_links.go
git commit -m "feat: add BookingLink CRUD handlers"
```

---

### Task 3.2: Create handler_polls.go

**Files:**
- Create: `api/handler_polls.go`

**Step 1: Create the file with Poll CRUD and PollOption handlers**

```go
package api

import (
	"context"

	"github.com/gosimple/slug"
	"meet-mesh/api/gen"
)

func (s *Server) ListPolls(ctx context.Context) ([]gen.Poll, error) {
	user := ctx.Value(userContextKey).(*User)
	var polls []Poll
	if err := s.db.Where("user_id = ?", user.ID).Find(&polls).Error; err != nil {
		return nil, err
	}
	result := make([]gen.Poll, len(polls))
	for i, poll := range polls {
		result[i] = pollToGen(poll)
	}
	return result, nil
}

func (s *Server) CreatePoll(ctx context.Context, req *gen.CreatePollReq) (*gen.Poll, error) {
	user := ctx.Value(userContextKey).(*User)

	pollSlug := slug.Make(req.Name)
	var count int64
	s.db.Model(&Poll{}).Where("slug LIKE ?", pollSlug+"%").Count(&count)
	if count > 0 {
		pollSlug = fmt.Sprintf("%s-%d", pollSlug, count+1)
	}

	poll := Poll{
		UserID:       user.ID,
		Slug:         pollSlug,
		Name:         req.Name,
		Description:  req.Description.Or(""),
		Status:       LinkStatusActive,
		ShowResults:  req.ShowResults.Or(true),
		RequireEmail: req.RequireEmail.Or(false),
	}

	if fields, ok := req.CustomFields.Get(); ok {
		poll.CustomFields = customFieldsFromGen(fields)
	}

	if err := s.db.Create(&poll).Error; err != nil {
		return nil, err
	}

	result := pollToGen(poll)
	return &result, nil
}

func (s *Server) GetPoll(ctx context.Context, params gen.GetPollParams) (*gen.Poll, error) {
	user := ctx.Value(userContextKey).(*User)
	var poll Poll
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&poll).Error; err != nil {
		return nil, err
	}
	result := pollToGen(poll)
	return &result, nil
}

func (s *Server) UpdatePoll(ctx context.Context, req *gen.UpdatePollReq, params gen.UpdatePollParams) (*gen.Poll, error) {
	user := ctx.Value(userContextKey).(*User)
	var poll Poll
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&poll).Error; err != nil {
		return nil, err
	}

	if name, ok := req.Name.Get(); ok {
		poll.Name = name
	}
	if desc, ok := req.Description.Get(); ok {
		poll.Description = desc
	}
	if status, ok := req.Status.Get(); ok {
		poll.Status = LinkStatus(status)
	}
	if showResults, ok := req.ShowResults.Get(); ok {
		poll.ShowResults = showResults
	}
	if requireEmail, ok := req.RequireEmail.Get(); ok {
		poll.RequireEmail = requireEmail
	}
	if fields, ok := req.CustomFields.Get(); ok {
		poll.CustomFields = customFieldsFromGen(fields)
	}

	if err := s.db.Save(&poll).Error; err != nil {
		return nil, err
	}

	result := pollToGen(poll)
	return &result, nil
}

func (s *Server) DeletePoll(ctx context.Context, params gen.DeletePollParams) error {
	user := ctx.Value(userContextKey).(*User)
	return s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).Delete(&Poll{}).Error
}

func (s *Server) GetPollOptions(ctx context.Context, params gen.GetPollOptionsParams) ([]gen.PollOption, error) {
	user := ctx.Value(userContextKey).(*User)
	var poll Poll
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&poll).Error; err != nil {
		return nil, err
	}

	var options []PollOption
	if err := s.db.Where("poll_id = ?", poll.ID).Order("start_time").Find(&options).Error; err != nil {
		return nil, err
	}

	result := make([]gen.PollOption, len(options))
	for i, opt := range options {
		result[i] = pollOptionToGen(opt)
	}
	return result, nil
}

func (s *Server) AddPollOption(ctx context.Context, req *gen.AddPollOptionReq, params gen.AddPollOptionParams) (*gen.PollOption, error) {
	user := ctx.Value(userContextKey).(*User)
	var poll Poll
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&poll).Error; err != nil {
		return nil, err
	}

	option := PollOption{
		PollID:    poll.ID,
		Type:      SlotType(req.Type),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := s.db.Create(&option).Error; err != nil {
		return nil, err
	}

	result := pollOptionToGen(option)
	return &result, nil
}

func (s *Server) DeletePollOption(ctx context.Context, params gen.DeletePollOptionParams) error {
	user := ctx.Value(userContextKey).(*User)
	var poll Poll
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&poll).Error; err != nil {
		return err
	}
	return s.db.Where("id = ? AND poll_id = ?", params.OptionID, poll.ID).Delete(&PollOption{}).Error
}

func (s *Server) GetPollVotes(ctx context.Context, params gen.GetPollVotesParams) ([]gen.Vote, error) {
	user := ctx.Value(userContextKey).(*User)
	var poll Poll
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&poll).Error; err != nil {
		return nil, err
	}

	var votes []Vote
	if err := s.db.Where("poll_id = ?", poll.ID).Find(&votes).Error; err != nil {
		return nil, err
	}

	result := make([]gen.Vote, len(votes))
	for i, v := range votes {
		result[i] = voteToGen(v)
	}
	return result, nil
}

func (s *Server) PickPollWinner(ctx context.Context, req *gen.PickPollWinnerReq, params gen.PickPollWinnerParams) error {
	user := ctx.Value(userContextKey).(*User)
	var poll Poll
	if err := s.db.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&poll).Error; err != nil {
		return err
	}

	poll.Status = LinkStatusClosed
	return s.db.Save(&poll).Error
}

func pollToGen(poll Poll) gen.Poll {
	result := gen.Poll{
		ID:           int(poll.ID),
		Slug:         poll.Slug,
		Name:         poll.Name,
		Status:       gen.LinkStatus(poll.Status),
		ShowResults:  gen.NewOptBool(poll.ShowResults),
		RequireEmail: gen.NewOptBool(poll.RequireEmail),
		CreatedAt:    gen.NewOptDateTime(poll.CreatedAt),
	}
	if poll.Description != "" {
		result.Description = gen.NewOptString(poll.Description)
	}
	if len(poll.CustomFields) > 0 {
		result.CustomFields = gen.NewOptNilCustomFieldArray(customFieldsToGen(poll.CustomFields))
	}
	return result
}

func pollOptionToGen(opt PollOption) gen.PollOption {
	return gen.PollOption{
		ID:        int(opt.ID),
		Type:      gen.SlotType(opt.Type),
		StartTime: opt.StartTime,
		EndTime:   opt.EndTime,
	}
}
```

**Step 2: Commit**

```bash
git add api/handler_polls.go
git commit -m "feat: add Poll and PollOption handlers"
```

---

### Task 3.3: Create handler_public_booking.go

**Files:**
- Create: `api/handler_public_booking.go`

**Step 1: Create public booking link handlers (adapted from handler_public.go)**

Copy and adapt the availability and booking logic from `handler_public.go` to work with BookingLink model.

**Step 2: Commit**

```bash
git add api/handler_public_booking.go
git commit -m "feat: add public booking link handlers"
```

---

### Task 3.4: Create handler_public_poll.go

**Files:**
- Create: `api/handler_public_poll.go`

**Step 1: Create public poll handlers (adapted from handler_public.go and handler_votes.go)**

Copy and adapt the voting and results logic to work with Poll model.

**Step 2: Commit**

```bash
git add api/handler_public_poll.go
git commit -m "feat: add public poll handlers"
```

---

### Task 3.5: Update handler_bookings.go

**Files:**
- Modify: `api/handler_bookings.go`

**Step 1: Update to reference BookingLink instead of Link**

Update the approve/decline handlers to work with the new BookingLink model.

**Step 2: Commit**

```bash
git add api/handler_bookings.go
git commit -m "refactor: update booking handlers for BookingLink model"
```

---

### Task 3.6: Remove Old Handlers

**Files:**
- Delete: `api/handler_links.go`
- Delete: `api/handler_slots.go` (slots now managed differently)
- Modify: `api/handler_public.go` (remove all, now split)
- Delete: `api/handler_votes.go` (merged into polls)

**Step 1: Delete old files**

```bash
rm api/handler_links.go api/handler_slots.go api/handler_votes.go
```

**Step 2: Commit**

```bash
git add -A
git commit -m "refactor: remove old unified link handlers"
```

---

### Task 3.7: Remove LinkType from models.go

**Files:**
- Modify: `api/models.go:9-14`

**Step 1: Remove LinkType enum**

Delete:
```go
type LinkType int

const (
	LinkTypeBooking LinkType = 1
	LinkTypePoll    LinkType = 2
)
```

**Step 2: Remove old Link model (lines 101-120)**

**Step 3: Commit**

```bash
git add api/models.go
git commit -m "refactor: remove unified Link model and LinkType enum"
```

---

## Phase 4: Frontend Updates

### Task 4.1: Update API Types

**Files:**
- Run: `cd frontend && pnpm generate:api`

**Step 1: Regenerate TypeScript types**

```bash
cd frontend && pnpm generate:api
```

**Step 2: Commit**

```bash
git add frontend/src/lib/api/generated/
git commit -m "chore: regenerate frontend API types"
```

---

### Task 4.2: Create Separate Dashboard Routes

**Files:**
- Create: `frontend/src/routes/(dashboard)/booking-links/+page.svelte`
- Create: `frontend/src/routes/(dashboard)/booking-links/new/+page.svelte`
- Create: `frontend/src/routes/(dashboard)/booking-links/[id]/+page.svelte`
- Create: `frontend/src/routes/(dashboard)/booking-links/[id]/edit/+page.svelte`
- Create: `frontend/src/routes/(dashboard)/polls/+page.svelte`
- Create: `frontend/src/routes/(dashboard)/polls/new/+page.svelte`
- Create: `frontend/src/routes/(dashboard)/polls/[id]/+page.svelte`
- Create: `frontend/src/routes/(dashboard)/polls/[id]/edit/+page.svelte`

**Step 1: Create booking-links routes (adapt from existing links routes)**

**Step 2: Create polls routes (adapt from existing links routes)**

**Step 3: Commit**

```bash
git add frontend/src/routes/\(dashboard\)/booking-links/ frontend/src/routes/\(dashboard\)/polls/
git commit -m "feat: add separate dashboard routes for booking links and polls"
```

---

### Task 4.3: Update Dashboard Index

**Files:**
- Modify: `frontend/src/routes/(dashboard)/+page.svelte`

**Step 1: Update to show separate sections for booking links and polls**

Replace single links list with two separate sections/cards.

**Step 2: Commit**

```bash
git add frontend/src/routes/\(dashboard\)/+page.svelte
git commit -m "feat: update dashboard to show booking links and polls separately"
```

---

### Task 4.4: Create Separate Public Routes

**Files:**
- Create: `frontend/src/routes/(public)/booking/[slug]/+page.svelte`
- Create: `frontend/src/routes/(public)/poll/[slug]/+page.svelte`

**Step 1: Create booking page (uses BookingPage component)**

```svelte
<script lang="ts">
	import { page } from '$app/stores';
	import { api } from '$lib/api/client';
	import { BookingPage } from '$lib/components/booking';
	import { Card } from '$lib/components/ui';

	const slug = $page.params.slug;
	let link = $state<Awaited<ReturnType<typeof loadLink>>>(null);
	let error = $state('');

	async function loadLink() {
		const { data, error: apiError } = await api.GET('/p/booking/{slug}', {
			params: { path: { slug } }
		});
		if (apiError) {
			error = 'Booking link not found';
			return null;
		}
		return data;
	}

	$effect(() => {
		loadLink().then((data) => (link = data));
	});
</script>

{#if error}
	<Card>
		<p class="text-red-600">{error}</p>
	</Card>
{:else if link}
	<BookingPage {link} {slug} />
{:else}
	<Card>Loading...</Card>
{/if}
```

**Step 2: Create poll page (uses PollPage component)**

```svelte
<script lang="ts">
	import { page } from '$app/stores';
	import { api } from '$lib/api/client';
	import { PollPage } from '$lib/components/poll';
	import { Card } from '$lib/components/ui';

	const slug = $page.params.slug;
	let poll = $state<Awaited<ReturnType<typeof loadPoll>>>(null);
	let error = $state('');

	async function loadPoll() {
		const { data, error: apiError } = await api.GET('/p/poll/{slug}', {
			params: { path: { slug } }
		});
		if (apiError) {
			error = 'Poll not found';
			return null;
		}
		return data;
	}

	$effect(() => {
		loadPoll().then((data) => (poll = data));
	});
</script>

{#if error}
	<Card>
		<p class="text-red-600">{error}</p>
	</Card>
{:else if poll}
	<PollPage {poll} {slug} />
{:else}
	<Card>Loading...</Card>
{/if}
```

**Step 3: Commit**

```bash
git add frontend/src/routes/\(public\)/booking/ frontend/src/routes/\(public\)/poll/
git commit -m "feat: add separate public routes for booking and poll"
```

---

### Task 4.5: Update BookingPage Component

**Files:**
- Modify: `frontend/src/lib/components/booking/BookingPage.svelte`

**Step 1: Update props to accept BookingLink type instead of Link**

**Step 2: Update API calls to use new endpoints**

**Step 3: Commit**

```bash
git add frontend/src/lib/components/booking/BookingPage.svelte
git commit -m "refactor: update BookingPage for new BookingLink type"
```

---

### Task 4.6: Update PollPage Component

**Files:**
- Modify: `frontend/src/lib/components/poll/PollPage.svelte`

**Step 1: Update props to accept Poll type instead of Link**

**Step 2: Update to use `options` instead of `slots`**

**Step 3: Update API calls to use new endpoints**

**Step 4: Commit**

```bash
git add frontend/src/lib/components/poll/PollPage.svelte
git commit -m "refactor: update PollPage for new Poll type"
```

---

### Task 4.7: Remove Old Routes

**Files:**
- Delete: `frontend/src/routes/(dashboard)/links/`
- Delete: `frontend/src/routes/(public)/p/`

**Step 1: Remove old unified routes**

```bash
rm -rf frontend/src/routes/\(dashboard\)/links/
rm -rf frontend/src/routes/\(public\)/p/
```

**Step 2: Commit**

```bash
git add -A
git commit -m "refactor: remove old unified link routes"
```

---

## Phase 5: Data Migration

### Task 5.1: Create Migration Script

**Files:**
- Create: `api/cmd/migrate/main.go`

**Step 1: Create migration script to copy data from links table to booking_links and polls**

```go
package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate booking links (type = 1)
	db.Exec(`
		INSERT INTO booking_links (id, user_id, slug, name, description, status, auto_confirm, availability_rules, require_email, custom_fields, event_template, created_at, updated_at)
		SELECT id, user_id, slug, name, description, status, auto_confirm, availability_rules, require_email, custom_fields, event_template, created_at, updated_at
		FROM links WHERE type = 1
	`)

	// Migrate polls (type = 2)
	db.Exec(`
		INSERT INTO polls (id, user_id, slug, name, description, status, show_results, require_email, custom_fields, created_at, updated_at)
		SELECT id, user_id, slug, name, description, status, show_results, require_email, custom_fields, created_at, updated_at
		FROM links WHERE type = 2
	`)

	// Migrate slots to poll_options for polls
	db.Exec(`
		INSERT INTO poll_options (id, poll_id, type, start_time, end_time, created_at)
		SELECT s.id, s.link_id, s.type, s.start_time, s.end_time, s.created_at
		FROM slots s
		JOIN links l ON s.link_id = l.id
		WHERE l.type = 2
	`)

	// Update bookings to reference booking_link_id
	db.Exec(`
		UPDATE bookings SET booking_link_id = link_id
	`)

	// Update votes to reference poll_id
	db.Exec(`
		UPDATE votes SET poll_id = link_id
	`)

	log.Println("Migration complete")
}
```

**Step 2: Commit**

```bash
git add api/cmd/migrate/
git commit -m "feat: add data migration script for link separation"
```

---

### Task 5.2: Run Migration and Verify

**Step 1: Backup database**

```bash
cp data.db data.db.backup
```

**Step 2: Run migration**

```bash
go run api/cmd/migrate/main.go
```

**Step 3: Verify data**

```bash
sqlite3 data.db "SELECT COUNT(*) FROM booking_links"
sqlite3 data.db "SELECT COUNT(*) FROM polls"
```

**Step 4: Commit backup note**

```bash
git commit --allow-empty -m "chore: data migration completed"
```

---

## Phase 6: Testing and Cleanup

### Task 6.1: Build and Test

**Step 1: Build the project**

```bash
make build
```

**Step 2: Run the application**

```bash
./meet-mesh
```

**Step 3: Test booking link creation via dashboard**

**Step 4: Test poll creation via dashboard**

**Step 5: Test public booking flow**

**Step 6: Test public voting flow**

---

### Task 6.2: Drop Old Tables (After Verification)

**Files:**
- Modify: `api/cmd/migrate/main.go`

**Step 1: Add cleanup phase to migration**

```go
// After verifying new tables work:
db.Exec("DROP TABLE IF EXISTS links")
db.Exec("ALTER TABLE bookings DROP COLUMN link_id")
db.Exec("ALTER TABLE votes DROP COLUMN link_id")
```

**Step 2: Commit**

```bash
git add api/cmd/migrate/
git commit -m "chore: add cleanup for old link tables"
```

---

### Task 6.3: Update CLAUDE.md

**Files:**
- Modify: `CLAUDE.md`

**Step 1: Update architecture section to reflect new structure**

Update handler organization section:
```markdown
### Handler Organization

Handlers are split by domain:
- `handler_auth.go` - Authentication
- `handler_booking_links.go` - Booking link CRUD
- `handler_polls.go` - Poll CRUD and options
- `handler_bookings.go` - Booking approval/decline
- `handler_public_booking.go` - Public booking endpoints
- `handler_public_poll.go` - Public poll endpoints
- `handler_calendar.go` - CalDAV integration
- `handler_actions.go` - Email action handlers
```

**Step 2: Update data model documentation**

**Step 3: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: update CLAUDE.md for separated booking links and polls"
```

---

## Summary

| Phase | Tasks | Files Changed |
|-------|-------|---------------|
| 1. Database Schema | 6 tasks | `api/models.go`, `api/cmd/main.go` |
| 2. OpenAPI Schema | 7 tasks | `api/openapi.yaml`, `api/gen/*` |
| 3. Backend Handlers | 7 tasks | `api/handler_*.go` |
| 4. Frontend Updates | 7 tasks | `frontend/src/routes/**`, `frontend/src/lib/components/**` |
| 5. Data Migration | 2 tasks | `api/cmd/migrate/main.go` |
| 6. Testing/Cleanup | 3 tasks | Various |

**Total: 32 tasks**

**Key breaking changes:**
- `/links` endpoints removed, replaced with `/booking-links` and `/polls`
- `/p/{slug}` replaced with `/p/booking/{slug}` and `/p/poll/{slug}`
- `Link` model removed, replaced with `BookingLink` and `Poll`
- `Slot` used only for booking availability; polls use `PollOption`
- Frontend routes restructured accordingly
