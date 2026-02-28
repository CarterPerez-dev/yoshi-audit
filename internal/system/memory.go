// ©AngelaMos | 2026
// memory.go

package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MemoryInfo struct {
	TotalRAM     uint64
	UsedRAM      uint64
	AvailableRAM uint64
	TotalSwap    uint64
	UsedSwap     uint64
}

func (m MemoryInfo) RAMPercent() float64 {
	if m.TotalRAM == 0 {
		return 0
	}
	return float64(m.UsedRAM) / float64(m.TotalRAM) * 100
}

func (m MemoryInfo) SwapPercent() float64 {
	if m.TotalSwap == 0 {
		return 0
	}
	return float64(m.UsedSwap) / float64(m.TotalSwap) * 100
}

func GetMemoryInfo() (MemoryInfo, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}, err
	}
	defer f.Close()

	values := make(map[string]uint64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		valStr := strings.TrimSpace(parts[1])
		valStr = strings.TrimSuffix(valStr, " kB")
		valStr = strings.TrimSpace(valStr)

		val, parseErr := strconv.ParseUint(valStr, 10, 64)
		if parseErr != nil {
			continue
		}

		values[key] = val * 1024
	}

	if err := scanner.Err(); err != nil {
		return MemoryInfo{}, err
	}

	memTotal, ok := values["MemTotal"]
	if !ok {
		return MemoryInfo{}, fmt.Errorf("MemTotal not found in /proc/meminfo")
	}

	memAvailable, ok := values["MemAvailable"]
	if !ok {
		return MemoryInfo{}, fmt.Errorf("MemAvailable not found in /proc/meminfo")
	}

	swapTotal := values["SwapTotal"]
	swapFree := values["SwapFree"]

	return MemoryInfo{
		TotalRAM:     memTotal,
		UsedRAM:      memTotal - memAvailable,
		AvailableRAM: memAvailable,
		TotalSwap:    swapTotal,
		UsedSwap:     swapTotal - swapFree,
	}, nil
}

func FormatBytes(b uint64) string {
	const gib = 1024 * 1024 * 1024
	const mib = 1024 * 1024

	if b >= gib {
		return fmt.Sprintf("%.1f GiB", float64(b)/float64(gib))
	}
	return fmt.Sprintf("%.1f MiB", float64(b)/float64(mib))
}
