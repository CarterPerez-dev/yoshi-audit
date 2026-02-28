// ©AngelaMos | 2026
// gpu.go

package system

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type GPUInfo struct {
	Utilization float64
	UsedVRAM    uint64
	TotalVRAM   uint64
}

func (g GPUInfo) VRAMPercent() float64 {
	if g.TotalVRAM == 0 {
		return 0
	}
	return float64(g.UsedVRAM) / float64(g.TotalVRAM) * 100
}

type GPUProcess struct {
	PID      int
	UsedVRAM uint64
}

func GetGPUInfo() (GPUInfo, error) {
	out, err := exec.Command(
		"nvidia-smi",
		"--query-gpu=utilization.gpu,memory.used,memory.total",
		"--format=csv,noheader,nounits",
	).Output()
	if err != nil {
		return GPUInfo{}, fmt.Errorf("nvidia-smi failed: %w", err)
	}

	line := strings.TrimSpace(string(out))
	parts := strings.Split(line, ",")
	if len(parts) < 3 {
		return GPUInfo{}, fmt.Errorf("unexpected nvidia-smi output: %q", line)
	}

	util, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return GPUInfo{}, fmt.Errorf("failed to parse utilization: %w", err)
	}

	usedMiB, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return GPUInfo{}, fmt.Errorf("failed to parse used memory: %w", err)
	}

	totalMiB, err := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)
	if err != nil {
		return GPUInfo{}, fmt.Errorf("failed to parse total memory: %w", err)
	}

	return GPUInfo{
		Utilization: util,
		UsedVRAM:    usedMiB * 1024 * 1024,
		TotalVRAM:   totalMiB * 1024 * 1024,
	}, nil
}

func GetGPUProcesses() ([]GPUProcess, error) {
	out, err := exec.Command(
		"nvidia-smi",
		"--query-compute-apps=pid,used_memory",
		"--format=csv,noheader,nounits",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("nvidia-smi failed: %w", err)
	}

	var procs []GPUProcess
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "[") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		pid, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		memMiB, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			continue
		}

		procs = append(procs, GPUProcess{
			PID:      pid,
			UsedVRAM: memMiB * 1024 * 1024,
		})
	}

	return procs, nil
}
