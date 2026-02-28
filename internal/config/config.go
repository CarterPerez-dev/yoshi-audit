// ©AngelaMos | 2026
// config.go

package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type PrunePreset struct {
	Name     string   `yaml:"name"`
	Patterns []string `yaml:"patterns"`
}

type Config struct {
	RefreshInterval    int           `yaml:"refresh_interval"`
	ProtectionPatterns []string      `yaml:"protection_patterns"`
	PrunePresets       []PrunePreset `yaml:"prune_presets"`
	BaselinePath       string        `yaml:"baseline_path"`
}

func Default() Config {
	return Config{
		RefreshInterval:    2,
		ProtectionPatterns: []string{"*certgames*", "*argos*", "*mongo*"},
		PrunePresets: []PrunePreset{
			{Name: "Dangling Only", Patterns: []string{"dangling"}},
			{Name: "Build Cache", Patterns: []string{"buildcache"}},
		},
		BaselinePath: DefaultBaselinePath(),
	}
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "yoshi-audit")
}

func DefaultPath() string {
	return filepath.Join(configDir(), "config.yaml")
}

func DefaultBaselinePath() string {
	return filepath.Join(configDir(), "baseline.json")
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Default(), nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
