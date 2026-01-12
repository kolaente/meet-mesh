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

// Placeholder type for Mailer that will be implemented later
type Mailer struct{}

// SendPollWinner sends notification about the poll winner to all voters
func (m *Mailer) SendPollWinner(link *Link, slot *Slot, votes []Vote) {
	// TODO: Implement email notification
}

// SendBookingConfirmation sends confirmation email to guest
func (m *Mailer) SendBookingConfirmation(booking *Booking, link *Link) error {
	// TODO: Implement email sending
	return nil
}

// SendBookingPending sends pending approval notification to guest
func (m *Mailer) SendBookingPending(booking *Booking, link *Link) error {
	// TODO: Implement email sending
	return nil
}

// SendBookingApproved sends approval notification to guest
func (m *Mailer) SendBookingApproved(booking *Booking, link *Link) error {
	// TODO: Implement email sending
	return nil
}

// SendBookingDeclined sends decline notification to guest
func (m *Mailer) SendBookingDeclined(booking *Booking, link *Link) error {
	// TODO: Implement email sending
	return nil
}
