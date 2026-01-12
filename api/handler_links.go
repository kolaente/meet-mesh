// api/handler_links.go
package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

func generateSlug() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:10]
}

// ListLinks returns all links for user
func (h *Handler) ListLinks(ctx context.Context) ([]gen.Link, error) {
	userID, _ := GetUserID(ctx)

	var links []Link
	if err := h.db.Where("user_id = ?", userID).Find(&links).Error; err != nil {
		return nil, err
	}

	return mapLinksToGen(links), nil
}

// CreateLink creates a new link
func (h *Handler) CreateLink(ctx context.Context, req *gen.CreateLinkReq) (*gen.Link, error) {
	userID, _ := GetUserID(ctx)

	link := Link{
		UserID:            userID,
		Slug:              generateSlug(),
		Type:              LinkType(req.Type),
		Name:              req.Name,
		Description:       req.Description.Value,
		Status:            LinkStatusActive,
		AutoConfirm:       req.AutoConfirm.Value,
		ShowResults:       req.ShowResults.Value,
		RequireEmail:      req.RequireEmail.Value,
		AvailabilityRules: mapAvailabilityRulesFromGen(req.AvailabilityRules),
		CustomFields:      mapCustomFieldsFromGen(req.CustomFields),
		EventTemplate:     mapEventTemplateFromGen(req.EventTemplate),
	}

	if err := h.db.Create(&link).Error; err != nil {
		return nil, err
	}

	return mapLinkToGen(&link), nil
}

// GetLink returns link details
func (h *Handler) GetLink(ctx context.Context, params gen.GetLinkParams) (*gen.Link, error) {
	userID, _ := GetUserID(ctx)

	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	return mapLinkToGen(&link), nil
}

// UpdateLink updates a link
func (h *Handler) UpdateLink(ctx context.Context, req *gen.UpdateLinkReq, params gen.UpdateLinkParams) (*gen.Link, error) {
	userID, _ := GetUserID(ctx)

	var link Link
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
	if req.ShowResults.Set {
		link.ShowResults = req.ShowResults.Value
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

	return mapLinkToGen(&link), nil
}

// DeleteLink deletes a link
func (h *Handler) DeleteLink(ctx context.Context, params gen.DeleteLinkParams) error {
	userID, _ := GetUserID(ctx)

	return h.db.Where("id = ? AND user_id = ?", params.ID, userID).Delete(&Link{}).Error
}

// Helper mapping functions
func mapLinksToGen(links []Link) []gen.Link {
	result := make([]gen.Link, len(links))
	for i, link := range links {
		result[i] = *mapLinkToGen(&link)
	}
	return result
}

func mapLinkToGen(link *Link) *gen.Link {
	return &gen.Link{
		ID:                int(link.ID),
		Slug:              link.Slug,
		Type:              gen.LinkType(link.Type),
		Name:              link.Name,
		Description:       gen.NewOptString(link.Description),
		Status:            gen.LinkStatus(link.Status),
		AutoConfirm:       gen.NewOptBool(link.AutoConfirm),
		ShowResults:       gen.NewOptBool(link.ShowResults),
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
