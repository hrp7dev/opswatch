package system

import "github.com/shirou/gopsutil/v4/mem"

type MemoryStats struct {
	Total uint64
	Used  uint64
	Free  uint64
	Usage float64
}

func MemoryUsage() (MemoryStats, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return MemoryStats{}, err
	}
	return MemoryStats{
		Total: memory.Total,
		Used:  memory.Used,
		Free:  memory.Available,
		Usage: memory.UsedPercent,
	}, nil
}