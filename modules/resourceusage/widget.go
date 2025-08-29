package resourceusage

import (
	"fmt"
	"math"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/wtfutil/wtf/view"
)

// Widget define wtf widget to register widget later
type Widget struct {
	settings *Settings
	tviewApp *tview.Application
	view.BarGraph

	// Network tracking
	lastNetStats     map[string]net.IOCountersStat
	lastNetTime      time.Time
	netUploadSpeed   float64
	netDownloadSpeed float64
}

// NewWidget Make new instance of widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, settings *Settings) *Widget {
	widget := Widget{
		BarGraph: view.NewBarGraph(tviewApp, redrawChan, settings.Name, settings.Common),

		tviewApp:     tviewApp,
		settings:     settings,
		lastNetStats: make(map[string]net.IOCountersStat),
		lastNetTime:  time.Now(),
	}

	widget.View.SetWrap(false)
	widget.View.SetWordWrap(false)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// MakeGraph - Load the dead drop stats
func MakeGraph(widget *Widget) {
	cpuStats, memInfo := getDataFromSystem(widget)

	// Update network stats
	widget.updateNetworkStats()

	var itemsCount = 0
	if widget.settings.showCPU {
		itemsCount += len(cpuStats)
	}

	if widget.settings.showMem {
		itemsCount++
	}

	if widget.settings.showSwp {
		itemsCount++
	}

	if widget.settings.showNet {
		itemsCount += 2 // Upload and Download
	}

	var stats = make([]view.Bar, itemsCount)
	var nextIndex = 0

	if widget.settings.showCPU && len(cpuStats) > 0 {
		for i, stat := range cpuStats {
			// Stats sometimes jump outside the 0-100 range, possibly due to timing
			stat = math.Min(100, stat)
			stat = math.Max(0, stat)

			var label string
			if widget.settings.cpuCombined {
				label = "CPU"
			} else {
				label = fmt.Sprint(i)
			}

			bar := view.Bar{
				Label:      label,
				Percent:    int(stat),
				ValueLabel: fmt.Sprintf("%d%%", int(stat)),
				LabelColor: "red",
			}

			stats[nextIndex] = bar
			nextIndex++
		}
	}

	if widget.settings.showMem {
		usedMemLabel := bytefmt.ByteSize(memInfo.Used)
		totalMemLabel := bytefmt.ByteSize(memInfo.Total)

		if usedMemLabel[len(usedMemLabel)-1] == totalMemLabel[len(totalMemLabel)-1] {
			usedMemLabel = usedMemLabel[:len(usedMemLabel)-1]
		}

		stats[nextIndex] = view.Bar{
			Label:      "Mem",
			Percent:    int(memInfo.UsedPercent),
			ValueLabel: fmt.Sprintf("%s/%s", usedMemLabel, totalMemLabel),
			LabelColor: "green",
		}
		nextIndex++
	}

	if widget.settings.showSwp {
		swapUsed := memInfo.SwapTotal - memInfo.SwapFree
		var swapPercent float64
		if memInfo.SwapTotal > 0 {
			swapPercent = float64(swapUsed) / float64(memInfo.SwapTotal)
		}

		usedSwapLabel := bytefmt.ByteSize(swapUsed)
		totalSwapLabel := bytefmt.ByteSize(memInfo.SwapTotal)

		if usedSwapLabel[len(usedSwapLabel)-1] == totalSwapLabel[len(totalSwapLabel)-1] {
			usedSwapLabel = usedSwapLabel[:len(usedSwapLabel)-1]
		}

		stats[nextIndex] = view.Bar{
			Label:      "Swp",
			Percent:    int(swapPercent * 100),
			ValueLabel: fmt.Sprintf("%s/%s", usedSwapLabel, totalSwapLabel),
			LabelColor: "yellow",
		}
		nextIndex++
	}

	if widget.settings.showNet {
		// Network Download
		downloadSpeed := widget.netDownloadSpeed
		downloadLabel := formatNetworkSpeed(downloadSpeed)

		stats[nextIndex] = view.Bar{
			Label:      "↓Net",
			Percent:    int(math.Min(100, downloadSpeed/1024/1024*10)), // Scale: 10MB/s = 100%
			ValueLabel: downloadLabel,
			LabelColor: "blue",
		}
		nextIndex++

		// Network Upload
		uploadSpeed := widget.netUploadSpeed
		uploadLabel := formatNetworkSpeed(uploadSpeed)

		stats[nextIndex] = view.Bar{
			Label:      "↑Net",
			Percent:    int(math.Min(100, uploadSpeed/1024/1024*10)), // Scale: 10MB/s = 100%
			ValueLabel: uploadLabel,
			LabelColor: "cyan",
		}
	}

	widget.BuildBars(stats)

}

// Refresh & update after interval time
func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	widget.View.Clear()
	MakeGraph(widget)
}

/* -------------------- Unexported Functions -------------------- */

func getDataFromSystem(widget *Widget) (cpuStats []float64, memInfo mem.VirtualMemoryStat) {
	if widget.settings.showCPU {
		rCPUStats, err := cpu.Percent(time.Duration(0), !widget.settings.cpuCombined)
		if err == nil {
			cpuStats = rCPUStats
		}
	}

	if widget.settings.showMem || widget.settings.showSwp {
		rMemInfo, err := mem.VirtualMemory()
		if err == nil {
			memInfo = *rMemInfo
		}
	}

	return cpuStats, memInfo
}

func (widget *Widget) updateNetworkStats() {
	if !widget.settings.showNet {
		return
	}

	var netStats []net.IOCountersStat
	var err error

	if widget.settings.netInterface != "" {
		// Get stats for specific interface
		netStats, err = net.IOCounters(true)
		if err != nil {
			return
		}

		// Filter for the specified interface
		var filteredStats []net.IOCountersStat
		for _, stat := range netStats {
			if stat.Name == widget.settings.netInterface {
				filteredStats = append(filteredStats, stat)
				break
			}
		}
		netStats = filteredStats
	} else {
		// Get total stats for all interfaces
		netStats, err = net.IOCounters(false)
		if err != nil {
			return
		}
	}

	if len(netStats) == 0 {
		return
	}

	currentTime := time.Now()
	currentStats := netStats[0]

	// Calculate speed if we have previous stats
	if lastStats, exists := widget.lastNetStats[currentStats.Name]; exists {
		timeDiff := currentTime.Sub(widget.lastNetTime).Seconds()
		if timeDiff > 0 {
			bytesRecvDiff := currentStats.BytesRecv - lastStats.BytesRecv
			bytesSentDiff := currentStats.BytesSent - lastStats.BytesSent

			widget.netDownloadSpeed = float64(bytesRecvDiff) / timeDiff
			widget.netUploadSpeed = float64(bytesSentDiff) / timeDiff
		}
	}

	// Store current stats for next calculation
	widget.lastNetStats[currentStats.Name] = currentStats
	widget.lastNetTime = currentTime
}

func formatNetworkSpeed(bytesPerSecond float64) string {
	if bytesPerSecond < 1024 {
		return fmt.Sprintf("%.0f B/s", bytesPerSecond)
	} else if bytesPerSecond < 1024*1024 {
		return fmt.Sprintf("%.1f KB/s", bytesPerSecond/1024)
	} else if bytesPerSecond < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB/s", bytesPerSecond/1024/1024)
	} else {
		return fmt.Sprintf("%.1f GB/s", bytesPerSecond/1024/1024/1024)
	}
}
