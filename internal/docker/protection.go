// ©AngelaMos | 2026
// protection.go

package docker

import (
	"path/filepath"
	"strings"
)

type ProtectionEngine struct {
	patterns  []string
	protected map[string]bool
}

func NewProtectionEngine(patterns []string) *ProtectionEngine {
	return &ProtectionEngine{
		patterns:  patterns,
		protected: make(map[string]bool),
	}
}

func (pe *ProtectionEngine) IsProtected(name string) bool {
	if pe.protected[name] {
		return true
	}
	return pe.MatchesPattern(name)
}

func (pe *ProtectionEngine) Protect(name string) {
	pe.protected[name] = true
}

func (pe *ProtectionEngine) Unprotect(name string) {
	delete(pe.protected, name)
}

func (pe *ProtectionEngine) MatchesPattern(name string) bool {
	for _, pattern := range pe.patterns {
		if matchGlob(pattern, name) {
			return true
		}
	}
	return false
}

func matchGlob(pattern, name string) bool {
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") &&
		!strings.Contains(pattern[1:len(pattern)-1], "*") {
		inner := pattern[1 : len(pattern)-1]
		return strings.Contains(strings.ToLower(name), strings.ToLower(inner))
	}

	matched, err := filepath.Match(
		strings.ToLower(pattern),
		strings.ToLower(name),
	)
	if err != nil {
		return false
	}
	return matched
}
