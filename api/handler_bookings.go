// api/handler_bookings.go
package api

import (
	"context"
	"errors"

	gen "github.com/kolaente/meet-mesh/api/gen"
	"github.com/ogen-go/ogen/ogenerrors"
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
		return nil, &ogenerrors.DecodeBodyError{
			ContentType: "application/json",
			Body:        nil,
			Err:         errors.New("not authorized"),
		}
	}

	booking.Status = BookingStatusConfirmed
	if err := h.db.Save(&booking).Error; err != nil {
		return nil, err
	}

	// Send confirmation email
	if h.mailer != nil {
		h.mailer.SendBookingApproved(&booking, &booking.Link)
	}

	// Create calendar event
	if h.caldav != nil {
		h.caldav.CreateEvent(ctx, booking.Link.UserID, &booking, &booking.Slot, booking.Link.EventTemplate)
	}

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
		return nil, &ogenerrors.DecodeBodyError{
			ContentType: "application/json",
			Body:        nil,
			Err:         errors.New("not authorized"),
		}
	}

	booking.Status = BookingStatusDeclined
	if err := h.db.Save(&booking).Error; err != nil {
		return nil, err
	}

	// Send decline email
	if h.mailer != nil {
		h.mailer.SendBookingDeclined(&booking, &booking.Link)
	}

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
		ID:           int(b.ID),
		Slot:         *mapSlotToGen(&b.Slot),
		GuestEmail:   b.GuestEmail,
		GuestName:    gen.NewOptString(b.GuestName),
		Status:       gen.BookingStatus(b.Status),
		CustomFields: mapBookingCustomFieldsToGen(b.CustomFields),
		CreatedAt:    gen.NewOptDateTime(b.CreatedAt),
	}
}

func mapBookingCustomFieldsToGen(fields map[string]string) gen.OptBookingCustomFields {
	if fields == nil {
		return gen.OptBookingCustomFields{}
	}
	return gen.NewOptBookingCustomFields(gen.BookingCustomFields(fields))
}
