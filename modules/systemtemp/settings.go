package systemtemp

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable       = false
	defaultRefreshInterval = "5s"
	defaultTitle           = "System Temperature"
)

type Settings struct {
	*cfg.Common

	showCPU     bool
	showGPU     bool
	showSSD     bool
	showMem     bool
	showBattery bool
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		showCPU:     ymlConfig.UBool("showCPU", true),
		showGPU:     ymlConfig.UBool("showGPU", true),
		showSSD:     ymlConfig.UBool("showSSD", true),
		showMem:     ymlConfig.UBool("showMem", true),
		showBattery: ymlConfig.UBool("showBattery", true),
	}
	settings.RefreshInterval = cfg.ParseTimeString(ymlConfig, "refreshInterval", defaultRefreshInterval)

	return &settings
}
