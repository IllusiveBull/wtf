package daysmatter

import (
	"fmt"
	"time"
)

// DaysEvent represents a memorable event with date calculations
type DaysEvent struct {
	Event
	targetDate time.Time
}

// NewDaysEvent creates a new DaysEvent from an Event
func NewDaysEvent(event Event) (*DaysEvent, error) {
	targetDate, err := time.Parse("2006-01-02", event.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format for event '%s': %v", event.Name, err)
	}

	return &DaysEvent{
		Event:      event,
		targetDate: targetDate,
	}, nil
}

// GetDisplayInfo returns the display text and color for this event
func (de *DaysEvent) GetDisplayInfo() (string, string) {
	now := time.Now()
	currentYear := now.Year()

	// Calculate the target date for this year if it's a yearly event
	var targetDate time.Time
	if de.Yearly {
		// For yearly events, use this year's date
		targetDate = time.Date(currentYear, de.targetDate.Month(), de.targetDate.Day(), 0, 0, 0, 0, time.Local)

		// If this year's date has passed, use next year's date
		if targetDate.Before(now) {
			targetDate = time.Date(currentYear+1, de.targetDate.Month(), de.targetDate.Day(), 0, 0, 0, 0, time.Local)
		}
	} else {
		// For one-time events, use the original date
		targetDate = de.targetDate
	}

	// Calculate days difference
	diff := targetDate.Sub(now)
	days := int(diff.Hours() / 24)

	var displayText string
	var color string

	if days > 0 {
		// Future event
		if days == 1 {
			displayText = fmt.Sprintf("%s - 还有 1 天", de.Name)
		} else {
			displayText = fmt.Sprintf("%s - 还有 %d 天", de.Name, days)
		}

		// Color based on urgency
		if days <= 7 {
			color = "red"
		} else if days <= 30 {
			color = "yellow"
		} else {
			color = "green"
		}
	} else if days == 0 {
		// Today
		displayText = fmt.Sprintf("%s - 今天！", de.Name)
		color = "red"
	} else {
		// Past event
		absDays := -days
		if de.Yearly {
			// For yearly events that have passed this year, show days since this year's occurrence
			thisYearDate := time.Date(currentYear, de.targetDate.Month(), de.targetDate.Day(), 0, 0, 0, 0, time.Local)
			if thisYearDate.After(now) {
				// This year's date hasn't occurred yet, calculate from last year
				lastYearDate := time.Date(currentYear-1, de.targetDate.Month(), de.targetDate.Day(), 0, 0, 0, 0, time.Local)
				absDays = int(now.Sub(lastYearDate).Hours() / 24)
			} else {
				absDays = int(now.Sub(thisYearDate).Hours() / 24)
			}
		}

		if absDays == 1 {
			displayText = fmt.Sprintf("%s - 过去了 1 天", de.Name)
		} else {
			displayText = fmt.Sprintf("%s - 过去了 %d 天", de.Name, absDays)
		}
		color = "gray"
	}

	return displayText, color
}

// DaysUntil returns the number of days until the event (negative if past)
func (de *DaysEvent) DaysUntil() int {
	now := time.Now()
	currentYear := now.Year()

	var targetDate time.Time
	if de.Yearly {
		targetDate = time.Date(currentYear, de.targetDate.Month(), de.targetDate.Day(), 0, 0, 0, 0, time.Local)
		if targetDate.Before(now) {
			targetDate = time.Date(currentYear+1, de.targetDate.Month(), de.targetDate.Day(), 0, 0, 0, 0, time.Local)
		}
	} else {
		targetDate = de.targetDate
	}

	diff := targetDate.Sub(now)
	return int(diff.Hours() / 24)
}
