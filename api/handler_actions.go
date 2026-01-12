// api/handler_actions.go
package api

import (
	"context"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

// ApproveViaEmail approves booking via email link
func (h *Handler) ApproveViaEmail(ctx context.Context) (gen.ApproveViaEmailRes, error) {
	bookingID, ok := GetBookingID(ctx)
	if !ok {
		return &gen.Error{Message: "Invalid token"}, nil
	}

	var booking Booking
	if err := h.db.Preload("BookingLink").Preload("Slot").First(&booking, bookingID).Error; err != nil {
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
	if h.mailer != nil {
		h.mailer.SendBookingApproved(&booking, &booking.BookingLink)
	}

	// Create calendar event
	if h.caldav != nil {
		h.caldav.CreateBookingEvent(ctx, booking.BookingLink.UserID, &booking, &booking.Slot, booking.BookingLink.EventTemplate)
	}

	return &gen.ApproveViaEmailOK{
		Message: gen.NewOptString("Booking approved successfully"),
	}, nil
}

// DeclineViaEmail declines booking via email link
func (h *Handler) DeclineViaEmail(ctx context.Context) (gen.DeclineViaEmailRes, error) {
	bookingID, ok := GetBookingID(ctx)
	if !ok {
		return &gen.Error{Message: "Invalid token"}, nil
	}

	var booking Booking
	if err := h.db.Preload("BookingLink").Preload("Slot").First(&booking, bookingID).Error; err != nil {
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
	if h.mailer != nil {
		h.mailer.SendBookingDeclined(&booking, &booking.BookingLink)
	}

	return &gen.DeclineViaEmailOK{
		Message: gen.NewOptString("Booking declined"),
	}, nil
}
