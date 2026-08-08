package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hrp7dev/opswatch/internal/ui"
)

func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ui.ExitAltScreen()
		os.Exit(0)
	}()
}
