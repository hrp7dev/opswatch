package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hrp7dev/opswatch/internal/system"
)

func RenderLogs(limit int) string {
	logs := system.GetSystemLogs(limit)

	if len(logs) == 0 {
		return infoStyle.Render("No system logs")
	}

	var rows []string

	for _, log := range logs {

		levelStyle := successStyle

		switch log.Level {
		case "ERROR":
			levelStyle = dangerStyle
		case "WARN":
			levelStyle = warningStyle
		}

		rows = append(rows,
			fmt.Sprintf(
				"%s %s",
				levelStyle.Render(padLogCell(log.Level, 7)),
				log.Message,
			),
		)
	}

	return strings.Join(rows, "\n")
}

func padLogCell(value string, width int) string {
	if lipgloss.Width(value) >= width {
		return value
	}

	return value + strings.Repeat(" ", width-lipgloss.Width(value))
}
