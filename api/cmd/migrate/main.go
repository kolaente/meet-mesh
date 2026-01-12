// cmd/migrate/main.go
// Migration script to split unified Link model into BookingLink and Poll
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Old models (for reading from existing database)
type OldLink struct {
	ID                uint
	UserID            uint
	Type              int // 1=booking, 2=poll
	Slug              string
	Name              string
	Description       string
	Status            int
	AutoConfirm       bool
	ShowResults       bool
	RequireEmail      bool
	AvailabilityRules []byte `gorm:"column:availability_rules"`
	CustomFields      []byte `gorm:"column:custom_fields"`
	EventTemplate     []byte `gorm:"column:event_template"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (OldLink) TableName() string {
	return "links"
}

type OldSlot struct {
	ID        uint
	LinkID    uint
	Type      int
	StartTime time.Time
	EndTime   time.Time
	Manual    bool
	CreatedAt time.Time
}

func (OldSlot) TableName() string {
	return "slots"
}

type OldBooking struct {
	ID           uint
	LinkID       uint
	SlotID       uint
	GuestEmail   string
	GuestName    string
	CustomFields []byte `gorm:"column:custom_fields"`
	Status       int
	ActionToken  string
	CalendarUID  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (OldBooking) TableName() string {
	return "bookings"
}

type OldVote struct {
	ID           uint
	LinkID       uint
	GuestEmail   string
	GuestName    string
	Responses    []byte `gorm:"column:responses"`
	CustomFields []byte `gorm:"column:custom_fields"`
	CreatedAt    time.Time
}

func (OldVote) TableName() string {
	return "votes"
}

// New models (for writing to new tables)
type BookingLink struct {
	ID                uint      `gorm:"primaryKey"`
	UserID            uint      `gorm:"index;not null"`
	Slug              string    `gorm:"uniqueIndex;not null"`
	Name              string    `gorm:"not null"`
	Description       string
	Status            int       `gorm:"not null;default:1"`
	AutoConfirm       bool
	AvailabilityRules []byte    `gorm:"serializer:json"`
	RequireEmail      bool
	CustomFields      []byte    `gorm:"serializer:json"`
	EventTemplate     []byte    `gorm:"serializer:json"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Poll struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index;not null"`
	Slug         string    `gorm:"uniqueIndex;not null"`
	Name         string    `gorm:"not null"`
	Description  string
	Status       int       `gorm:"not null;default:1"`
	ShowResults  bool
	RequireEmail bool
	CustomFields []byte    `gorm:"serializer:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PollOption struct {
	ID        uint      `gorm:"primaryKey"`
	PollID    uint      `gorm:"index;not null"`
	Type      int       `gorm:"not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type Slot struct {
	ID            uint      `gorm:"primaryKey"`
	BookingLinkID uint      `gorm:"index;not null"`
	Type          int       `gorm:"not null"`
	StartTime     time.Time `gorm:"not null"`
	EndTime       time.Time `gorm:"not null"`
	Manual        bool
	CreatedAt     time.Time
}

