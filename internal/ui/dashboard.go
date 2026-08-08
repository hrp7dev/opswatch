package ui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrp7dev/opswatch/internal/system"
	"golang.org/x/term"
)

var (
	titleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38BDF8"))
	subtitleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
	labelStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CBD5E1"))
	valueStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8FAFC"))
	infoStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
	panelStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#334155")).Padding(0, 1)
	successStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E"))
	warningStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308"))
	dangerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	emptyBarStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E293B"))
	pageInfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Italic(true)
	tableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CBD5E1"))
	rowSepStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E293B"))
	idStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B"))
)

var (
	dockerCacheMu    sync.RWMutex
	cachedDocker     system.DockerStats
	cachedContainers []system.ContainerStats
	dockerPage       int
	dockerPageSize   = 4
)

func SetDockerCache(info system.DockerStats, containers []system.ContainerStats) {
	dockerCacheMu.Lock()
	cachedDocker = info
	cachedContainers = containers

	maxPage := dockerMaxPageLocked()

	if dockerPage > maxPage {
		dockerPage = maxPage
	}

	dockerCacheMu.Unlock()
}

func getDockerCache() (system.DockerStats, []system.ContainerStats) {
	dockerCacheMu.RLock()
	defer dockerCacheMu.RUnlock()

	return cachedDocker, cachedContainers
}

func dockerMaxPageLocked() int {
	total := len(cachedContainers)
	if total == 0 {
		return 0
	}
	return (total - 1) / dockerPageSize
}

func NextDockerPage() {
	dockerCacheMu.Lock()
	defer dockerCacheMu.Unlock()

	max := dockerMaxPageLocked()

	if dockerPage < max {
		dockerPage++
	}
}

func PrevDockerPage() {
	dockerCacheMu.Lock()
	defer dockerCacheMu.Unlock()

	if dockerPage > 0 {
		dockerPage--
	}
}

func ResetDockerPage() {
	dockerCacheMu.Lock()
	defer dockerCacheMu.Unlock()

	dockerPage = 0
}

func RenderDashboard(cpuUsage float64, memory system.MemoryStats, disk system.DiskStats) {
	Clear()

	width, height, err := term.GetSize(0)

	if err != nil {
		width = 90
		height = 50
	}

	if width < 90 || height < 50 {
		errorStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EF4444"))

		labelErrorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FACC15"))

		message := lipgloss.JoinVertical(
			lipgloss.Center,
			errorStyle.Render("⚠ OPSWATCH TERMINAL SIZE ERROR"),
			"",
			labelErrorStyle.Render(fmt.Sprintf(
				"Current Size:\nWidth  : %d\nHeight : %d",
				width,
				height,
			)),
			"",
			labelErrorStyle.Render(
				"Required Size:\nWidth  : 90+\nHeight : 50+",
			),
			"",
			infoStyle.Render("Please resize your terminal."),
		)

		fmt.Println(
			lipgloss.NewStyle().
				Width(width).
				Height(height).
				Align(
					lipgloss.Center,
					lipgloss.Center,
				).
				Render(message),
		)

		return
	}

	dashboardWidth := width - 4

	creatorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#38BDF8"))

	creatorLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94A3B8"))

	header := lipgloss.NewStyle().
		Width(dashboardWidth).
		Align(lipgloss.Center).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				titleStyle.Render("⚡ OPSWATCH"),
				subtitleStyle.Render("Real-time Server Monitoring"),
				infoStyle.Render("CPU • MEMORY • DISK • DOCKER • LIVE METRICS"),
				"",
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					creatorLabelStyle.Render("Created by "),
					creatorStyle.Render(
						terminalLink(
							"Hrp7dev",
							"https://github.com/hrp7dev",
						),
					),
				),
			),
		)

	cardWidth := (dashboardWidth - 2) / 3

	cardStyle := func(content string) string {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#334155")).
			Padding(0, 2).
			Width(cardWidth).
			Render(content)
	}

	stats := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cardStyle(renderMetric("CPU", cpuUsage, fmt.Sprintf("%d Cores", system.CPUCores()), cardWidth)),
		cardStyle(renderMetric("MEMORY", memory.Usage, fmt.Sprintf("%s / %s", system.FormatBytes(memory.Used), system.FormatBytes(memory.Total)), cardWidth)),
		cardStyle(renderMetric("DISK", disk.Usage, fmt.Sprintf("%s / %s", system.FormatBytes(disk.Used), system.FormatBytes(disk.Total)), cardWidth)),
	)

	chartWidth := int(float64(dashboardWidth) / 2.02)

	cpuChart := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(0, 2).
		Width(chartWidth).
		Render(RenderCPUChart(chartWidth - 4))

	ramChart := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(0, 2).
		Width(chartWidth).
		Render(RenderRAMChart(chartWidth - 4))

	charts := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpuChart,
		ramChart,
	)

	docker := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(0, 2).
		Width(dashboardWidth).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				labelStyle.Render("DOCKER OVERVIEW"),
				renderDockerPanel(dashboardWidth),
			),
		)

	logs := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(0, 2).
		Width(dashboardWidth).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				labelStyle.Render("EVENT LOGS"),
				RenderLogs(10),
			),
		)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		stats,
		charts,
		docker,
		logs,
	)

	fmt.Println(
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(body),
	)
}

func renderMetric(label string, percentage float64, detail string, width int) string {
	value := valueStyle.Render(
		fmt.Sprintf("%5.1f%%", percentage),
	)

	barWidth := width - lipgloss.Width(value) - 4

	if barWidth < 10 {
		barWidth = 10
	}

	bar := progressBar(percentage, barWidth)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render(label),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			bar,
			" ",
			value,
		),
		infoStyle.Render(detail),
	)
}

func cardWidthForMetric() int {
	return 18
}

func progressBar(percentage float64, width int) string {
	if percentage < 0 {
		percentage = 0
	}

	if percentage > 100 {
		percentage = 100
	}

	filled := int((percentage / 100) * float64(width))

	style := successStyle

	if percentage >= 70 {
		style = warningStyle
	}

	if percentage >= 90 {
		style = dangerStyle
	}

	return style.Render(strings.Repeat("█", filled)) +
		emptyBarStyle.Render(strings.Repeat("░", width-filled))
}

func terminalLink(text string, url string) string {
	return fmt.Sprintf(
		"\033]8;;%s\033\\%s\033]8;;\033\\",
		url,
		text,
	)
}
