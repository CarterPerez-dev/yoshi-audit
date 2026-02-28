// ©AngelaMos | 2026
// gpu_test.go

package system

import (
	"testing"
)

func TestGetGPUInfo(t *testing.T) {
	info, err := GetGPUInfo()
	if err != nil {
		t.Skipf("nvidia-smi not available: %v", err)
	}
	if info.TotalVRAM == 0 {
		t.Error("TotalVRAM should not be 0")
	}
	if info.Utilization < 0 || info.Utilization > 100 {
		t.Errorf("out of range: %f", info.Utilization)
	}
}

func TestGetGPUProcesses(t *testing.T) {
	procs, err := GetGPUProcesses()
	if err != nil {
		t.Skipf("nvidia-smi not available: %v", err)
	}
	_ = procs
}

func TestVRAMPercent(t *testing.T) {
	info := GPUInfo{TotalVRAM: 12000, UsedVRAM: 6000}
	pct := info.VRAMPercent()
	if pct < 49.9 || pct > 50.1 {
		t.Errorf("expected ~50%%, got %f", pct)
	}
}

func TestVRAMPercentZeroTotal(t *testing.T) {
	info := GPUInfo{TotalVRAM: 0, UsedVRAM: 0}
	pct := info.VRAMPercent()
	if pct != 0 {
		t.Errorf("expected 0%%, got %f", pct)
	}
}
