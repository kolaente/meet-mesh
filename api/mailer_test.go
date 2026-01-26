package api

import (
	"testing"
)

func TestEmailAttachment_Format(t *testing.T) {
	attachment := &EmailAttachment{
		Filename:    "event.ics",
		ContentType: "text/calendar; charset=utf-8; method=REQUEST",
		Data:        []byte("BEGIN:VCALENDAR\r\nEND:VCALENDAR"),
	}

	if attachment.Filename != "event.ics" {
		t.Error("Expected filename 'event.ics'")
	}
	if attachment.ContentType != "text/calendar; charset=utf-8; method=REQUEST" {
		t.Error("Expected correct content type for ICS")
	}
	if len(attachment.Data) == 0 {
		t.Error("Expected non-empty data")
	}
}
