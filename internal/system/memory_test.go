// ©AngelaMos | 2026
// memory_test.go

package system

import (
	"testing"
)

func TestGetMemoryInfo(t *testing.T) {
	info, err := GetMemoryInfo()
	if err != nil {
		t.Fatalf("GetMemoryInfo failed: %v", err)
	}
	if info.TotalRAM == 0 {
		t.Error("TotalRAM should not be 0")
	}
	if info.UsedRAM > info.TotalRAM {
		t.Error("UsedRAM exceeds TotalRAM")
	}
}

func TestRAMPercent(t *testing.T) {
	info := MemoryInfo{TotalRAM: 1000, UsedRAM: 500}
	pct := info.RAMPercent()
	if pct < 49.9 || pct > 50.1 {
		t.Errorf("expected ~50%%, got %f", pct)
	}
}

func TestSwapPercent(t *testing.T) {
	info := MemoryInfo{TotalSwap: 1000, UsedSwap: 250}
	pct := info.SwapPercent()
	if pct < 24.9 || pct > 25.1 {
		t.Errorf("expected ~25%%, got %f", pct)
	}
}

func TestSwapPercentZeroTotal(t *testing.T) {
	info := MemoryInfo{TotalSwap: 0, UsedSwap: 0}
	pct := info.SwapPercent()
	if pct != 0 {
		t.Errorf("expected 0%%, got %f", pct)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0.0 MiB"},
		{1048576, "1.0 MiB"},
		{1073741824, "1.0 GiB"},
		{2147483648, "2.0 GiB"},
		{536870912, "512.0 MiB"},
	}
	for _, tc := range tests {
		got := FormatBytes(tc.input)
		if got != tc.expected {
			t.Errorf(
				"FormatBytes(%d) = %q, want %q",
				tc.input,
				got,
				tc.expected,
			)
		}
	}
}
