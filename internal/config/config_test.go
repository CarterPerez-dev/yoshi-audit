// ©AngelaMos | 2026
// config_test.go

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	if cfg.RefreshInterval != 2 {
		t.Errorf("expected RefreshInterval=2, got %d", cfg.RefreshInterval)
	}

	if len(cfg.ProtectionPatterns) != 3 {
		t.Fatalf("expected 3 protection patterns, got %d", len(cfg.ProtectionPatterns))
	}

	expected := []string{"*certgames*", "*argos*", "*mongo*"}
	for i, p := range expected {
		if cfg.ProtectionPatterns[i] != p {
			t.Errorf("protection pattern %d: expected %q, got %q", i, p, cfg.ProtectionPatterns[i])
		}
	}

	if len(cfg.PrunePresets) != 2 {
		t.Fatalf("expected 2 prune presets, got %d", len(cfg.PrunePresets))
	}

	if cfg.PrunePresets[0].Name != "Dangling Only" {
		t.Errorf("expected first preset name 'Dangling Only', got %q", cfg.PrunePresets[0].Name)
	}

	if cfg.PrunePresets[1].Name != "Build Cache" {
		t.Errorf("expected second preset name 'Build Cache', got %q", cfg.PrunePresets[1].Name)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test-config.yaml")

	cfg := Config{
		RefreshInterval:    5,
		ProtectionPatterns: []string{"*custom*"},
		PrunePresets: []PrunePreset{
			{Name: "Test Preset", Patterns: []string{"test-pattern"}},
		},
		BaselinePath: "/tmp/test-baseline.json",
	}

	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.RefreshInterval != 5 {
		t.Errorf("expected RefreshInterval=5, got %d", loaded.RefreshInterval)
	}

	if len(loaded.ProtectionPatterns) != 1 || loaded.ProtectionPatterns[0] != "*custom*" {
		t.Errorf("unexpected ProtectionPatterns: %v", loaded.ProtectionPatterns)
	}

	if len(loaded.PrunePresets) != 1 || loaded.PrunePresets[0].Name != "Test Preset" {
		t.Errorf("unexpected PrunePresets: %v", loaded.PrunePresets)
	}

	if loaded.BaselinePath != "/tmp/test-baseline.json" {
		t.Errorf("expected BaselinePath '/tmp/test-baseline.json', got %q", loaded.BaselinePath)
	}
}

func TestLoadNonExistent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.yaml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error for nonexistent file, got: %v", err)
	}

	def := Default()

	if cfg.RefreshInterval != def.RefreshInterval {
		t.Errorf("expected default RefreshInterval=%d, got %d", def.RefreshInterval, cfg.RefreshInterval)
	}

	if len(cfg.ProtectionPatterns) != len(def.ProtectionPatterns) {
		t.Errorf("expected %d protection patterns, got %d", len(def.ProtectionPatterns), len(cfg.ProtectionPatterns))
	}
}
