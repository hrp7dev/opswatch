package system

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"time"
)

func CPUUsage() (float64, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) == 0 {
		return 0, nil
	}
	return percentages[0], nil
}
func CPUCores() int {
	cores, err := cpu.Counts(true)
	if err != nil {
		return 0
	}
	return cores
}
