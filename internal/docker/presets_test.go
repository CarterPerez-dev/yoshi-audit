// ©AngelaMos | 2026
// presets_test.go

package docker

import (
	"testing"

	"github.com/CarterPerez-dev/yoshi-audit/internal/config"
)

func TestApplyPreset_Dangling(t *testing.T) {
	preset := config.PrunePreset{
		Name:     "Dangling Only",
		Patterns: []string{"dangling"},
	}
	images := []ImageInfo{
		{ID: "img1", Repository: "<none>", Tag: "<none>", Dangling: true},
		{ID: "img2", Repository: "nginx", Tag: "latest", Dangling: false},
	}
	pe := NewProtectionEngine(nil)

	imageIDs, _, _ := ApplyPreset(preset, images, nil, nil, pe)
	if len(imageIDs) != 1 || imageIDs[0] != "img1" {
		t.Errorf("expected only dangling img1, got %v", imageIDs)
	}
}

func TestApplyPreset_Stopped(t *testing.T) {
	preset := config.PrunePreset{Name: "Stopped", Patterns: []string{"stopped"}}
	containers := []ContainerInfo{
		{ID: "c1", Name: "old-worker", Image: "worker:v1", Running: false},
		{ID: "c2", Name: "web", Image: "nginx:latest", Running: true},
		{ID: "c3", Name: "another-stopped", Image: "redis:7", Running: false},
	}
	pe := NewProtectionEngine(nil)

	_, _, containerIDs := ApplyPreset(preset, nil, nil, containers, pe)
	if len(containerIDs) != 2 {
		t.Errorf(
			"expected 2 stopped containers, got %d: %v",
			len(containerIDs),
			containerIDs,
		)
	}
}

func TestApplyPreset_Stopped_ProtectedSkipped(t *testing.T) {
	preset := config.PrunePreset{Name: "Stopped", Patterns: []string{"stopped"}}
	containers := []ContainerInfo{
		{
			ID:      "c1",
			Name:    "certgames-api",
			Image:   "certgames:latest",
			Running: false,
		},
		{ID: "c2", Name: "old-worker", Image: "worker:v1", Running: false},
	}
	pe := NewProtectionEngine([]string{"*certgames*"})

	_, _, containerIDs := ApplyPreset(preset, nil, nil, containers, pe)
	if len(containerIDs) != 1 || containerIDs[0] != "c2" {
		t.Errorf("expected only old-worker, got %v", containerIDs)
	}
}
