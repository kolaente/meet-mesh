// api/models.go
package api

import (
	"time"
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
	ID           uint      `gorm:"primaryKey"`
	OIDCSub      string    `gorm:"column:oidc_sub;uniqueIndex;not null"`
	Email        string    `gorm:"not null"`
	Name         string
	CreatedAt    time.Time
	Calendars    []CalendarConnection `gorm:"foreignKey:UserID"`
	BookingLinks []BookingLink        `gorm:"foreignKey:UserID"`
	Polls        []Poll               `gorm:"foreignKey:UserID"`
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

type BookingLink struct {
	ID                   uint               `gorm:"primaryKey"`
	UserID               uint               `gorm:"index;not null"`
	Slug                 string             `gorm:"uniqueIndex;not null"`
	Name                 string             `gorm:"not null"`
	Description          string
	Status               LinkStatus         `gorm:"not null;default:1"`
	AutoConfirm          bool
	SlotDurationMinutes  int                `gorm:"not null;default:30"`
	SlotDurationsMinutes []int              `gorm:"serializer:json"`
	BufferMinutes        int                `gorm:"not null;default:0"`
	AvailabilityRules    []AvailabilityRule `gorm:"serializer:json"`
	RequireEmail         bool
	MeetingLink          string
	CustomFields         []CustomField      `gorm:"serializer:json"`
	EventTemplate        *EventTemplate     `gorm:"serializer:json"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Bookings             []Booking          `gorm:"foreignKey:BookingLinkID"`
}

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

type PollOption struct {
	ID        uint      `gorm:"primaryKey"`
	PollID    uint      `gorm:"index;not null"`
	Type      SlotType  `gorm:"not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type Slot struct {
	ID            uint      `gorm:"primaryKey"`
	BookingLinkID uint      `gorm:"index;not null;default:0"`
	Type          SlotType  `gorm:"not null"`
	StartTime     time.Time `gorm:"not null"`
	EndTime       time.Time `gorm:"not null"`
	Manual        bool
	CreatedAt     time.Time
}

type Booking struct {
	ID            uint              `gorm:"primaryKey"`
	BookingLinkID uint              `gorm:"index;not null;default:0"`
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

type Vote struct {
	ID           uint                      `gorm:"primaryKey"`
	PollID       uint                      `gorm:"index;not null;default:0"`
	GuestEmail   string
	GuestName    string
	Responses    map[uint]VoteResponseType `gorm:"serializer:json;not null"`
	CustomFields map[string]string         `gorm:"serializer:json"`
	CreatedAt    time.Time
	Poll         Poll                      `gorm:"foreignKey:PollID"`
}
