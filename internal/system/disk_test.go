// ©AngelaMos | 2026
// disk_test.go

package system

import (
	"testing"
)

func TestGetDiskUsage(t *testing.T) {
	info, err := GetDiskUsage("/")
	if err != nil {
		t.Fatalf("GetDiskUsage failed: %v", err)
	}
	if info.Total == 0 {
		t.Error("Total should not be 0")
	}
	if info.Used > info.Total {
		t.Error("Used exceeds Total")
	}
}

func TestDiskPercent(t *testing.T) {
	info := DiskInfo{Total: 1000, Used: 400, Free: 600}
	pct := info.Percent()
	if pct < 39.9 || pct > 40.1 {
		t.Errorf("expected ~40%%, got %f", pct)
	}
}

func TestDiskPercentZeroTotal(t *testing.T) {
	info := DiskInfo{Total: 0, Used: 0, Free: 0}
	pct := info.Percent()
	if pct != 0 {
		t.Errorf("expected 0%%, got %f", pct)
	}
}
