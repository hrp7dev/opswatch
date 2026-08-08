package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

const (
	chartRows = 3
)

var (
	cpuHistory []float64
	ramHistory []float64
)

func AddCPUHistory(value float64) {
	cpuHistory = pushHistory(cpuHistory, value)
}
func AddRAMHistory(value float64) {
	ramHistory = pushHistory(ramHistory, value)
}

// Keep a rolling buffer
func pushHistory(h []float64, value float64) []float64 {
	maxSamples := 400
	if len(h) >= maxSamples {
		h = h[1:]
	}
	return append(h, value)
}

var sparkChars = []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

func barColorFor(value float64) lipgloss.Color {
	switch {
	case value >= 90:
		return lipgloss.Color("#EF4444")
	case value >= 70:
		return lipgloss.Color("#EAB308")
	default:
		return lipgloss.Color("#22C55E")
	}
}
func renderChartPanel(title string, accent lipgloss.Color, history []float64, width int) string {
	current := 0.0
	if len(history) > 0 {
		current = history[len(history)-1]
	}
	titleLine := lipgloss.NewStyle().Bold(true).Foreground(accent).Render(title)
	valueLine := valueStyle.Render(fmt.Sprintf("%.1f%%", current))
	grid := renderGrid(history, width)
	minV, maxV, avgV := statsOf(history)
	statsLine := infoStyle.Render(fmt.Sprintf("min %.0f%%  avg %.0f%%  max %.0f%%", minV, avgV, maxV))
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Bottom, titleLine, "  ", valueLine),
		"",
		grid,
		statsLine,
	)
}
func resampleToWidth(data []float64, width int) []float64 {
	if width < 1 {
		return []float64{}
	}
	out := make([]float64, width)
	n := len(data)
	if n == 0 {
		return out
	}
	if n == 1 {
		for i := range out {
			out[i] = data[0]
		}
		return out
	}
	for i := 0; i < width; i++ {
		pos := float64(i) * float64(n-1) / float64(width-1)
		lo := int(pos)
		hi := lo + 1
		if hi >= n {
			out[i] = data[n-1]
			continue
		}
		frac := pos - float64(lo)
		out[i] = data[lo]*(1-frac) + data[hi]*frac
	}
	return out
}
func renderGrid(history []float64, width int) string {
	data := resampleToWidth(history, width)
	var rows []string
	for r := chartRows; r >= 1; r-- {
		levelFloor := float64(r-1) * (100.0 / float64(chartRows))
		levelCeil := float64(r) * (100.0 / float64(chartRows))
		var row strings.Builder
		for i := 0; i < width; i++ {
			v := data[i]
			style := lipgloss.NewStyle().Foreground(barColorFor(v))
			switch {
			case v >= levelCeil:
				row.WriteString(style.Render("█"))
			case v > levelFloor:
				frac := (v - levelFloor) / (levelCeil - levelFloor)
				idx := int(frac * float64(len(sparkChars)-1))
				row.WriteString(style.Render(string(sparkChars[idx])))
			default:
				row.WriteString(emptyBarStyle.Render("·"))
			}
		}
		rows = append(rows, row.String())
	}
	rows = append(rows, strings.Repeat("─", width))
	return strings.Join(rows, "\n")
}
func statsOf(history []float64) (min, max, avg float64) {
	if len(history) == 0 {
		return 0, 0, 0
	}
	min, max = history[0], history[0]
	sum := 0.0
	for _, v := range history {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	avg = sum / float64(len(history))
	return
}
func RenderCPUChart(width int) string {
	return renderChartPanel("CPU", lipgloss.Color("#38BDF8"), cpuHistory, width)
}
func RenderRAMChart(width int) string {
	return renderChartPanel("RAM", lipgloss.Color("#A78BFA"), ramHistory, width)
}
