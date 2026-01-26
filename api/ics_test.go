package api

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateICSData(t *testing.T) {
	booking := &Booking{
		GuestEmail: "guest@example.com",
		GuestName:  "John Doe",
	}
	slot := &Slot{
		StartTime: time.Date(2026, 2, 15, 14, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 2, 15, 15, 0, 0, 0, time.UTC),
	}
	template := &EventTemplate{
		TitleTemplate:       "Meeting with {{guest_name}}",
		DescriptionTemplate: "Booking with {{guest_email}}",
		Location:            "Conference Room A",
	}
	organizerEmail := "organizer@example.com"

	icsData, err := GenerateICSData(booking, slot, template, organizerEmail)
	if err != nil {
		t.Fatalf("GenerateICSData failed: %v", err)
	}

	// Verify required components
	checks := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Meet Mesh//EN",
		"BEGIN:VEVENT",
		"SUMMARY:Meeting with John Doe",
		"DESCRIPTION:Booking with guest@example.com",
		"LOCATION:Conference Room A",
		"DTSTART:20260215T140000Z",
		"DTEND:20260215T150000Z",
		"ORGANIZER:mailto:organizer@example.com",
		"guest@example.com",
		"END:VEVENT",
		"END:VCALENDAR",
	}

	for _, check := range checks {
		if !strings.Contains(icsData, check) {
			t.Errorf("ICS data missing %q", check)
		}
	}
}

func TestGenerateICSData_MinimalTemplate(t *testing.T) {
	booking := &Booking{
		GuestEmail: "guest@example.com",
	}
	slot := &Slot{
		StartTime: time.Date(2026, 2, 15, 14, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 2, 15, 15, 0, 0, 0, time.UTC),
	}
	organizerEmail := "organizer@example.com"

	// nil template should use defaults
	icsData, err := GenerateICSData(booking, slot, nil, organizerEmail)
	if err != nil {
		t.Fatalf("GenerateICSData failed: %v", err)
	}

	if !strings.Contains(icsData, "SUMMARY:Meeting") {
		t.Errorf("Expected default summary 'Meeting', got: %s", icsData)
	}
}
