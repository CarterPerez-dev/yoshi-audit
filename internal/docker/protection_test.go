// ©AngelaMos | 2026
// protection_test.go

package docker

import (
	"testing"
)

func TestMatchesPattern(t *testing.T) {
	pe := NewProtectionEngine([]string{"*certgames*", "*argos*", "*mongo*"})

	shouldMatch := []struct {
		name string
		desc string
	}{
		{"certgames-prod_redis_data", "should match certgames"},
		{"certgamesdb-argos_backend_backups", "should match argos"},
		{"oneisnun__mongo_data", "should match mongo"},
	}
	for _, tt := range shouldMatch {
		if !pe.MatchesPattern(tt.name) {
			t.Errorf("%s: %q did not match", tt.desc, tt.name)
		}
	}

	shouldNotMatch := []struct {
		name string
		desc string
	}{
		{"redis:7-alpine", "should not match redis"},
		{"nginx:latest", "should not match nginx"},
		{"python:3.12-slim", "should not match python"},
	}
	for _, tt := range shouldNotMatch {
		if pe.MatchesPattern(tt.name) {
			t.Errorf("%s: %q should not have matched", tt.desc, tt.name)
		}
	}
}

func TestProtectUnprotect(t *testing.T) {
	pe := NewProtectionEngine(nil)
	pe.Protect("my-important-volume")
	if !pe.IsProtected("my-important-volume") {
		t.Error("should be protected after Protect()")
	}
	pe.Unprotect("my-important-volume")
	if pe.IsProtected("my-important-volume") {
		t.Error("should not be protected after Unprotect()")
	}
}

func TestIsProtected_PatternMatch(t *testing.T) {
	pe := NewProtectionEngine([]string{"*mongo*"})
	if !pe.IsProtected("my-mongo-volume") {
		t.Error("pattern match should make it protected")
	}
}

func TestIsProtected_Explicit(t *testing.T) {
	pe := NewProtectionEngine(nil)
	if pe.IsProtected("random-vol") {
		t.Error("should not be protected by default")
	}
	pe.Protect("random-vol")
	if !pe.IsProtected("random-vol") {
		t.Error("should be protected after explicit Protect()")
	}
}

func TestMatchesPattern_EmptyPatterns(t *testing.T) {
	pe := NewProtectionEngine(nil)
	if pe.MatchesPattern("anything") {
		t.Error("no patterns should match nothing")
	}
}

func TestMatchesPattern_ExactGlob(t *testing.T) {
	pe := NewProtectionEngine([]string{"redis:*"})
	if !pe.MatchesPattern("redis:7-alpine") {
		t.Error("redis:* should match redis:7-alpine")
	}
	if pe.MatchesPattern("nginx:latest") {
		t.Error("redis:* should not match nginx:latest")
	}
}
