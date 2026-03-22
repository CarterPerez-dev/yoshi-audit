// ©AngelaMos | 2026
// scanner_test.go

package audit

import (
	"testing"

	"github.com/CarterPerez-dev/yoshi-audit/internal/system"
)

func TestScanZombies(t *testing.T) {
	s := NewScanner(nil)
	procs := []system.ProcessInfo{
		{PID: 1, Name: "systemd", State: "S", PPID: 0},
		{PID: 100, Name: "defunct", State: "Z", PPID: 1},
	}
	findings := s.ScanZombies(procs)
	if len(findings) != 1 {
		t.Errorf("expected 1 zombie, got %d", len(findings))
	}
	if findings[0].Type != FindingZombie {
		t.Error("wrong type")
	}
}

func TestScanOrphans(t *testing.T) {
	s := NewScanner(nil)
	procs := []system.ProcessInfo{
		{PID: 1, Name: "systemd", State: "S", PPID: 0},
		{PID: 100, Name: "mystery_proc", State: "S", PPID: 1},
		{PID: 200, Name: "dockerd", State: "S", PPID: 1},
	}
	findings := s.ScanOrphans(procs)
	foundMystery := false
	for _, f := range findings {
		if f.Name == "mystery_proc" {
			foundMystery = true
		}
		if f.Name == "dockerd" {
			t.Error("dockerd should not be flagged")
		}
	}
	if !foundMystery {
		t.Error("mystery_proc should be flagged as orphan")
	}
}

func TestScanDaemons(t *testing.T) {
	s := NewScanner(nil)
	procs := []system.ProcessInfo{
		{
			PID:   100,
			Name:  "some_weird_daemon",
			State: "S",
			PPID:  1,
			RSS:   50 * 1024 * 1024,
		},
		{PID: 200, Name: "sshd", State: "S", PPID: 1, RSS: 10 * 1024 * 1024},
	}
	findings := s.ScanDaemons(procs)
	for _, f := range findings {
		if f.Name == "sshd" {
			t.Error("sshd should not be flagged")
		}
	}
}

func TestScanMemoryLeaks(t *testing.T) {
	s := NewScanner(nil)
	procs := []system.ProcessInfo{
		{PID: 100, Name: "leaky", State: "S"},
	}
	history := map[int][]uint64{
		100: {100 * 1024 * 1024, 150 * 1024 * 1024, 250 * 1024 * 1024},
	}
	findings := s.ScanMemoryLeaks(procs, history)
	if len(findings) != 1 {
		t.Errorf("expected 1 leak, got %d", len(findings))
	}
}

func TestScanZombiesEmpty(t *testing.T) {
	s := NewScanner(nil)
	findings := s.ScanZombies(nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for nil procs, got %d", len(findings))
	}
}

func TestScanGPUShadows(t *testing.T) {
	s := NewScanner(nil)
	procs := []system.ProcessInfo{
		{PID: 100, Name: "app1", State: "S"},
	}
	gpuProcs := []system.GPUProcess{
		{PID: 100, UsedVRAM: 512 * 1024 * 1024},
		{PID: 999, UsedVRAM: 256 * 1024 * 1024},
	}
	findings := s.ScanGPUShadows(procs, gpuProcs)
	if len(findings) != 1 {
		t.Errorf("expected 1 gpu shadow, got %d", len(findings))
	}
	if len(findings) > 0 && findings[0].PID != 999 {
		t.Errorf("expected shadow PID 999, got %d", findings[0].PID)
	}
}

func TestScanKernelThreads(t *testing.T) {
	s := NewScanner(nil)
	procs := []system.ProcessInfo{
		{PID: 10, Name: "kworker/0:0", State: "S", PPID: 2},
		{PID: 11, Name: "kworker/1:0", State: "S", PPID: 2},
		{PID: 12, Name: "rcu_sched", State: "S", PPID: 2},
		{PID: 100, Name: "bash", State: "S", PPID: 1},
	}
	findings := s.ScanKernelThreads(procs)
	if len(findings) != 1 {
		t.Errorf("expected 1 informational finding, got %d", len(findings))
	}
	if len(findings) > 0 && findings[0].Severity != SeverityOK {
		t.Errorf("expected OK severity, got %s", findings[0].Severity)
	}
}

func TestScanMemoryLeaksNoLeak(t *testing.T) {
	s := NewScanner(nil)
	procs := []system.ProcessInfo{
		{PID: 100, Name: "stable", State: "S"},
	}
	history := map[int][]uint64{
		100: {50 * 1024 * 1024, 52 * 1024 * 1024, 51 * 1024 * 1024},
	}
	findings := s.ScanMemoryLeaks(procs, history)
	if len(findings) != 0 {
		t.Errorf("expected 0 leaks for stable process, got %d", len(findings))
	}
}

func TestScanAllCombined(t *testing.T) {
	s := NewScanner(nil)
	procs := []system.ProcessInfo{
		{PID: 1, Name: "systemd", State: "S", PPID: 0},
		{PID: 100, Name: "defunct", State: "Z", PPID: 1},
		{PID: 200, Name: "mystery", State: "S", PPID: 1},
	}
	gpuProcs := []system.GPUProcess{
		{PID: 999, UsedVRAM: 256 * 1024 * 1024},
	}
	findings := s.ScanAll(procs, gpuProcs)
	if len(findings) == 0 {
		t.Error("expected some findings from ScanAll")
	}
}
