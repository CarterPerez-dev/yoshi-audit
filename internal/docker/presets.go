// ©AngelaMos | 2026
// presets.go

package docker

import (
	"github.com/CarterPerez-dev/yoshi-audit/internal/config"
)

func ApplyPreset(
	preset config.PrunePreset,
	images []ImageInfo,
	volumes []VolumeInfo,
	containers []ContainerInfo,
	protection *ProtectionEngine,
) (imageIDs, volumeNames, containerIDs []string) {
	hasDangling := false
	hasStopped := false
	for _, p := range preset.Patterns {
		switch p {
		case "dangling":
			hasDangling = true
		case "stopped":
			hasStopped = true
		}
	}

	if hasDangling {
		for _, img := range images {
			if img.Dangling && !protection.IsProtected(img.Repository) &&
				!protection.IsProtected(img.ID) {
				imageIDs = append(imageIDs, img.ID)
			}
		}
	}

	if hasStopped {
		for _, ctr := range containers {
			if !ctr.Running && !protection.IsProtected(ctr.Name) &&
				!protection.IsProtected(ctr.Image) {
				containerIDs = append(containerIDs, ctr.ID)
			}
		}
	}

	for _, p := range preset.Patterns {
		if p == "dangling" || p == "buildcache" || p == "stopped" {
			continue
		}
		pe := NewProtectionEngine([]string{p})
		for _, img := range images {
			if !img.Dangling && pe.MatchesPattern(img.Repository) &&
				!protection.IsProtected(img.Repository) &&
				!protection.IsProtected(img.ID) {
				imageIDs = append(imageIDs, img.ID)
			}
		}
		for _, vol := range volumes {
			if pe.MatchesPattern(vol.Name) &&
				!protection.IsProtected(vol.Name) {
				volumeNames = append(volumeNames, vol.Name)
			}
		}
	}

	return imageIDs, volumeNames, containerIDs
}

func EstimateReclaimable(
	images []ImageInfo,
	selectedImageIDs []string,
	volumes []VolumeInfo,
	selectedVolumeNames []string,
) int64 {
	imageSet := make(map[string]bool, len(selectedImageIDs))
	for _, id := range selectedImageIDs {
		imageSet[id] = true
	}

	volumeSet := make(map[string]bool, len(selectedVolumeNames))
	for _, name := range selectedVolumeNames {
		volumeSet[name] = true
	}

	var total int64
	for _, img := range images {
		if imageSet[img.ID] {
			total += img.Size
		}
	}
	for _, vol := range volumes {
		if volumeSet[vol.Name] {
			total += vol.Size
		}
	}
	return total
}
