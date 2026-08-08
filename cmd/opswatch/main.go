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

	go listenKeys()

	for {
		cpuUsage, err := system.CPUUsage()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		memory, err := system.MemoryUsage()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		disk, err := system.DiskUsage("/")
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		ui.AddCPUHistory(cpuUsage)
		ui.AddRAMHistory(memory.Usage)

		ui.Clear()
		ui.RenderDashboard(cpuUsage, memory, disk)

		time.Sleep(time.Second)
	}
}