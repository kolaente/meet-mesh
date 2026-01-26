// api/caldav.go
package api

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"gorm.io/gorm"
)

type CalDAVClient struct {
	db *gorm.DB
}

func NewCalDAVClient(db *gorm.DB) *CalDAVClient {
	return &CalDAVClient{db: db}
}

type TimePeriod struct {
	Start time.Time
	End   time.Time
}

func (c *CalDAVClient) createClient(conn *CalendarConnection) (*caldav.Client, error) {
	httpClient := webdav.HTTPClientWithBasicAuth(http.DefaultClient, conn.Username, conn.Password)
	return caldav.NewClient(httpClient, conn.ServerURL)
}

// GetBusyTimes fetches busy periods from all connected calendars
func (c *CalDAVClient) GetBusyTimes(ctx context.Context, userID uint, start, end time.Time) ([]TimePeriod, error) {
	var connections []CalendarConnection
	if err := c.db.Where("user_id = ?", userID).Find(&connections).Error; err != nil {
		return nil, err
	}

	var allBusy []TimePeriod
	for _, conn := range connections {
		client, err := c.createClient(&conn)
		if err != nil {
			continue
		}

		for _, calURL := range conn.CalendarURLs {
			query := &caldav.CalendarQuery{
				CompRequest: caldav.CalendarCompRequest{
					Name: "VCALENDAR",
					Comps: []caldav.CalendarCompRequest{{
						Name:     "VEVENT",
						AllProps: true,
					}},
				},
				CompFilter: caldav.CompFilter{
					Name: "VCALENDAR",
					Comps: []caldav.CompFilter{{
						Name:  "VEVENT",
						Start: start,
						End:   end,
					}},
				},
			}

			events, err := client.QueryCalendar(ctx, calURL, query)
			if err != nil {
				continue
			}

			for _, obj := range events {
				for _, event := range obj.Data.Events() {
					dtstart, _ := event.DateTimeStart(nil)
					dtend, _ := event.DateTimeEnd(nil)
					allBusy = append(allBusy, TimePeriod{
						Start: dtstart,
						End:   dtend,
					})
				}
			}
		}
	}

	return mergePeriods(allBusy), nil
}

// CalendarEvent represents a single calendar event for test results
type CalendarEvent struct {
	Title string
	Start time.Time
	End   time.Time
}

// DiscoveredCalendar represents a calendar found during discovery
type DiscoveredCalendar struct {
	URL                 string
	Name                string
	Description         string
	SupportedComponents []string
}

// DiscoverCalendars discovers available calendars from a CalDAV server
func (c *CalDAVClient) DiscoverCalendars(ctx context.Context, serverURL, username, password string) ([]DiscoveredCalendar, error) {
	httpClient := webdav.HTTPClientWithBasicAuth(http.DefaultClient, username, password)

	client, err := caldav.NewClient(httpClient, serverURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create CalDAV client: %w", err)
	}

	// Step 1: Find current user principal
	principal, err := client.FindCurrentUserPrincipal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find user principal: %w", err)
	}

	// Step 2: Find calendar home set
	homeSet, err := client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		return nil, fmt.Errorf("failed to find calendar home set: %w", err)
	}

	// Step 3: Find all calendars
	calendars, err := client.FindCalendars(ctx, homeSet)
	if err != nil {
		return nil, fmt.Errorf("failed to find calendars: %w", err)
	}

	result := make([]DiscoveredCalendar, len(calendars))
	for i, cal := range calendars {
		result[i] = DiscoveredCalendar{
			URL:                 cal.Path,
			Name:                cal.Name,
			Description:         cal.Description,
			SupportedComponents: cal.SupportedComponentSet,
		}
	}

	return result, nil
}

