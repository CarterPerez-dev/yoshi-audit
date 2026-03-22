// ©AngelaMos | 2026
// audittab_test.go

package audittab

import (
	"testing"
)

func TestPruneRSSHistory_RemovesDeadPIDs(t *testing.T) {
	history := map[int][]uint64{
		100: {1, 2, 3},
		200: {4, 5, 6},
	}
	active := map[int]bool{100: true}
	pruneRSSHistory(history, active, 100)

	if _, exists := history[200]; exists {
		t.Error("dead PID 200 should have been removed")
	}
	if _, exists := history[100]; !exists {
		t.Error("active PID 100 should be kept")
	}
}

func TestPruneRSSHistory_CapsReadings(t *testing.T) {
	history := map[int][]uint64{
		100: {1, 2, 3, 4, 5},
	}
	active := map[int]bool{100: true}
	pruneRSSHistory(history, active, 3)

	if len(history[100]) != 3 {
		t.Errorf("expected 3 readings, got %d", len(history[100]))
	}
	if history[100][0] != 3 || history[100][2] != 5 {
		t.Errorf("should keep most recent readings, got %v", history[100])
	}
}

func TestPruneRSSHistory_UnderCap(t *testing.T) {
	history := map[int][]uint64{
		100: {1, 2},
	}
	active := map[int]bool{100: true}
	pruneRSSHistory(history, active, 100)

	if len(history[100]) != 2 {
		t.Errorf(
			"should not truncate when under cap, got %d",
			len(history[100]),
		)
	}
}
