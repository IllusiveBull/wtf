package daysmatter

import (
	"fmt"
)

// display renders the widget's content
func (widget *Widget) display() {
	widget.Redraw(widget.content)
}

// content generates the display content for the widget
func (widget *Widget) content() (string, string, bool) {
	title := widget.CommonSettings().Title

	if widget.Error != "" {
		return title, widget.Error, true
	}

	if len(widget.events) == 0 {
		return title, "No events configured", false
	}

	var str string

	for _, event := range widget.events {
		displayText, color := event.GetDisplayInfo()

		line := fmt.Sprintf("[%s]%s[white]\n", color, displayText)
		str += line
	}

	return title, str, false
}
