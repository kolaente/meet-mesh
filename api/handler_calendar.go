// api/handler_calendar.go
package api

import (
	"context"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

// ListCalendars returns calendar connections
func (h *Handler) ListCalendars(ctx context.Context) ([]gen.CalendarConnection, error) {
	userID, _ := GetUserID(ctx)

	var connections []CalendarConnection
	if err := h.db.Where("user_id = ?", userID).Find(&connections).Error; err != nil {
		return nil, err
	}

	result := make([]gen.CalendarConnection, len(connections))
	for i, conn := range connections {
		result[i] = gen.CalendarConnection{
			ID:           int(conn.ID),
			ServerURL:    conn.ServerURL,
			Username:     conn.Username,
			CalendarUrls: conn.CalendarURLs,
			WriteURL:     gen.NewOptString(conn.WriteURL),
		}
	}

	return result, nil
}

// AddCalendar adds a calendar connection
func (h *Handler) AddCalendar(ctx context.Context, req *gen.AddCalendarReq) (*gen.CalendarConnection, error) {
	userID, _ := GetUserID(ctx)

	conn := CalendarConnection{
		UserID:       userID,
		ServerURL:    req.ServerURL,
		Username:     req.Username,
		Password:     req.Password,
		CalendarURLs: req.CalendarUrls,
		WriteURL:     req.WriteURL.Value,
	}

	if err := h.db.Create(&conn).Error; err != nil {
		return nil, err
	}

	return &gen.CalendarConnection{
		ID:           int(conn.ID),
		ServerURL:    conn.ServerURL,
		Username:     conn.Username,
		CalendarUrls: conn.CalendarURLs,
		WriteURL:     gen.NewOptString(conn.WriteURL),
	}, nil
}

// RemoveCalendar removes a calendar connection
func (h *Handler) RemoveCalendar(ctx context.Context, params gen.RemoveCalendarParams) error {
	userID, _ := GetUserID(ctx)

	return h.db.Where("id = ? AND user_id = ?", params.ID, userID).Delete(&CalendarConnection{}).Error
}

// TestCalendar tests a calendar connection by fetching events
func (h *Handler) TestCalendar(ctx context.Context, params gen.TestCalendarParams) (gen.TestCalendarRes, error) {
	userID, _ := GetUserID(ctx)

	events, err := h.caldav.TestCalendarConnection(ctx, uint(params.ID), userID)
	if err != nil {
		return &gen.CalendarTestResult{
			Success: false,
			Error:   gen.NewOptString(err.Error()),
			Events:  []gen.CalendarTestResultEventsItem{},
		}, nil
	}

	resultEvents := make([]gen.CalendarTestResultEventsItem, len(events))
	for i, e := range events {
		resultEvents[i] = gen.CalendarTestResultEventsItem{
			Title: e.Title,
			Start: e.Start,
			End:   e.End,
		}
	}

	return &gen.CalendarTestResult{
		Success: true,
		Events:  resultEvents,
	}, nil
}

// DiscoverCalendars discovers available calendars from a CalDAV server
func (h *Handler) DiscoverCalendars(ctx context.Context, req *gen.DiscoverCalendarsReq) (*gen.CalendarDiscoveryResult, error) {
	calendars, err := h.caldav.DiscoverCalendars(ctx, req.ServerURL, req.Username, req.Password)
	if err != nil {
		return &gen.CalendarDiscoveryResult{
			Success:   false,
			Error:     gen.NewOptString(err.Error()),
			Calendars: []gen.DiscoveredCalendar{},
		}, nil
	}

	result := make([]gen.DiscoveredCalendar, len(calendars))
	for i, cal := range calendars {
		result[i] = gen.DiscoveredCalendar{
			URL:  cal.URL,
			Name: cal.Name,
		}
		if cal.Description != "" {
			result[i].Description = gen.NewOptString(cal.Description)
		}
		if len(cal.SupportedComponents) > 0 {
			result[i].SupportedComponents = cal.SupportedComponents
		}
	}

	return &gen.CalendarDiscoveryResult{
		Success:   true,
		Calendars: result,
	}, nil
}
