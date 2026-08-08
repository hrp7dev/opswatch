package ui

import (
	"fmt"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

var logHistory []LogEntry

func AddLog(level, message string) {
	if len(logHistory) >= 20 {
		logHistory = logHistory[1:]
	}
	logHistory = append(logHistory, LogEntry{
		Timestamp: time.Now(),
		Level:     strings.ToUpper(level),
		Message:   message,
	})
}

func RenderLogs(maxLines int) string {
	entries := logHistory
	if len(entries) > maxLines {
		entries = entries[len(entries)-maxLines:]
	}
	if len(entries) == 0 {
		return infoStyle.Render("No logs yet. Waiting for activity...")
	}

	var rows []string
	for _, entry := range entries {
		timestamp := entry.Timestamp.Format("15:04:05")
		levelStyle := successStyle
		switch entry.Level {
		case "ERROR", "ERR":
			levelStyle = dangerStyle
		case "WARN", "WARNING":
			levelStyle = warningStyle
		}

		rows = append(rows, fmt.Sprintf(
			"%s %s %s",
			subtitleStyle.Render(timestamp),
			levelStyle.Render(entry.Level),
			entry.Message,
		))
	}

	return strings.Join(rows, "\n")
}
