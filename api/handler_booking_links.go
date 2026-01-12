// api/handler_booking_links.go
package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

// ListBookingLinks returns all booking links for user
func (h *Handler) ListBookingLinks(ctx context.Context) ([]gen.BookingLink, error) {
	userID, _ := GetUserID(ctx)

	var links []BookingLink
	if err := h.db.Where("user_id = ?", userID).Find(&links).Error; err != nil {
		return nil, err
	}

	return mapBookingLinksToGen(links), nil
}

// CreateBookingLink creates a new booking link
func (h *Handler) CreateBookingLink(ctx context.Context, req *gen.CreateBookingLinkReq) (*gen.BookingLink, error) {
	userID, _ := GetUserID(ctx)

	link := BookingLink{
		UserID:            userID,
		Slug:              generateSlug(),
		Name:              req.Name,
		Description:       req.Description.Value,
		Status:            LinkStatusActive,
		AutoConfirm:       req.AutoConfirm.Value,
		RequireEmail:      req.RequireEmail.Value,
		AvailabilityRules: mapAvailabilityRulesFromGen(req.AvailabilityRules),
		CustomFields:      mapCustomFieldsFromGen(req.CustomFields),
		EventTemplate:     mapEventTemplateFromGen(req.EventTemplate),
	}

	if err := h.db.Create(&link).Error; err != nil {
		return nil, err
	}

	return mapBookingLinkToGen(&link), nil
}

// GetBookingLink returns booking link details
func (h *Handler) GetBookingLink(ctx context.Context, params gen.GetBookingLinkParams) (*gen.BookingLink, error) {
	userID, _ := GetUserID(ctx)

	var link BookingLink
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	return mapBookingLinkToGen(&link), nil
}

// UpdateBookingLink updates a booking link
func (h *Handler) UpdateBookingLink(ctx context.Context, req *gen.UpdateBookingLinkReq, params gen.UpdateBookingLinkParams) (*gen.BookingLink, error) {
	userID, _ := GetUserID(ctx)

	var link BookingLink
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

	return mapBookingLinkToGen(&link), nil
}

// DeleteBookingLink deletes a booking link
func (h *Handler) DeleteBookingLink(ctx context.Context, params gen.DeleteBookingLinkParams) error {
	userID, _ := GetUserID(ctx)

	return h.db.Where("id = ? AND user_id = ?", params.ID, userID).Delete(&BookingLink{}).Error
}

// GetBookingLinkBookings returns bookings for a booking link
func (h *Handler) GetBookingLinkBookings(ctx context.Context, params gen.GetBookingLinkBookingsParams) ([]gen.Booking, error) {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link BookingLink
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	var bookings []Booking
	if err := h.db.Preload("Slot").Where("booking_link_id = ?", params.ID).Order("created_at DESC").Find(&bookings).Error; err != nil {
		return nil, err
	}

	return mapBookingsToGen(bookings), nil
}

// Helper functions
func generateSlug() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:10]
}

func mapBookingLinksToGen(links []BookingLink) []gen.BookingLink {
	result := make([]gen.BookingLink, len(links))
	for i, link := range links {
		result[i] = *mapBookingLinkToGen(&link)
	}
	return result
}

func mapBookingLinkToGen(link *BookingLink) *gen.BookingLink {
	return &gen.BookingLink{
		ID:                int(link.ID),
		Slug:              link.Slug,
		Name:              link.Name,
		Description:       gen.NewOptString(link.Description),
		Status:            gen.LinkStatus(link.Status),
		AutoConfirm:       gen.NewOptBool(link.AutoConfirm),
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
