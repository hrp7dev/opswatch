package main

import (
	"os"

	"golang.org/x/sys/unix"
)

func restoreTerminal() {
	fd := int(os.Stdin.Fd())

	state, err := unix.IoctlGetTermios(fd, ioctlGetTermios)
	if err != nil {
		return
	}

	state.Lflag |= unix.ICANON
	state.Lflag |= unix.ECHO

	unix.IoctlSetTermios(
		fd,
		ioctlSetTermios,
		state,
	)
}
