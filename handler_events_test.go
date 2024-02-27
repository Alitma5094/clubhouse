package main

import (
	"errors"
	"testing"
	"time"
)

func TestValidateEvent(t *testing.T) {

	//startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	//endTime, _ := time.Parse(time.RFC3339, "2021-01-01T02:00:00Z")
	//
	//got := ValidateEventInput("Hello, World", startTime, endTime)
	//
	//if got != nil {
	//	t.Errorf("got %s, wanted nil", got)
	//}
	//
	//startTime, _ = time.Parse(time.RFC3339, "")
	//endTime, _ = time.Parse(time.RFC3339, "")
	//
	//got = ValidateEventInput("Hello, World", startTime, endTime)
	//
	//if got != nil {
	//	t.Errorf("got %s, wanted nil", got)
	//}
	var tests = []struct {
		name      string
		title     string
		startTime string
		endTime   string
		want      error
	}{
		{"Valid Event", "First Event", "2021-01-01T00:00:00Z", "2021-01-01T02:00:00Z", nil},
		{"Title Too Long", "Second Event Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua", "2021-01-01T00:00:00Z", "2021-01-01T02:00:00Z", ErrEventTitleTooLong},
		{"Start Time Missing", "Third Event", "", "2021-01-01T02:00:00Z", ErrEventStartAtMissing},
		{"End Time Missing", "Fourth Event", "2021-01-01T00:00:00Z", "", ErrEventEndAtMissing},
		{"Start and End Missing", "Fifth Event", "", "", ErrEventStartAtMissing},
		{"End Time Before Start Time", "Sixth Event", "2021-01-01T02:00:00Z", "2021-01-001T00:00:00Z", ErrEventEndAtMissing},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime, _ := time.Parse(time.RFC3339, tt.startTime)
			endTime, _ := time.Parse(time.RFC3339, tt.endTime)
			ans := ValidateEventInput(tt.title, startTime, endTime)
			if !errors.Is(ans, tt.want) {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
