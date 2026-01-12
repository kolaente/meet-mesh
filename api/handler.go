// api/handler.go
package api

import (
	"gorm.io/gorm"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

// Handler implements the generated Handler interface
type Handler struct {
	gen.UnimplementedHandler
	db     *gorm.DB
	auth   *AuthService
	caldav *CalDAVClient
	mailer *Mailer
	config *Config
}

var _ gen.Handler = (*Handler)(nil)

func NewHandler(db *gorm.DB, auth *AuthService, caldav *CalDAVClient, mailer *Mailer, config *Config) *Handler {
	return &Handler{
		db:     db,
		auth:   auth,
		caldav: caldav,
		mailer: mailer,
		config: config,
	}
}
