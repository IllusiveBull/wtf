package resourceusage

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable       = false
	defaultRefreshInterval = "1s"
	defaultTitle           = "ResourceUsage"
)

type Settings struct {
	*cfg.Common

	cpuCombined  bool
	showCPU      bool
	showMem      bool
	showSwp      bool
	showNet      bool
	netInterface string
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		cpuCombined:  ymlConfig.UBool("cpuCombined", false),
		showCPU:      ymlConfig.UBool("showCPU", true),
		showMem:      ymlConfig.UBool("showMem", true),
		showSwp:      ymlConfig.UBool("showSwp", true),
		showNet:      ymlConfig.UBool("showNet", true),
		netInterface: ymlConfig.UString("netInterface", ""),
	}
	settings.RefreshInterval = cfg.ParseTimeString(ymlConfig, "refreshInterval", defaultRefreshInterval)

	return &settings
}
