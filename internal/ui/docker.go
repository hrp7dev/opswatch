package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrp7dev/opswatch/internal/system"
)

const (
	colIdx    = 3
	colID     = 20
	colName   = 20
	colStatus = 18
	colCPU    = 8
	colMem    = 20
)

func padCell(s string, width int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	w := lipgloss.Width(s)

	if w >= width {
		return s[:width]
	}

	return s + strings.Repeat(" ", width-w)
}

func renderDockerPanel(width int) string {
	info, containers := getDockerCache()

	if !info.Available && len(containers) == 0 {
		return infoStyle.Render("Checking Docker...")
	}

	var mem uint64

	for _, c := range containers {
		mem += c.MemBytes
	}

	header := fmt.Sprintf(
		"Version:%s │ CPU:%d │ MEM:%s │ Containers:%d/%d",
		info.ServerVersion,
		info.CPUs,
		system.FormatBytes(mem),
		info.RunningContainers,
		info.TotalContainers,
	)

	if len(containers) == 0 {
		return infoStyle.Render(header)
	}

	dockerCacheMu.RLock()
	page := dockerPage
	maxPage := dockerMaxPageLocked()
	dockerCacheMu.RUnlock()

	start := page * dockerPageSize

	if start > len(containers) {
		start = 0
	}

	end := start + dockerPageSize

	if end > len(containers) {
		end = len(containers)
	}

	pageItems := containers[start:end]

	colIdx := 3
	colID := 12
	colName := 20
	colStatus := 18
	colCPU := 8
	colMem := 12

	total := colIdx + colID + colName + colStatus + colCPU + colMem + 5

	if total > width-4 {

		colID = 8
		colName = 14
		colStatus = 12
		colCPU = 6
		colMem = 10

		total = colIdx + colID + colName + colStatus + colCPU + colMem + 5
	}

	headerLine := strings.Join([]string{
		padCell(tableHeaderStyle.Render("#"), colIdx),
		padCell(tableHeaderStyle.Render("ID"), colID),
		padCell(tableHeaderStyle.Render("NAME"), colName),
		padCell(tableHeaderStyle.Render("STATUS"), colStatus),
		padCell(tableHeaderStyle.Render("CPU"), colCPU),
		padCell(tableHeaderStyle.Render("MEM"), colMem),
	}, " ")

	rows := []string{
		headerLine,
		rowSepStyle.Render(strings.Repeat("─", total)),
	}

	for i, c := range pageItems {

		cpuStyle := successStyle

		if c.CPUPercent >= 70 {
			cpuStyle = warningStyle
		}

		if c.CPUPercent >= 90 {
			cpuStyle = dangerStyle
		}

		row := strings.Join([]string{
			padCell(fmt.Sprintf("%d", start+i+1), colIdx),
			padCell(truncate(c.ID, colID-1), colID),
			padCell(truncate(c.Name, colName-1), colName),
			padCell(statusStyle(c.Status).Render(truncate(c.Status, colStatus-1)), colStatus),
			padCell(cpuStyle.Render(truncate(c.CPUUsage, colCPU-1)), colCPU),
			padCell(truncate(c.MemUsage, colMem-1), colMem),
		}, " ")

		rows = append(rows, row)
		rows = append(rows, rowSepStyle.Render(strings.Repeat("─", total)))
	}

	rows = append(
		rows,
		pageInfoStyle.Render(
			fmt.Sprintf(
				"Page %d/%d (%d containers) ←/p prev →/n next",
				page+1,
				maxPage+1,
				len(containers),
			),
		),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		infoStyle.Render(header),
		strings.Join(rows, "\n"),
	)
}

func statusStyle(status string) lipgloss.Style {
	status = strings.ToLower(status)

	if strings.Contains(status, "up") {
		return successStyle
	}

	if strings.Contains(status, "exit") || strings.Contains(status, "dead") {
		return dangerStyle
	}

	return warningStyle
}

func truncate(value string, width int) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.TrimSpace(value)

	if lipgloss.Width(value) <= width {
		return value
	}

	return value[:width-1] + "…"
}
