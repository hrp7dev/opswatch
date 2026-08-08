package system

import (
	"os/exec"
	"strings"
)

type SystemLog struct {
	Time    string
	Level   string
	Message string
}

func GetSystemLogs(limit int) []SystemLog {
	var logs []SystemLog

	commands := [][]string{
		{"journalctl", "-n", "50", "--no-pager"},
		{"tail", "-n", "50", "/var/log/syslog"},
		{"tail", "-n", "50", "/var/log/messages"},
	}

	var output string

	for _, cmd := range commands {
		out, err := exec.Command(cmd[0], cmd[1:]...).Output()

		if err == nil {
			output = string(out)
			break
		}
	}

	if output == "" {
		return []SystemLog{
			{
				Time:    "--:--:--",
				Level:   "WARN",
				Message: "Unable to read system logs",
			},
		}
	}

	lines := strings.Split(output, "\n")

	start := 0
	if len(lines) > limit {
		start = len(lines) - limit
	}

	for _, line := range lines[start:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		level := "INFO"

		lower := strings.ToLower(line)

		if strings.Contains(lower, "error") ||
			strings.Contains(lower, "failed") {
			level = "ERROR"
		}

		if strings.Contains(lower, "warn") {
			level = "WARN"
		}

		logs = append(logs, SystemLog{
			Level:   level,
			Message: line,
		})
	}

	return logs
}