type Booking struct {
	ID            uint
	BookingLinkID uint `gorm:"column:booking_link_id"`
	SlotID        uint
	GuestEmail    string
	GuestName     string
	CustomFields  []byte `gorm:"column:custom_fields"`
	Status        int
	ActionToken   string
	CalendarUID   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Vote struct {
	ID           uint
	PollID       uint `gorm:"column:poll_id"`
	GuestEmail   string
	GuestName    string
	Responses    []byte `gorm:"column:responses"`
	CustomFields []byte `gorm:"column:custom_fields"`
	CreatedAt    time.Time
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <database-path>")
		fmt.Println("Example: migrate ./data.db")
		os.Exit(1)
	}

	dbPath := os.Args[1]

	// Open database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Check if migration is needed
	if !tableExists(db, "links") {
		fmt.Println("No 'links' table found - migration not needed or already completed")
		return
	}

	fmt.Println("Starting migration...")

	// Create new tables
	fmt.Println("Creating new tables...")
	if err := db.AutoMigrate(&BookingLink{}, &Poll{}, &PollOption{}); err != nil {
		log.Fatal("Failed to create new tables:", err)
	}

	// Map old link IDs to new IDs
	linkToBookingLink := make(map[uint]uint)
	linkToPoll := make(map[uint]uint)
	slotToPollOption := make(map[uint]uint)

	// Migrate links
	var oldLinks []OldLink
	if err := db.Find(&oldLinks).Error; err != nil {
		log.Fatal("Failed to read old links:", err)
	}

	fmt.Printf("Found %d links to migrate\n", len(oldLinks))

	for _, oldLink := range oldLinks {
		if oldLink.Type == 1 {
			// Booking link
			newLink := BookingLink{
				UserID:            oldLink.UserID,
				Slug:              oldLink.Slug,
				Name:              oldLink.Name,
				Description:       oldLink.Description,
				Status:            oldLink.Status,
				AutoConfirm:       oldLink.AutoConfirm,
				AvailabilityRules: oldLink.AvailabilityRules,
				RequireEmail:      oldLink.RequireEmail,
				CustomFields:      oldLink.CustomFields,
				EventTemplate:     oldLink.EventTemplate,
				CreatedAt:         oldLink.CreatedAt,
				UpdatedAt:         oldLink.UpdatedAt,
			}
			if err := db.Create(&newLink).Error; err != nil {
				log.Printf("Warning: Failed to create booking link for %s: %v", oldLink.Slug, err)
				continue
			}
			linkToBookingLink[oldLink.ID] = newLink.ID
			fmt.Printf("  Migrated booking link: %s (old ID: %d -> new ID: %d)\n", oldLink.Slug, oldLink.ID, newLink.ID)
		} else {
			// Poll
			newPoll := Poll{
				UserID:       oldLink.UserID,
				Slug:         oldLink.Slug,
				Name:         oldLink.Name,
				Description:  oldLink.Description,
				Status:       oldLink.Status,
				ShowResults:  oldLink.ShowResults,
				RequireEmail: oldLink.RequireEmail,
				CustomFields: oldLink.CustomFields,
				CreatedAt:    oldLink.CreatedAt,
				UpdatedAt:    oldLink.UpdatedAt,
			}
			if err := db.Create(&newPoll).Error; err != nil {
				log.Printf("Warning: Failed to create poll for %s: %v", oldLink.Slug, err)
				continue
			}
			linkToPoll[oldLink.ID] = newPoll.ID
			fmt.Printf("  Migrated poll: %s (old ID: %d -> new ID: %d)\n", oldLink.Slug, oldLink.ID, newPoll.ID)
		}
	}

	// Migrate slots
	var oldSlots []OldSlot
	if err := db.Find(&oldSlots).Error; err != nil {
		log.Fatal("Failed to read old slots:", err)
	}

	fmt.Printf("Found %d slots to migrate\n", len(oldSlots))

	for _, oldSlot := range oldSlots {
		if newPollID, ok := linkToPoll[oldSlot.LinkID]; ok {
			// This slot belongs to a poll - create PollOption
			newOption := PollOption{
				PollID:    newPollID,
				Type:      oldSlot.Type,
				StartTime: oldSlot.StartTime,
				EndTime:   oldSlot.EndTime,
				CreatedAt: oldSlot.CreatedAt,
			}
			if err := db.Create(&newOption).Error; err != nil {
				log.Printf("Warning: Failed to create poll option: %v", err)
				continue
			}
			slotToPollOption[oldSlot.ID] = newOption.ID
			fmt.Printf("  Migrated slot to poll option: %d -> %d\n", oldSlot.ID, newOption.ID)
		} else if newLinkID, ok := linkToBookingLink[oldSlot.LinkID]; ok {
			// This slot belongs to a booking link - update to use new foreign key
			newSlot := Slot{
				BookingLinkID: newLinkID,
				Type:          oldSlot.Type,
				StartTime:     oldSlot.StartTime,
				EndTime:       oldSlot.EndTime,
				Manual:        oldSlot.Manual,
				CreatedAt:     oldSlot.CreatedAt,
			}
			if err := db.Create(&newSlot).Error; err != nil {
				log.Printf("Warning: Failed to create slot: %v", err)
				continue
			}
			fmt.Printf("  Migrated slot to booking link: %d -> %d\n", oldSlot.ID, newSlot.ID)
		}
	}

	// Update bookings to use new booking_link_id
	fmt.Println("Updating bookings...")
	for oldLinkID, newLinkID := range linkToBookingLink {
		result := db.Model(&Booking{}).Where("link_id = ?", oldLinkID).Update("booking_link_id", newLinkID)
		if result.Error != nil {
			log.Printf("Warning: Failed to update bookings for link %d: %v", oldLinkID, result.Error)
		} else {
			fmt.Printf("  Updated %d bookings for link %d -> booking_link %d\n", result.RowsAffected, oldLinkID, newLinkID)
		}
	}

	// Update votes to use new poll_id
	fmt.Println("Updating votes...")
	for oldLinkID, newPollID := range linkToPoll {
		result := db.Model(&Vote{}).Where("link_id = ?", oldLinkID).Update("poll_id", newPollID)
		if result.Error != nil {
			log.Printf("Warning: Failed to update votes for link %d: %v", oldLinkID, result.Error)
		} else {
			fmt.Printf("  Updated %d votes for link %d -> poll %d\n", result.RowsAffected, oldLinkID, newPollID)
		}
	}

	fmt.Println("\nMigration completed!")
	fmt.Println("\nSummary:")
	fmt.Printf("  Booking links created: %d\n", len(linkToBookingLink))
	fmt.Printf("  Polls created: %d\n", len(linkToPoll))
	fmt.Printf("  Poll options created: %d\n", len(slotToPollOption))
	fmt.Println("\nNote: You can now safely drop the old 'links' table after verifying the migration.")
	fmt.Println("      You may also need to drop the old 'link_id' columns from bookings and votes tables.")
}

func tableExists(db *gorm.DB, tableName string) bool {
	return db.Migrator().HasTable(tableName)
}
