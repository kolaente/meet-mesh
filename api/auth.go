// api/auth.go
package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type AuthService struct {
	config       *OIDCConfig
	provider     *oidc.Provider
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewAuthService(ctx context.Context, cfg *OIDCConfig) (*AuthService, error) {
	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return &AuthService{
		config:       cfg,
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
	}, nil
}

func (a *AuthService) GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (a *AuthService) AuthCodeURL(state string) string {
	return a.oauth2Config.AuthCodeURL(state)
}

type UserClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (a *AuthService) Exchange(ctx context.Context, code string) (*UserClaims, error) {
	token, err := a.oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Printf("OAuth token exchange failed: %v", err)
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Printf("No id_token in OAuth response")
		return nil, errors.New("no id_token in response")
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Printf("ID token verification failed: %v", err)
		return nil, err
	}

	var claims UserClaims
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("Failed to extract claims: %v", err)
		return nil, err
	}

	return &claims, nil
}

// Session management
type Session struct {
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (a *AuthService) CreateSessionCookie(session *Session) (*http.Cookie, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	// In production, encrypt this with a secret key
	encoded := base64.URLEncoding.EncodeToString(data)

	return &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: Make configurable for production
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	}, nil
}

func (a *AuthService) ParseSessionCookie(cookie *http.Cookie) (*Session, error) {
	data, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, err
	}

	return &session, nil
}
