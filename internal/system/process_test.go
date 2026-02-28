// ©AngelaMos | 2026
// process_test.go

package system

import (
	"testing"
)

func TestGetProcesses(t *testing.T) {
	procs, err := GetProcesses()
	if err != nil {
		t.Fatalf("GetProcesses failed: %v", err)
	}
	if len(procs) == 0 {
		t.Error("expected at least one process")
	}

	foundNamed := false
	for _, p := range procs {
		if p.PID > 0 && p.Name != "" {
			foundNamed = true
			break
		}
	}
	if !foundNamed {
		t.Error("expected at least one named process")
	}
}

func TestProcessState(t *testing.T) {
	procs, err := GetProcesses()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range procs {
		if p.State == "" {
			t.Errorf("PID %d has empty state", p.PID)
		}
	}
}

func TestGetZombies(t *testing.T) {
	procs := []ProcessInfo{
		{PID: 1, State: "S"},
		{PID: 2, State: "Z"},
		{PID: 3, State: "R"},
		{PID: 4, State: "Z"},
	}
	zombies := GetZombies(procs)
	if len(zombies) != 2 {
		t.Errorf("expected 2 zombies, got %d", len(zombies))
	}
}

func TestGetOrphans(t *testing.T) {
	procs := []ProcessInfo{
		{PID: 10, PPID: 1, State: "S"},
		{PID: 20, PPID: 1, State: "Z"},
		{PID: 30, PPID: 5, State: "S"},
		{PID: 40, PPID: 1, State: "R"},
	}
	orphans := GetOrphans(procs)
	if len(orphans) != 2 {
		t.Errorf("expected 2 orphans, got %d", len(orphans))
	}
}

func TestGetKernelThreads(t *testing.T) {
	procs := []ProcessInfo{
		{PID: 10, PPID: 2, State: "S"},
		{PID: 20, PPID: 1, State: "S"},
		{PID: 30, PPID: 2, State: "S"},
	}
	kthreads := GetKernelThreads(procs)
	if len(kthreads) != 2 {
		t.Errorf("expected 2 kernel threads, got %d", len(kthreads))
	}
}

func TestFormatProcessName(t *testing.T) {
	if got := FormatProcessName("bash", 10); got != "bash      " {
		t.Errorf("expected padded name, got %q", got)
	}
	if got := FormatProcessName("very-long-process-name", 10); got != "very-lo..." {
		t.Errorf("expected truncated name, got %q", got)
	}
	if got := FormatProcessName("exact_fit!", 10); got != "exact_fit!" {
		t.Errorf("expected exact fit, got %q", got)
	}
}

func TestSortByMemory(t *testing.T) {
	procs := []ProcessInfo{
		{PID: 1, RSS: 100},
		{PID: 2, RSS: 500},
		{PID: 3, RSS: 300},
	}
	SortByMemory(procs)
	if procs[0].RSS != 500 || procs[1].RSS != 300 || procs[2].RSS != 100 {
		t.Errorf("not sorted descending by RSS: %v", procs)
	}
}

func TestTopN(t *testing.T) {
	procs := []ProcessInfo{
		{PID: 1}, {PID: 2}, {PID: 3}, {PID: 4}, {PID: 5},
	}
	top := TopN(procs, 3)
	if len(top) != 3 {
		t.Errorf("expected 3, got %d", len(top))
	}
	top2 := TopN(procs, 10)
	if len(top2) != 5 {
		t.Errorf("expected 5, got %d", len(top2))
	}
}
