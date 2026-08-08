package main

import (
	"os"

	"golang.org/x/sys/unix"

	"github.com/hrp7dev/opswatch/internal/ui"
)

// listenKeys puts stdin into "cbreak" mode: no line buffering, no echo —
// but OPOST stays ON, so newlines still render correctly and the dashboard
// doesn't get garbled like it does under full raw mode.
func listenKeys() {
	fd := int(os.Stdin.Fd())

	oldState, err := unix.IoctlGetTermios(fd, ioctlGetTermios)
	if err != nil {
		// Not a real terminal (e.g. pipe/CI) — skip silently.
		return
	}

	newState := *oldState
	newState.Lflag &^= unix.ICANON | unix.ECHO
	newState.Cc[unix.VMIN] = 1
	newState.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(fd, ioctlSetTermios, &newState); err != nil {
		return
	}
	defer unix.IoctlSetTermios(fd, ioctlSetTermios, oldState)

	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		switch {
		case buf[0] == 'n', buf[0] == 'l':
			ui.NextDockerPage()

		case buf[0] == 'p', buf[0] == 'h':
			ui.PrevDockerPage()

		case n == 3 && buf[0] == 27 && buf[1] == 91 && buf[2] == 67: // Right arrow
			ui.NextDockerPage()

		case n == 3 && buf[0] == 27 && buf[1] == 91 && buf[2] == 68: // Left arrow
			ui.PrevDockerPage()

		case buf[0] == 'q', buf[0] == 3: // 'q' or Ctrl+C
			unix.IoctlSetTermios(fd, ioctlSetTermios, oldState)
			ui.ExitAltScreen()
			os.Exit(0)
		}
	}
}
