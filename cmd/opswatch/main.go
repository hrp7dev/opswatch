package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hrp7dev/opswatch/internal/system"
	"github.com/hrp7dev/opswatch/internal/ui"
	"github.com/hrp7dev/opswatch/internal/version"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("OpsWatch", version.Version)
		return
	}

	ui.EnterAltScreen()

	defer func() {
		ui.ExitAltScreen()
		restoreTerminal()
	}()

	setupSignalHandler()

	go func() {
		for {
			info, containers, _ := system.DockerInfo()
			ui.SetDockerCache(info, containers)
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				restoreTerminal()
				ui.ExitAltScreen()
				fmt.Fprintf(os.Stderr, "panic in listenKeys: %v\n", r)
				os.Exit(1)
			}
		}()
		listenKeys()
	}()

	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					restoreTerminal()
					ui.ExitAltScreen()
					fmt.Fprintf(os.Stderr, "panic: %v\n", r)
					os.Exit(1)
				}
			}()

			cpuUsage, err := system.CPUUsage()
			if err != nil {
				time.Sleep(time.Second)
				return
			}

			memory, err := system.MemoryUsage()
			if err != nil {
				time.Sleep(time.Second)
				return
			}

			disk, err := system.DiskUsage("/")
			if err != nil {
				time.Sleep(time.Second)
				return
			}

			ui.AddCPUHistory(cpuUsage)
			ui.AddRAMHistory(memory.Usage)

			ui.Clear()
			ui.RenderDashboard(cpuUsage, memory, disk)

			time.Sleep(time.Second)
		}()
	}
}
