package system

import (
	"bufio"
	"context"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
)

type DockerStats struct {
	Available         bool
	ServerVersion     string
	CPUs              int
	TotalMemory       uint64
	RunningContainers int
	TotalContainers   int
}

type ContainerStats struct {
	ID         string
	Name       string
	Status     string
	CPUUsage   string
	MemUsage   string
	CPUPercent float64
	MemBytes   uint64
}

func DockerInfo() (DockerStats, []ContainerStats, error) {
	info := DockerStats{
		Available: false,
	}

	if _, err := exec.LookPath("docker"); err != nil {
		return info, nil, nil
	}

	info.Available = true

	if out, err := exec.Command("docker", "version", "--format", "{{.Server.Version}}").Output(); err == nil {
		info.ServerVersion = strings.TrimSpace(string(out))
	}

	if out, err := exec.Command("docker", "info", "--format", "{{.NCPU}};{{.MemTotal}};{{.Containers}};{{.ContainersRunning}}").Output(); err == nil {
		parts := strings.Split(strings.TrimSpace(string(out)), ";")

		if len(parts) >= 4 {
			info.CPUs, _ = strconv.Atoi(parts[0])

			info.TotalMemory, _ = strconv.ParseUint(parts[1], 10, 64)

			info.TotalContainers, _ = strconv.Atoi(parts[2])

			info.RunningContainers, _ = strconv.Atoi(parts[3])

			if info.TotalMemory == 0 {
				m, _ := mem.VirtualMemory()
				info.TotalMemory = m.Total
			}
		}
	}

	containers, err := fetchContainers()

	return info, containers, err
}

func fetchContainers() ([]ContainerStats, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		3500*time.Millisecond,
	)

	defer cancel()

	psCmd := exec.CommandContext(
		ctx,
		"docker",
		"ps",
		"--format",
		"{{.Names}};{{.Status}};{{.ID}}",
	)

	out, err := psCmd.Output()

	if err != nil {
		return nil, err
	}

	var list []struct {
		Name   string
		Status string
		ID     string
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ";", 3)

		if len(parts) != 3 {
			continue
		}

		list = append(list, struct {
			Name   string
			Status string
			ID     string
		}{
			Name:   parts[0],
			Status: parts[1],
			ID:     parts[2],
		})
	}

	statsCmd := exec.CommandContext(
		ctx,
		"docker",
		"stats",
		"--no-stream",
		"--format",
		"{{.Name}};{{.CPUPerc}};{{.MemUsage}}",
	)

	statsOut, _ := statsCmd.Output()

	stats := map[string][2]string{}

	ss := bufio.NewScanner(
		strings.NewReader(string(statsOut)),
	)

	for ss.Scan() {
		line := strings.TrimSpace(ss.Text())

		if line == "" {
			continue
		}

		p := strings.SplitN(line, ";", 3)

		if len(p) < 3 {
			continue
		}

		stats[p[0]] = [2]string{
			p[1],
			p[2],
		}
	}

	var result []ContainerStats

	for _, item := range list {

		cpu := "N/A"
		mem := "N/A"

		if v, ok := stats[item.Name]; ok {
			cpu = v[0]
			mem = v[1]
		}

		id := item.ID

		if len(id) > 12 {
			id = id[:12]
		}

		result = append(result, ContainerStats{
			ID:         id,
			Name:       item.Name,
			Status:     item.Status,
			CPUUsage:   cpu,
			MemUsage:   mem,
			CPUPercent: parseCPUPercent(cpu),
			MemBytes:   parseMemUsageBytes(mem),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CPUPercent > result[j].CPUPercent
	})

	return result, nil
}

func parseCPUPercent(s string) float64 {
	s = strings.TrimSpace(
		strings.TrimSuffix(s, "%"),
	)

	v, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0
	}

	return v
}

func parseMemUsageBytes(s string) uint64 {
	used := strings.TrimSpace(
		strings.SplitN(s, "/", 2)[0],
	)

	if used == "" {
		return 0
	}

	var number string
	var unit string

	for i, r := range used {
		if !(r >= '0' && r <= '9' || r == '.') {
			number = used[:i]
			unit = strings.TrimSpace(used[i:])
			break
		}
	}

	if number == "" {
		number = used
	}

	value, err := strconv.ParseFloat(number, 64)

	if err != nil {
		return 0
	}

	multiplier := float64(1)

	switch strings.ToUpper(unit) {
	case "KB", "KIB":
		multiplier = 1024
	case "MB", "MIB":
		multiplier = 1024 * 1024
	case "GB", "GIB":
		multiplier = 1024 * 1024 * 1024
	case "TB", "TIB":
		multiplier = 1024 * 1024 * 1024 * 1024
	}

	return uint64(value * multiplier)
}
