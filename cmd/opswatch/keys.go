package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"

	"github.com/hrp7dev/opswatch/internal/ui"
)

func listenKeys() {
	fd := int(os.Stdin.Fd())

	oldState, err := unix.IoctlGetTermios(fd, ioctlGetTermios)
	if err != nil {
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

	debugLog, _ := os.OpenFile("/tmp/opswatch_keys.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer func() {
		if debugLog != nil {
			debugLog.Close()
		}
	}()

	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		if debugLog != nil {
			fmt.Fprintf(debugLog, "n=%d bytes=%v\n", n, buf[:n])
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
			if debugLog != nil {
				fmt.Fprintf(debugLog, "EXIT TRIGGERED by byte %d\n", buf[0])
			}
			unix.IoctlSetTermios(fd, ioctlSetTermios, oldState)
			ui.ExitAltScreen()
			os.Exit(0)
		}
	}
}
