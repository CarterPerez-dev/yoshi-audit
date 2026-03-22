// ©AngelaMos | 2026
// categorize.go

package docker

import (
	"strings"
)

type SafetyCategory int

const (
	CategorySafe         SafetyCategory = 1
	CategoryProbablySafe SafetyCategory = 2
	CategoryCheckFirst   SafetyCategory = 3
	CategoryDoNotTouch   SafetyCategory = 4
)

func CategorizeImage(
	img ImageInfo,
	protectionPatterns []string,
) SafetyCategory {
	if img.Dangling && img.Containers == 0 {
		return CategorySafe
	}

	if img.Containers > 0 {
		return CategoryDoNotTouch
	}

	pe := NewProtectionEngine(protectionPatterns)
	fullName := img.Repository
	if img.Tag != noneTag {
		fullName = img.Repository + ":" + img.Tag
	}
	if pe.MatchesPattern(fullName) || pe.MatchesPattern(img.Repository) {
		return CategoryDoNotTouch
	}

	return CategoryProbablySafe
}

func CategorizeVolume(
	vol VolumeInfo,
	protectionPatterns []string,
) SafetyCategory {
	pe := NewProtectionEngine(protectionPatterns)
	if pe.MatchesPattern(vol.Name) {
		return CategoryDoNotTouch
	}

	if strings.Contains(strings.ToLower(vol.Name), "backup") {
		return CategoryCheckFirst
	}

	if vol.Links > 0 {
		return CategoryDoNotTouch
	}

	if looksLikeHash(vol.Name) && vol.Size <= 0 {
		return CategorySafe
	}

	if vol.Name != "" && vol.Size > 0 {
		return CategoryProbablySafe
	}

	return CategorySafe
}

func CategorizeContainer(
	ctr ContainerInfo,
	protectionPatterns []string,
) SafetyCategory {
	if ctr.Running {
		return CategoryDoNotTouch
	}

	pe := NewProtectionEngine(protectionPatterns)
	if pe.MatchesPattern(ctr.Name) || pe.MatchesPattern(ctr.Image) {
		return CategoryDoNotTouch
	}

	return CategoryProbablySafe
}

func CategorizeNetwork(
	net NetworkInfo,
	protectionPatterns []string,
) SafetyCategory {
	switch net.Name {
	case "bridge", "host", "none":
		return CategoryDoNotTouch
	}

	if net.Containers > 0 {
		return CategoryDoNotTouch
	}

	pe := NewProtectionEngine(protectionPatterns)
	if pe.MatchesPattern(net.Name) {
		return CategoryDoNotTouch
	}

	return CategorySafe
}

func CategoryLabel(cat SafetyCategory) string {
	switch cat {
	case CategorySafe:
		return "SAFE"
	case CategoryProbablySafe:
		return "PROBABLY SAFE"
	case CategoryCheckFirst:
		return "CHECK FIRST"
	case CategoryDoNotTouch:
		return "DO NOT TOUCH"
	default:
		return "UNKNOWN"
	}
}

func looksLikeHash(name string) bool {
	if len(name) < 32 {
		return false
	}
	for _, c := range name {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
