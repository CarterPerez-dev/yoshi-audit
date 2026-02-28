// ©AngelaMos | 2026
// cpu_test.go

package system

import (
	"testing"
)

func TestParseCPULine(t *testing.T) {
	idle, total, err := parseCPULine("cpu  4705 356 584 3699 23 0 5 0 0 0")
	if err != nil {
		t.Fatalf("parseCPULine failed: %v", err)
	}
	if idle != 3699 {
		t.Errorf("expected idle=3699, got %d", idle)
	}
	if total != 9372 {
		t.Errorf("expected total=9372, got %d", total)
	}
}

func TestParseCPULineInvalid(t *testing.T) {
	_, _, err := parseCPULine("cpu  abc def")
	if err == nil {
		t.Error("expected error for invalid cpu line")
	}
}

func TestGetCPUUsage(t *testing.T) {
	usage, err := GetCPUUsage()
	if err != nil {
		t.Fatalf("GetCPUUsage failed: %v", err)
	}
	if usage < 0 || usage > 100 {
		t.Errorf("CPU usage out of range: %f", usage)
	}
}
