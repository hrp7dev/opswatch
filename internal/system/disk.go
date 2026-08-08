package system

import (
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskStats struct {
	Total uint64
	Used  uint64
	Free  uint64
	Usage float64
}

func DiskUsage(path string) (DiskStats, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return DiskStats{}, err
	}
	return DiskStats{
		Total: usage.Total,
		Used: usage.Used,
		Free: usage.Free,
		Usage: usage.UsedPercent,
	}, nil
}