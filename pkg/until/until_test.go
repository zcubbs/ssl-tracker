package until

import (
	"testing"
	"time"
)

func TestTimeUntil(t *testing.T) {
	now := time.Date(2023, time.August, 11, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		time     time.Time
		expected string
	}{
		{now.Add(-time.Hour), "Time Passed"},
		{now.Add(time.Hour), "Less than a day remaining"},
		{now.Add(24 * time.Hour), "1 day(s) remaining"},
		{now.Add(7 * 24 * time.Hour), "1 week(s) remaining"},
		{now.Add(30 * 24 * time.Hour), "1 month(s) remaining"},
	}

	for _, tt := range tests {
		actual := TimeUntil(now, tt.time)
		if actual != tt.expected {
			t.Errorf("TimeUntil(%v, %v) = %v; want %v", now, tt.time, actual, tt.expected)
		}
	}
}

func TestHasDatePassed(t *testing.T) {
	now := time.Date(2023, time.August, 11, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		time     time.Time
		expected bool
	}{
		{now.Add(-time.Hour), false},     // Same day, so not "passed"
		{now.Add(time.Hour), false},      // Same day, so not "passed"
		{now.Add(-25 * time.Hour), true}, // Previous day, so "passed"
		{now.Add(25 * time.Hour), false}, // Next day, so not "passed"
	}

	for _, tt := range tests {
		actual := HasDatePassed(now, tt.time)
		if actual != tt.expected {
			t.Errorf("HasDatePassed(%v, %v) = %v; want %v", now, tt.time, actual, tt.expected)
		}
	}
}
