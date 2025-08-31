package daysmatter

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable = false
	defaultTitle     = "Days Matter"
)

// Event represents a single memorable event
type Event struct {
	Name   string `yaml:"name" json:"name"`
	Date   string `yaml:"date" json:"date"`
	Yearly bool   `yaml:"yearly" json:"yearly"`
}

// Settings defines the configuration properties for this module
type Settings struct {
	*cfg.Common

	events []Event
}

// NewSettingsFromYAML creates a new settings instance from a YAML config block
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	common := cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig)

	settings := Settings{
		Common: common,
		events: []Event{},
	}

	// Parse events from config
	if eventsConfig := ymlConfig.UList("events"); eventsConfig != nil {
		for _, eventInterface := range eventsConfig {
			if eventMap, ok := eventInterface.(map[string]interface{}); ok {
				event := Event{
					Name:   getString(eventMap, "name"),
					Date:   getString(eventMap, "date"),
					Yearly: getBool(eventMap, "yearly"),
				}
				settings.events = append(settings.events, event)
			}
		}
	}

	return &settings
}

// Helper functions to safely extract values from interface{} map
func getString(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if val, exists := m[key]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
