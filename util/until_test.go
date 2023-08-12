package util

import (
	"testing"
	"time"
)

func TestTimeUntil(t *testing.T) {
	tests := []struct {
		time     time.Time
		expected string
	}{
		{time.Now().Add(-time.Hour), "Time Passed"},
		{time.Now().Add(time.Hour), "Less than a day remaining"},
		{time.Now().Add(24 * time.Hour), "1 day(s) remaining"},
		{time.Now().Add(7 * 24 * time.Hour), "1 week(s) remaining"},
		{time.Now().Add(30 * 24 * time.Hour), "1 month(s) remaining"},
	}

	for _, tt := range tests {
		actual := TimeUntil(tt.time)
		if actual != tt.expected {
			t.Errorf("TimeUntil(%v) = %v; want %v", tt.time, actual, tt.expected)
		}
	}
}

func TestHasDatePassed(t *testing.T) {
	tests := []struct {
		time     time.Time
		expected bool
	}{
		{time.Now().Add(-time.Hour * 2), true},
		{time.Now().Add(2 * time.Hour), false},
	}

	for _, tt := range tests {
		actual := HasDatePassed(tt.time)
		if actual != tt.expected {
			t.Errorf("HasDatePassed(%v) = %v; want %v", tt.time, actual, tt.expected)
		}
	}
}
