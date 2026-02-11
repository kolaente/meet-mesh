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
			return &gen.Error{Message: "Booking link not found"}, nil
		}
		return nil, err
	}

	// Determine available durations
	var durations []int
	if len(link.SlotDurationsMinutes) > 0 {
		durations = link.SlotDurationsMinutes
	} else {
		durations = []int{link.SlotDurationMinutes}
	}

	return &gen.GetPublicBookingLinkOK{
		Name:                 link.Name,
		Description:          gen.NewOptString(link.Description),
		CustomFields:         mapCustomFieldsToGen(link.CustomFields),
		RequireEmail:         gen.NewOptBool(link.RequireEmail),
		SlotDurationsMinutes: durations,
	}, nil
}

// GetBookingAvailability returns real-time availability for a booking link
func (h *Handler) GetBookingAvailability(ctx context.Context, params gen.GetBookingAvailabilityParams) (*gen.GetBookingAvailabilityOK, error) {
	var link BookingLink
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		return nil, err
	}

	// Determine which duration to use
	duration := link.SlotDurationMinutes // Default
	if params.Duration.Set {
		requestedDuration := params.Duration.Value
		// Validate requested duration is allowed
		allowed := false
		if len(link.SlotDurationsMinutes) > 0 {
			for _, d := range link.SlotDurationsMinutes {
				if d == requestedDuration {
					allowed = true
					break
				}
			}
		} else if requestedDuration == link.SlotDurationMinutes {
			allowed = true
		}
		if !allowed {
			return &gen.GetBookingAvailabilityOK{Slots: []gen.Slot{}}, nil
		}
		duration = requestedDuration
	} else if len(link.SlotDurationsMinutes) > 0 {
		// Use first duration as default
		duration = link.SlotDurationsMinutes[0]
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
	slots := h.generateAvailableSlotsWithDuration(link, params.Start, params.End, busyTimes, duration)

	return &gen.GetBookingAvailabilityOK{
		Slots: mapSlotsToGen(slots),
	}, nil
}



func (h *Handler) generateAvailableSlotsWithDuration(link BookingLink, start, end time.Time, busyTimes []TimePeriod, durationMinutes int) []Slot {
	var slots []Slot

	// If no availability rules, return empty
	if len(link.AvailabilityRules) == 0 {
		return slots
	}

	slotDuration := time.Duration(durationMinutes) * time.Minute
	bufferDuration := time.Duration(link.BufferMinutes) * time.Minute

	// Generate slots for each day in the range
	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		weekday := int(day.Weekday())

		for _, rule := range link.AvailabilityRules {
			if !containsDay(rule.DaysOfWeek, weekday) {
				continue
			}

			// Parse availability window times
			ruleStart, _ := time.Parse("15:04", rule.StartTime)
			ruleEnd, _ := time.Parse("15:04", rule.EndTime)

			windowStart := time.Date(day.Year(), day.Month(), day.Day(),
				ruleStart.Hour(), ruleStart.Minute(), 0, 0, day.Location())
			windowEnd := time.Date(day.Year(), day.Month(), day.Day(),
				ruleEnd.Hour(), ruleEnd.Minute(), 0, 0, day.Location())

			// Generate individual slots within this window
			for slotStart := windowStart; slotStart.Add(slotDuration).Before(windowEnd) || slotStart.Add(slotDuration).Equal(windowEnd); slotStart = slotStart.Add(slotDuration + bufferDuration) {
				slotEnd := slotStart.Add(slotDuration)

				// Skip past slots
				if slotStart.Before(time.Now()) {
					continue
				}

				// Check if slot conflicts with busy times
				if isSlotBusy(slotStart, slotEnd, busyTimes) {
					continue
				}

				slots = append(slots, Slot{
					Type:      SlotTypeTime,
					StartTime: slotStart,
					EndTime:   slotEnd,
				})
			}
		}
	}

	return slots
}

func containsDay(days []int, day int) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}

func isSlotBusy(start, end time.Time, busyTimes []TimePeriod) bool {
	for _, busy := range busyTimes {
		if start.Before(busy.End) && end.After(busy.Start) {
			return true
		}
	}
	return false
}

// isSlotWithinAvailability checks if a slot falls within the configured availability rules
func (h *Handler) isSlotWithinAvailability(start, end time.Time, rules []AvailabilityRule) bool {
	weekday := int(start.Weekday())
	slotStartTime := start.Format("15:04")
	slotEndTime := end.Format("15:04")

	for _, rule := range rules {
		if !containsDay(rule.DaysOfWeek, weekday) {
			continue
		}
		if slotStartTime >= rule.StartTime && slotEndTime <= rule.EndTime {
			return true
		}
	}
	return false
}

// CreateBooking creates a booking
func (h *Handler) CreateBooking(ctx context.Context, req *gen.CreateBookingReq, params gen.CreateBookingParams) (gen.CreateBookingRes, error) {
	var link BookingLink
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		return nil, err
	}

	// Validate slot duration matches one of the booking link's configurations
	requestedDuration := req.EndTime.Sub(req.StartTime)
	requestedMinutes := int(requestedDuration.Minutes())

	validDuration := false
	if len(link.SlotDurationsMinutes) > 0 {
		for _, d := range link.SlotDurationsMinutes {
			if requestedMinutes == d {
				validDuration = true
				break
			}
		}
	} else if requestedMinutes == link.SlotDurationMinutes {
		validDuration = true
	}

	if !validDuration {
		return &gen.Error{Message: "Invalid slot duration"}, nil
	}

	// Check the slot is not in the past
	if req.StartTime.Before(time.Now()) {
		return &gen.Error{Message: "Cannot book slots in the past"}, nil
	}

	// Validate that the slot falls within availability rules
	if !h.isSlotWithinAvailability(req.StartTime, req.EndTime, link.AvailabilityRules) {
		return &gen.Error{Message: "Slot not within available hours"}, nil
	}

	// Create slot from request times
	slot := Slot{
		BookingLinkID: link.ID,
		Type:          SlotTypeTime,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	}

	// Check CalDAV availability
	if h.caldav != nil {
		busyTimes, err := h.caldav.GetBusyTimes(ctx, link.UserID, req.StartTime, req.EndTime)
		if err == nil && len(busyTimes) > 0 {
			return &gen.Error{Message: "Slot no longer available"}, nil
		}
	}

	// Save the slot
	if err := h.db.Create(&slot).Error; err != nil {
		return nil, err
	}

	// Generate action token
	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)
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

	// Get organizer for emails
	var organizer User
	h.db.First(&organizer, link.UserID)

	// Send notification email and create calendar event if auto-confirmed
	if link.AutoConfirm {
		if h.mailer != nil {
			_ = h.mailer.SendBookingConfirmationWithICS(&booking, &link, organizer.Email)
		}
		// Create calendar event
		if h.caldav != nil {
			uid, err := h.caldav.CreateBookingEvent(ctx, link.UserID, &booking, &slot, link.EventTemplate, link.MeetingLink)
			if err == nil && uid != "" {
				booking.CalendarUID = uid
				h.db.Save(&booking)
			}
		}
	} else {
		if h.mailer != nil {
			_ = h.mailer.SendBookingPending(&booking, &link, &organizer)
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
