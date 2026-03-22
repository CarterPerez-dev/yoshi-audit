// ©AngelaMos | 2026
// baseline.go

package audit

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/CarterPerez-dev/yoshi-audit/internal/system"
)

type Baseline struct {
	Processes map[string]bool `json:"processes"`
}

func NewBaseline() *Baseline {
	return &Baseline{
		Processes: make(map[string]bool),
	}
}

func (b *Baseline) Load(path string) error {
	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, b)
}

func (b *Baseline) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}

	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

func (b *Baseline) IsKnown(name string) bool {
	return b.Processes[name]
}

func (b *Baseline) Add(name string) {
	b.Processes[name] = true
}

func (b *Baseline) Remove(name string) {
	delete(b.Processes, name)
}

func (b *Baseline) BuildFromCurrent(procs []system.ProcessInfo) {
	for _, p := range procs {
		b.Processes[p.Name] = true
	}
}
