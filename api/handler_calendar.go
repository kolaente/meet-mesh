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
