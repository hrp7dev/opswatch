package ui

import (
	"os"

	"golang.org/x/term"
)

// TerminalSize returns the current terminal width and height in columns/rows.
// Falls back to a sane default if it can't be determined (e.g. not a TTY).
func TerminalSize() (width, height int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 || h <= 0 {
		return 110, 40
	}
	return w, h
}
