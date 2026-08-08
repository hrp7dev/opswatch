package main

import (
	"os"

	"golang.org/x/sys/unix"

	"github.com/hrp7dev/opswatch/internal/ui"
)

// listenKeys puts stdin into "cbreak" mode: no line buffering, no echo,
// no signal generation from the terminal (ISIG off) and no XON/XOFF flow
// control (IXON off) — but OPOST stays ON, so newlines still render
// correctly and the dashboard doesn't get garbled like it does under
// full raw mode.
func listenKeys() {
	fd := int(os.Stdin.Fd())

	oldState, err := unix.IoctlGetTermios(fd, ioctlGetTermios)
	if err != nil {
		// Not a real terminal (e.g. pipe/CI) — skip silently.
		return
	}

	newState := *oldState
	newState.Lflag &^= unix.ICANON | unix.ECHO | unix.ISIG
	newState.Iflag &^= unix.IXON
	newState.Cc[unix.VMIN] = 1
	newState.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(fd, ioctlSetTermios, &newState); err != nil {
		return
	}
	defer unix.IoctlSetTermios(fd, ioctlSetTermios, oldState)

	readByte := func() (byte, bool) {
		b := make([]byte, 1)
		n, err := os.Stdin.Read(b)
		if err != nil || n == 0 {
			return 0, false
		}
		return b[0], true
	}

	for {
		b1, ok := readByte()
		if !ok {
			continue
		}

		switch b1 {
		case 'n', 'l':
			ui.NextDockerPage()
			continue

		case 'p', 'h':
			ui.PrevDockerPage()
			continue

		case 'q', 3: // 'q' or Ctrl+C
			unix.IoctlSetTermios(fd, ioctlSetTermios, oldState)
			ui.ExitAltScreen()
			os.Exit(0)

		case 27: // ESC — could be a lone Escape or the start of an
			// arrow-key sequence. Arrow keys always send more bytes
			// right behind it, so block for the next byte instead of
			// requiring all bytes to land in a single Read() — over
			// SSH they often arrive in separate reads.
			b2, ok := readByte()
			if !ok {
				continue
			}

			// Normal mode: ESC [ C / ESC [ D
			// Application mode (some terminals/SSH setups): ESC O C / ESC O D
			if b2 != '[' && b2 != 'O' {
				continue
			}

			b3, ok := readByte()
			if !ok {
				continue
			}

			switch b3 {
			case 'C': // Right arrow
				ui.NextDockerPage()
			case 'D': // Left arrow
				ui.PrevDockerPage()
			}
		}
	}
}