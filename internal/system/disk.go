// ©AngelaMos | 2026
// disk.go

package system

import (
	"syscall"
)

type DiskInfo struct {
	Total uint64
	Used  uint64
	Free  uint64
}

func (d DiskInfo) Percent() float64 {
	if d.Total == 0 {
		return 0
	}
	return float64(d.Used) / float64(d.Total) * 100
}

func GetDiskUsage(path string) (DiskInfo, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskInfo{}, err
	}

	total := stat.Blocks * uint64(stat.Bsize) //nolint:gosec
	free := stat.Bavail * uint64(stat.Bsize)  //nolint:gosec
	used := total - free

	return DiskInfo{
		Total: total,
		Used:  used,
		Free:  free,
	}, nil
}
