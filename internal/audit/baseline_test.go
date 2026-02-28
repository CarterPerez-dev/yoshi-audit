// ©AngelaMos | 2026
// baseline_test.go

package audit

import (
	"path/filepath"
	"testing"

	"github.com/CarterPerez-dev/yoshi-audit/internal/system"
)

func TestBaselineSaveLoad(t *testing.T) {
	b := NewBaseline()
	b.Add("my_process")

	path := filepath.Join(t.TempDir(), "baseline.json")
	if err := b.Save(path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	b2 := NewBaseline()
	if err := b2.Load(path); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if !b2.IsKnown("my_process") {
		t.Error("should be known after load")
	}
}

func TestBaselineBuildFromCurrent(t *testing.T) {
	b := NewBaseline()
	procs := []system.ProcessInfo{
		{Name: "proc1"}, {Name: "proc2"}, {Name: "proc1"},
	}
	b.BuildFromCurrent(procs)
	if !b.IsKnown("proc1") || !b.IsKnown("proc2") {
		t.Error("should know both")
	}
}

func TestBaselineRemove(t *testing.T) {
	b := NewBaseline()
	b.Add("foo")
	b.Add("bar")
	b.Remove("foo")
	if b.IsKnown("foo") {
		t.Error("foo should be removed")
	}
	if !b.IsKnown("bar") {
		t.Error("bar should still be known")
	}
}

func TestBaselineLoadMissing(t *testing.T) {
	b := NewBaseline()
	err := b.Load("/nonexistent/path/baseline.json")
	if err != nil {
		t.Error("loading missing file should not error")
	}
}

func TestBaselineNewIsEmpty(t *testing.T) {
	b := NewBaseline()
	if b.IsKnown("anything") {
		t.Error("new baseline should not know anything")
	}
}

func TestBaselineSaveCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "deep", "baseline.json")
	b := NewBaseline()
	b.Add("test")
	if err := b.Save(path); err != nil {
		t.Fatalf("save should create directories: %v", err)
	}

	b2 := NewBaseline()
	if err := b2.Load(path); err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if !b2.IsKnown("test") {
		t.Error("should be known after save/load through nested dirs")
	}
}
