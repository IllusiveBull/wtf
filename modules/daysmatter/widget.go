package daysmatter

import (
	"sort"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
)

// Widget represents a DaysMatter widget
type Widget struct {
	view.TextWidget

	settings *Settings
	events   []*DaysEvent
	tviewApp *tview.Application
	Error    string
}

// NewWidget creates a new instance of a widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, nil, settings.Common),

		settings: settings,
		tviewApp: tviewApp,
	}

	widget.loadEvents()
	widget.View.SetWrap(true)
	widget.View.SetWordWrap(true)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the widget's content
func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	widget.loadEvents()
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

// loadEvents loads and parses events from settings
func (widget *Widget) loadEvents() {
	widget.events = []*DaysEvent{}
	widget.Error = ""

	for _, event := range widget.settings.events {
		daysEvent, err := NewDaysEvent(event)
		if err != nil {
			widget.Error = err.Error()
			return
		}
		widget.events = append(widget.events, daysEvent)
	}

	// Sort events by days until (closest events first)
	sort.Slice(widget.events, func(i, j int) bool {
		daysI := widget.events[i].DaysUntil()
		daysJ := widget.events[j].DaysUntil()

		// Put future events before past events
		if daysI >= 0 && daysJ < 0 {
			return true
		}
		if daysI < 0 && daysJ >= 0 {
			return false
		}

		// For events in the same category (both future or both past)
		if daysI >= 0 && daysJ >= 0 {
			// Both future: closer events first
			return daysI < daysJ
		} else {
			// Both past: more recent events first
			return daysI > daysJ
		}
	})
}
