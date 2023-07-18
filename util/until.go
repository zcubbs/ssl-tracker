package util

import (
	"fmt"
	"time"
)

func TimeUntil(t time.Time) string {
	now := time.Now()

	// if the time is in the past, return "Time Passed"
	if t.Before(now) {
		return "Time Passed"
	}

	duration := t.Sub(now)

	days := int(duration.Hours() / 24)
	weeks := days / 7
	months := days / 30

	if months > 0 {
		return fmt.Sprintf("%d month(s) remaining", months)
	}

	if weeks > 0 {
		return fmt.Sprintf("%d week(s) remaining", weeks)
	}

	if days > 0 {
		return fmt.Sprintf("%d day(s) remaining", days)
	}

	// if less than a day is remaining
	return "Less than a day remaining"
}

func HasDatePassed(t time.Time) bool {
	now := time.Now()

	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if now.After(t) {
		return true
	}

	return false
}