// TestCalendarConnection tests a specific calendar by fetching events for the next 7 days
func (c *CalDAVClient) TestCalendarConnection(ctx context.Context, connID uint, userID uint) ([]CalendarEvent, error) {
	var conn CalendarConnection
	if err := c.db.Where("id = ? AND user_id = ?", connID, userID).First(&conn).Error; err != nil {
		return nil, err
	}

	client, err := c.createClient(&conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create CalDAV client: %w", err)
	}

	start := time.Now()
	end := start.AddDate(0, 0, 7) // 7 days ahead

	var events []CalendarEvent
	for _, calURL := range conn.CalendarURLs {
		query := &caldav.CalendarQuery{
			CompRequest: caldav.CalendarCompRequest{
				Name: "VCALENDAR",
				Comps: []caldav.CalendarCompRequest{{
					Name:     "VEVENT",
					AllProps: true,
				}},
			},
			CompFilter: caldav.CompFilter{
				Name: "VCALENDAR",
				Comps: []caldav.CompFilter{{
					Name:  "VEVENT",
					Start: start,
					End:   end,
				}},
			},
		}

		objs, err := client.QueryCalendar(ctx, calURL, query)
		if err != nil {
			return nil, fmt.Errorf("failed to query calendar %s: %w", calURL, err)
		}

		for _, obj := range objs {
			for _, event := range obj.Data.Events() {
				dtstart, _ := event.DateTimeStart(nil)
				dtend, _ := event.DateTimeEnd(nil)
				summary := event.Props.Get(ical.PropSummary)
				title := "(No title)"
				if summary != nil {
					title = summary.Value
				}
				events = append(events, CalendarEvent{
					Title: title,
					Start: dtstart,
					End:   dtend,
				})
			}
		}
	}

	// Sort by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	return events, nil
}

// CreateBookingEvent creates a calendar event for a confirmed booking
func (c *CalDAVClient) CreateBookingEvent(ctx context.Context, userID uint, booking *Booking, slot *Slot, template *EventTemplate, meetingLink string) (string, error) {
	var conn CalendarConnection
	if err := c.db.Where("user_id = ? AND write_url != ''", userID).First(&conn).Error; err != nil {
		return "", err
	}

	client, err := c.createClient(&conn)
	if err != nil {
		return "", err
	}

	// Build iCal event
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropProductID, "-//Meet Mesh//EN")
	cal.Props.SetText(ical.PropVersion, "2.0")

	event := ical.NewEvent()
	uid := generateUID()
	event.Props.SetText(ical.PropUID, uid)
	event.Props.SetDateTime(ical.PropDateTimeStart, slot.StartTime)
	event.Props.SetDateTime(ical.PropDateTimeEnd, slot.EndTime)

	title := "Meeting"
	if template != nil && template.TitleTemplate != "" {
		title = expandTemplate(template.TitleTemplate, booking, meetingLink)
	}
	event.Props.SetText(ical.PropSummary, title)

	if template != nil && template.DescriptionTemplate != "" {
		event.Props.SetText(ical.PropDescription, expandTemplate(template.DescriptionTemplate, booking, meetingLink))
	}

	// Use meeting link as location if provided and no location set in template
	if meetingLink != "" {
		if template == nil || template.Location == "" {
			event.Props.SetText(ical.PropLocation, meetingLink)
		} else {
			event.Props.SetText(ical.PropLocation, template.Location)
		}
	} else if template != nil && template.Location != "" {
		event.Props.SetText(ical.PropLocation, template.Location)
	}

	cal.Children = append(cal.Children, event.Component)

	// Put the event
	path := conn.WriteURL + "/" + uid + ".ics"
	_, err = client.PutCalendarObject(ctx, path, cal)
	if err != nil {
		return "", err
	}

	return uid, nil
}

// Helper functions
func mergePeriods(periods []TimePeriod) []TimePeriod {
	if len(periods) == 0 {
		return periods
	}

	// Sort by start time
	sort.Slice(periods, func(i, j int) bool {
		return periods[i].Start.Before(periods[j].Start)
	})

	var merged []TimePeriod
	current := periods[0]
	for i := 1; i < len(periods); i++ {
		if periods[i].Start.Before(current.End) || periods[i].Start.Equal(current.End) {
			if periods[i].End.After(current.End) {
				current.End = periods[i].End
			}
		} else {
			merged = append(merged, current)
			current = periods[i]
		}
	}
	merged = append(merged, current)
	return merged
}

func generateUID() string {
	// Generate unique ID for calendar event
	return time.Now().Format("20060102T150405") + "@meet-mesh"
}

func expandTemplate(template string, booking *Booking, meetingLink string) string {
	// Simple template expansion
	result := template
	if booking != nil {
		// Replace {{guest_name}} with booking.GuestName
		// Replace {{guest_email}} with booking.GuestEmail
		// Basic string replacement - in production use text/template
		result = strings.ReplaceAll(result, "{{guest_name}}", booking.GuestName)
		result = strings.ReplaceAll(result, "{{guest_email}}", booking.GuestEmail)
	}
	// Replace {{meeting_link}} with the meeting link
	result = strings.ReplaceAll(result, "{{meeting_link}}", meetingLink)
	return result
}
