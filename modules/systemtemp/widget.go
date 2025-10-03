package systemtemp

import (
	"fmt"
	"strings"

	"github.com/dkorunic/iSMC/smc"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
)

// Widget define wtf widget to register widget later
type Widget struct {
	settings *Settings
	tviewApp *tview.Application
	view.ScrollableWidget
}

// NewWidget Make new instance of widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),

		tviewApp: tviewApp,
		settings: settings,
	}

	widget.View.SetWrap(true)
	widget.View.SetWordWrap(true)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh & update after interval time
func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	widget.View.Clear()
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) display() {
	widget.Redraw(widget.content)
}

func (widget *Widget) content() (string, string, bool) {
	temps := getTemperatures(widget.settings)

	if len(temps) == 0 {
		errMsg := "[red]无法获取温度信息[white]\n\n"
		errMsg += "可能的原因:\n"
		errMsg += "1. 需要管理员权限\n"
		errMsg += "2. SMC 访问被限制\n"
		errMsg += "3. 不支持的硬件型号"
		return widget.CommonSettings().Title, errMsg, false
	}

	content := ""
	for _, temp := range temps {
		color := getColorForTemp(temp.Value)
		content += fmt.Sprintf("[%s]%-10s[white] %s\n", color, temp.Label+":", temp.Display)
	}

	return widget.CommonSettings().Title, content, false
}

type TempInfo struct {
	Label   string
	Value   float64
	Display string
}

func getTemperatures(settings *Settings) []TempInfo {
	var temps []TempInfo

	// 获取所有温度传感器
	allTemps := smc.GetTemperature()
	if len(allTemps) == 0 {
		return temps
	}

	// CPU 温度 - 取 CPU Performance Core 的平均值
	if settings.showCPU {
		cpuTemps := []float64{}
		for key, value := range allTemps {
			keyLower := strings.ToLower(key)
			if strings.Contains(keyLower, "cpu") &&
				(strings.Contains(keyLower, "performance") || strings.Contains(keyLower, "efficiency")) &&
				!strings.Contains(keyLower, "proximity") {
				if temp, ok := extractTempFromSMCValue(value); ok {
					cpuTemps = append(cpuTemps, temp)
				}
			}
		}
		if len(cpuTemps) > 0 {
			avgTemp := average(cpuTemps)
			temps = append(temps, TempInfo{
				Label:   "CPU",
				Value:   avgTemp,
				Display: fmt.Sprintf("%.1f°C", avgTemp),
			})
		}
	}

	// GPU 温度
	if settings.showGPU {
		for key, value := range allTemps {
			keyLower := strings.ToLower(key)
			if strings.Contains(keyLower, "gpu") && !strings.Contains(keyLower, "heatsink") {
				if temp, ok := extractTempFromSMCValue(value); ok {
					temps = append(temps, TempInfo{
						Label:   "GPU",
						Value:   temp,
						Display: fmt.Sprintf("%.1f°C", temp),
					})
					break
				}
			}
		}
	}

	// SSD/NAND 温度
	if settings.showSSD {
		for key, value := range allTemps {
			keyLower := strings.ToLower(key)
			if strings.Contains(keyLower, "nand") ||
				(strings.Contains(keyLower, "drive") && !strings.Contains(keyLower, "proximity")) {
				if temp, ok := extractTempFromSMCValue(value); ok {
					temps = append(temps, TempInfo{
						Label:   "SSD",
						Value:   temp,
						Display: fmt.Sprintf("%.1f°C", temp),
					})
					break
				}
			}
		}
	}

	// Memory 温度
	if settings.showMem {
		for key, value := range allTemps {
			keyLower := strings.ToLower(key)
			if strings.Contains(keyLower, "memory") && strings.Contains(keyLower, "proximity") {
				if temp, ok := extractTempFromSMCValue(value); ok {
					temps = append(temps, TempInfo{
						Label:   "Memory",
						Value:   temp,
						Display: fmt.Sprintf("%.1f°C", temp),
					})
					break
				}
			}
		}
	}

	// 电池温度
	if settings.showBattery {
		battTemps := []float64{}
		for key, value := range allTemps {
			keyLower := strings.ToLower(key)
			if strings.Contains(keyLower, "battery") {
				if temp, ok := extractTempFromSMCValue(value); ok {
					battTemps = append(battTemps, temp)
				}
			}
		}
		if len(battTemps) > 0 {
			avgTemp := average(battTemps)
			temps = append(temps, TempInfo{
				Label:   "Battery",
				Value:   avgTemp,
				Display: fmt.Sprintf("%.1f°C", avgTemp),
			})
		}
	}

	return temps
}

// extractTempFromSMCValue 从 iSMC 返回的 map[string]interface{} 中提取温度值
func extractTempFromSMCValue(value any) (float64, bool) {
	// iSMC 返回的格式: map[key:Tp01 type:flt value:74.1 °C]
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return 0, false
	}

	// 获取 value 字段
	tempValue, ok := valueMap["value"]
	if !ok {
		return 0, false
	}

	// value 是字符串格式，如 "74.1 °C"
	tempStr, ok := tempValue.(string)
	if !ok {
		return 0, false
	}

	// 解析温度值
	tempStr = strings.TrimSpace(tempStr)
	tempStr = strings.TrimSuffix(tempStr, "°C")
	tempStr = strings.TrimSuffix(tempStr, "C")
	tempStr = strings.TrimSpace(tempStr)

	var temp float64
	if _, err := fmt.Sscanf(tempStr, "%f", &temp); err != nil {
		return 0, false
	}

	return temp, true
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// getColorForTemp 根据温度返回颜色
func getColorForTemp(temp float64) string {
	if temp < 0 {
		return "gray"
	} else if temp < 50 {
		return "green"
	} else if temp < 70 {
		return "yellow"
	} else if temp < 85 {
		return "orange"
	} else {
		return "red"
	}
}
