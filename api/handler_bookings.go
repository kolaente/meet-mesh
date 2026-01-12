// api/handler_bookings.go
package api

import (
	"context"
	"errors"

	gen "github.com/kolaente/meet-mesh/api/gen"
	"github.com/ogen-go/ogen/ogenerrors"
)

// ApproveBooking approves a booking
func (h *Handler) ApproveBooking(ctx context.Context, params gen.ApproveBookingParams) (*gen.Booking, error) {
	userID, _ := GetUserID(ctx)

	var booking Booking
	if err := h.db.Preload("BookingLink").Preload("Slot").First(&booking, params.ID).Error; err != nil {
		return nil, err
	}

	// Verify ownership
	if booking.BookingLink.UserID != userID {
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
		h.mailer.SendBookingApproved(&booking, &booking.BookingLink)
	}

	// Create calendar event
	if h.caldav != nil {
		h.caldav.CreateBookingEvent(ctx, booking.BookingLink.UserID, &booking, &booking.Slot, booking.BookingLink.EventTemplate)
	}

	return mapBookingToGen(&booking), nil
}

// DeclineBooking declines a booking
func (h *Handler) DeclineBooking(ctx context.Context, params gen.DeclineBookingParams) (*gen.Booking, error) {
	userID, _ := GetUserID(ctx)

	var booking Booking
	if err := h.db.Preload("BookingLink").Preload("Slot").First(&booking, params.ID).Error; err != nil {
		return nil, err
	}

	// Verify ownership
	if booking.BookingLink.UserID != userID {
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
		h.mailer.SendBookingDeclined(&booking, &booking.BookingLink)
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
