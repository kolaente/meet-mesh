// api/mailer.go
package api

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	config    *SMTPConfig
	baseURL   string
	dialer    *gomail.Dialer
	templates *template.Template
}

func NewMailer(cfg *SMTPConfig, baseURL string) (*Mailer, error) {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	// Parse templates
	tmpl, err := template.New("emails").Parse(emailTemplates)
	if err != nil {
		return nil, err
	}

	return &Mailer{
		config:    cfg,
		baseURL:   baseURL,
		dialer:    dialer,
		templates: tmpl,
	}, nil
}

func (m *Mailer) send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.config.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)
}

func (m *Mailer) renderTemplate(name string, data any) string {
	var buf bytes.Buffer
	m.templates.ExecuteTemplate(&buf, name, data)
	return buf.String()
}

// SendBookingConfirmation sends confirmation to guest
func (m *Mailer) SendBookingConfirmation(booking *Booking, link *Link) error {
	body := m.renderTemplate("booking_confirmed_guest", map[string]any{
		"LinkName":  link.Name,
		"GuestName": booking.GuestName,
		"Time":      booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
	})
	return m.send(booking.GuestEmail, "Booking Confirmed: "+link.Name, body)
}

// SendBookingPending sends pending notification to organizer
func (m *Mailer) SendBookingPending(booking *Booking, link *Link) error {
	approveURL := fmt.Sprintf("%s/api/v1/actions/approve?token=%s", m.baseURL, booking.ActionToken)
	declineURL := fmt.Sprintf("%s/api/v1/actions/decline?token=%s", m.baseURL, booking.ActionToken)

	body := m.renderTemplate("booking_pending", map[string]any{
		"LinkName":   link.Name,
		"GuestEmail": booking.GuestEmail,
		"GuestName":  booking.GuestName,
		"Time":       booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
		"ApproveURL": approveURL,
		"DeclineURL": declineURL,
	})

	// For now, use a placeholder - in production, fetch organizer email from user
	return m.send("organizer@example.com", "New Booking Request: "+link.Name, body)
}

// SendBookingApproved sends approval notification to guest
func (m *Mailer) SendBookingApproved(booking *Booking, link *Link) error {
	body := m.renderTemplate("booking_approved", map[string]any{
		"LinkName":  link.Name,
		"GuestName": booking.GuestName,
		"Time":      booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
	})
	return m.send(booking.GuestEmail, "Booking Approved: "+link.Name, body)
}

// SendBookingDeclined sends decline notification to guest
func (m *Mailer) SendBookingDeclined(booking *Booking, link *Link) error {
	body := m.renderTemplate("booking_declined", map[string]any{
		"LinkName":  link.Name,
		"GuestName": booking.GuestName,
		"Time":      booking.Slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
	})
	return m.send(booking.GuestEmail, "Booking Declined: "+link.Name, body)
}

// SendPollWinner sends winner notification to all voters
func (m *Mailer) SendPollWinner(link *Link, slot *Slot, votes []Vote) error {
	body := m.renderTemplate("poll_winner", map[string]any{
		"LinkName": link.Name,
		"Time":     slot.StartTime.Format("Monday, January 2 at 3:04 PM"),
	})

	for _, vote := range votes {
		if vote.GuestEmail != "" {
			m.send(vote.GuestEmail, "Date Selected: "+link.Name, body)
		}
	}
	return nil
}

const emailTemplates = `
{{define "booking_confirmed_guest"}}
<html>
<body>
<h1>Booking Confirmed!</h1>
<p>Hi {{.GuestName}},</p>
<p>Your booking for <strong>{{.LinkName}}</strong> has been confirmed.</p>
<p><strong>When:</strong> {{.Time}}</p>
</body>
</html>
{{end}}

{{define "booking_pending"}}
<html>
<body>
<h1>New Booking Request</h1>
<p>You have a new booking request for <strong>{{.LinkName}}</strong>.</p>
<p><strong>Guest:</strong> {{.GuestName}} ({{.GuestEmail}})</p>
<p><strong>Requested time:</strong> {{.Time}}</p>
<p>
<a href="{{.ApproveURL}}" style="background:#22c55e;color:white;padding:10px 20px;text-decoration:none;border-radius:5px;">Approve</a>
<a href="{{.DeclineURL}}" style="background:#ef4444;color:white;padding:10px 20px;text-decoration:none;border-radius:5px;margin-left:10px;">Decline</a>
</p>
</body>
</html>
{{end}}

{{define "booking_approved"}}
<html>
<body>
<h1>Booking Approved!</h1>
<p>Hi {{.GuestName}},</p>
<p>Great news! Your booking for <strong>{{.LinkName}}</strong> has been approved.</p>
<p><strong>When:</strong> {{.Time}}</p>
</body>
</html>
{{end}}

{{define "booking_declined"}}
<html>
<body>
<h1>Booking Update</h1>
<p>Hi {{.GuestName}},</p>
<p>Unfortunately, your booking request for <strong>{{.LinkName}}</strong> could not be accommodated.</p>
<p><strong>Requested time:</strong> {{.Time}}</p>
</body>
</html>
{{end}}

{{define "poll_winner"}}
<html>
<body>
<h1>Date Selected!</h1>
<p>The organizer has selected a date for <strong>{{.LinkName}}</strong>.</p>
<p><strong>Selected time:</strong> {{.Time}}</p>
</body>
</html>
{{end}}
`
