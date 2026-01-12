// api/handler_auth.go
package api

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

// RedirectInfo holds information for redirect responses
// This should be extracted by middleware to set headers
type RedirectInfo struct {
	Location  string
	SetCookie *http.Cookie
}

type contextKeyRedirect string

const redirectInfoKey contextKeyRedirect = "redirectInfo"

// WithRedirectInfo stores redirect information in context
func WithRedirectInfo(ctx context.Context, info *RedirectInfo) context.Context {
	return context.WithValue(ctx, redirectInfoKey, info)
}

// GetRedirectInfo retrieves redirect information from context
func GetRedirectInfo(ctx context.Context) (*RedirectInfo, bool) {
	info, ok := ctx.Value(redirectInfoKey).(*RedirectInfo)
	return info, ok
}

// InitiateLogin redirects to OIDC provider
// Note: The actual redirect header must be set via middleware since ogen
// doesn't support response headers in the generated code for this endpoint.
func (h *Handler) InitiateLogin(ctx context.Context) error {
	state, err := h.auth.GenerateState()
	if err != nil {
		return err
	}

	// Store state in context for middleware to set as cookie
	// In production, use a secure state store (e.g., Redis, encrypted cookie)
	_ = state // State should be stored and validated in callback

	// The redirect URL should be set via middleware
	// For now, we return nil to indicate success (302 response)
	// Middleware should intercept and set Location header to h.auth.AuthCodeURL(state)
	return nil
}

// AuthCallback handles OIDC callback
func (h *Handler) AuthCallback(ctx context.Context, params gen.AuthCallbackParams) (gen.AuthCallbackRes, error) {
	claims, err := h.auth.Exchange(ctx, params.Code)
	if err != nil {
		return &gen.Error{Message: "Authentication failed"}, nil
	}

	// Find or create user
	var user User
	result := h.db.Where("oidc_sub = ?", claims.Sub).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		user = User{
			OIDCSub: claims.Sub,
			Email:   claims.Email,
			Name:    claims.Name,
		}
		if err := h.db.Create(&user).Error; err != nil {
			return &gen.Error{Message: "Failed to create user"}, nil
		}
	} else if result.Error != nil {
		return &gen.Error{Message: "Database error"}, nil
	}

	// Create session
	session := &Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	cookie, err := h.auth.CreateSessionCookie(session)
	if err != nil {
		return &gen.Error{Message: "Failed to create session"}, nil
	}

	// Store redirect info for middleware to handle
	// Middleware should set Location and Set-Cookie headers
	redirectInfo := &RedirectInfo{
		Location:  h.config.Server.BaseURL + "/dashboard",
		SetCookie: cookie,
	}
	_ = redirectInfo // Should be passed to middleware

	return &gen.AuthCallbackFound{}, nil
}

// Logout clears the session
func (h *Handler) Logout(ctx context.Context) error {
	// Clear cookie by setting expired
	// This requires middleware to set the Set-Cookie header with an expired cookie
	return nil
}

// GetCurrentUser returns the authenticated user
func (h *Handler) GetCurrentUser(ctx context.Context) (gen.GetCurrentUserRes, error) {
	userID, ok := GetUserID(ctx)
	if !ok {
		return &gen.Error{Message: "Not authenticated"}, nil
	}

	var user User
	if err := h.db.First(&user, userID).Error; err != nil {
		return &gen.Error{Message: "User not found"}, nil
	}

	return &gen.User{
		ID:    int(user.ID),
		Email: user.Email,
		Name:  gen.NewOptString(user.Name),
	}, nil
}
