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
