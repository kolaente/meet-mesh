// api/security.go
package api

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

type contextKey string

const (
	userIDKey    contextKey = "userID"
	bookingIDKey contextKey = "bookingID"
)

func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDKey).(uint)
	return userID, ok
}

func WithBookingID(ctx context.Context, bookingID uint) context.Context {
	return context.WithValue(ctx, bookingIDKey, bookingID)
}

func GetBookingID(ctx context.Context) (uint, bool) {
	bookingID, ok := ctx.Value(bookingIDKey).(uint)
	return bookingID, ok
}

type SecurityHandler struct {
	db   *gorm.DB
	auth *AuthService
}

func NewSecurityHandler(db *gorm.DB, auth *AuthService) *SecurityHandler {
	return &SecurityHandler{db: db, auth: auth}
}

var _ gen.SecurityHandler = (*SecurityHandler)(nil)

func (s *SecurityHandler) HandleCookieAuth(ctx context.Context, operationName gen.OperationName, t gen.CookieAuth) (context.Context, error) {
	// t.APIKey contains the session cookie value
	cookie := &http.Cookie{Value: t.APIKey}

	session, err := s.auth.ParseSessionCookie(cookie)
	if err != nil {
		return ctx, errors.New("invalid session")
	}

	return WithUserID(ctx, session.UserID), nil
}

func (s *SecurityHandler) HandleActionToken(ctx context.Context, operationName gen.OperationName, t gen.ActionToken) (context.Context, error) {
	var booking Booking
	if err := s.db.Where("action_token = ?", t.APIKey).First(&booking).Error; err != nil {
		return ctx, errors.New("invalid action token")
	}

	return WithBookingID(ctx, booking.ID), nil
}
