// api/handler_public_booking.go
package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

// GetPublicBookingLink returns public booking link info
func (h *Handler) GetPublicBookingLink(ctx context.Context, params gen.GetPublicBookingLinkParams) (gen.GetPublicBookingLinkRes, error) {
	var link BookingLink
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &gen.GetPublicBookingLinkNotFound{Message: "Booking link not found"}, nil
		}
		return nil, err
	}

	return &gen.GetPublicBookingLinkOK{
		Name:         link.Name,
		Description:  gen.NewOptString(link.Description),
		CustomFields: mapCustomFieldsToGen(link.CustomFields),
		RequireEmail: gen.NewOptBool(link.RequireEmail),
	}, nil
}

// GetBookingAvailability returns real-time availability for a booking link
func (h *Handler) GetBookingAvailability(ctx context.Context, params gen.GetBookingAvailabilityParams) (*gen.GetBookingAvailabilityOK, error) {
	var link BookingLink
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		return nil, err
	}

	// Fetch busy times from CalDAV
	var busyTimes []TimePeriod
	if h.caldav != nil {
		var err error
		busyTimes, err = h.caldav.GetBusyTimes(ctx, link.UserID, params.Start, params.End)
		if err != nil {
			// Log error but continue - availability without calendar integration
			busyTimes = nil
		}
	}

	// Generate available slots based on availability rules
	slots := h.generateAvailableSlots(link, params.Start, params.End, busyTimes)

	return &gen.GetBookingAvailabilityOK{
		Slots: mapSlotsToGen(slots),
	}, nil
}

func (h *Handler) generateAvailableSlots(link BookingLink, start, end time.Time, busyTimes []TimePeriod) []Slot {
	var slots []Slot

	// If no availability rules, return empty
	if len(link.AvailabilityRules) == 0 {
		return slots
	}

	// Generate slots for each day in the range
	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		weekday := int(day.Weekday())

		for _, rule := range link.AvailabilityRules {
			// Check if this day is in the rule's days of week
			for _, dow := range rule.DaysOfWeek {
				if dow == weekday {
					// Parse start and end times
					startTime, _ := time.Parse("15:04", rule.StartTime)
					endTime, _ := time.Parse("15:04", rule.EndTime)

					slotStart := time.Date(day.Year(), day.Month(), day.Day(),
						startTime.Hour(), startTime.Minute(), 0, 0, day.Location())
					slotEnd := time.Date(day.Year(), day.Month(), day.Day(),
						endTime.Hour(), endTime.Minute(), 0, 0, day.Location())

					// Check if slot conflicts with busy times
					isBusy := false
					for _, busy := range busyTimes {
						if slotStart.Before(busy.End) && slotEnd.After(busy.Start) {
							isBusy = true
							break
						}
					}

					if !isBusy && slotStart.After(time.Now()) {
						slots = append(slots, Slot{
							Type:      SlotTypeTime,
							StartTime: slotStart,
							EndTime:   slotEnd,
						})
					}
				}
			}
		}
	}

	return slots
}

// CreateBooking creates a booking
func (h *Handler) CreateBooking(ctx context.Context, req *gen.CreateBookingReq, params gen.CreateBookingParams) (gen.CreateBookingRes, error) {
	var link BookingLink
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		return nil, err
	}

	// Verify slot exists and is available
	var slot Slot
	if err := h.db.Where("id = ? AND booking_link_id = ?", req.SlotID, link.ID).First(&slot).Error; err != nil {
		return &gen.CreateBookingConflict{Message: "Slot not available"}, nil
	}

	// Check CalDAV availability
	if h.caldav != nil {
		busyTimes, err := h.caldav.GetBusyTimes(ctx, link.UserID, slot.StartTime, slot.EndTime)
		if err == nil && len(busyTimes) > 0 {
			return &gen.CreateBookingConflict{Message: "Slot no longer available"}, nil
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

	// Convert custom fields
	var customFields map[string]string
	if req.CustomFields.Set {
		customFields = make(map[string]string)
		for k, v := range req.CustomFields.Value {
			customFields[k] = v
		}
	}

	booking := Booking{
		BookingLinkID: link.ID,
		SlotID:        slot.ID,
		GuestEmail:    req.GuestEmail,
		GuestName:     req.GuestName.Value,
		CustomFields:  customFields,
		Status:        status,
		ActionToken:   actionToken,
	}

	if err := h.db.Create(&booking).Error; err != nil {
		return nil, err
	}

	// Send notification email and create calendar event if auto-confirmed
	if link.AutoConfirm {
		if h.mailer != nil {
			h.mailer.SendBookingConfirmation(&booking, &link)
		}
		// Create calendar event
		if h.caldav != nil {
			h.caldav.CreateBookingEvent(ctx, link.UserID, &booking, &slot, link.EventTemplate)
		}
	} else {
		if h.mailer != nil {
			h.mailer.SendBookingPending(&booking, &link)
		}
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

// Helper for mapping slots
func mapSlotsToGen(slots []Slot) []gen.Slot {
	result := make([]gen.Slot, len(slots))
	for i, slot := range slots {
		result[i] = *mapSlotToGen(&slot)
	}
	return result
}

func mapSlotToGen(slot *Slot) *gen.Slot {
	return &gen.Slot{
		ID:        int(slot.ID),
		Type:      gen.SlotType(slot.Type),
		StartTime: slot.StartTime,
		EndTime:   slot.EndTime,
	}
}
