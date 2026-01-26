// api/ics.go
package api

import (
	"bytes"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

// GenerateICSData creates an ICS calendar file for a booking
func GenerateICSData(booking *Booking, slot *Slot, template *EventTemplate, organizerEmail string) (string, error) {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropProductID, "-//Meet Mesh//EN")
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropMethod, "REQUEST")

	event := ical.NewEvent()

	// Required: UID and DTSTAMP
	uid := generateUID()
	event.Props.SetText(ical.PropUID, uid)
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())

	// Time
	event.Props.SetDateTime(ical.PropDateTimeStart, slot.StartTime.UTC())
	event.Props.SetDateTime(ical.PropDateTimeEnd, slot.EndTime.UTC())

	// Title
	title := "Meeting"
	if template != nil && template.TitleTemplate != "" {
		title = expandTemplateICS(template.TitleTemplate, booking)
	}
	event.Props.SetText(ical.PropSummary, title)

	// Description (optional)
	if template != nil && template.DescriptionTemplate != "" {
		event.Props.SetText(ical.PropDescription, expandTemplateICS(template.DescriptionTemplate, booking))
	}

	// Location (optional)
	if template != nil && template.Location != "" {
		event.Props.SetText(ical.PropLocation, template.Location)
	}

	// Organizer
	organizerProp := ical.NewProp(ical.PropOrganizer)
	organizerProp.Value = "mailto:" + organizerEmail
	event.Props.Set(organizerProp)

	// Attendee (guest)
	attendeeProp := ical.NewProp(ical.PropAttendee)
	attendeeProp.Value = "mailto:" + booking.GuestEmail
	if booking.GuestName != "" {
		attendeeProp.Params.Set(ical.ParamCommonName, booking.GuestName)
	}
	event.Props.Set(attendeeProp)

	cal.Children = append(cal.Children, event.Component)

	// Encode to string
	var buf bytes.Buffer
	enc := ical.NewEncoder(&buf)
	if err := enc.Encode(cal); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// expandTemplateICS expands template variables (same as caldav.go expandTemplate)
func expandTemplateICS(template string, booking *Booking) string {
	result := template
	if booking != nil {
		result = strings.ReplaceAll(result, "{{guest_name}}", booking.GuestName)
		result = strings.ReplaceAll(result, "{{guest_email}}", booking.GuestEmail)
	}
	return result
}
