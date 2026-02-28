// ©AngelaMos | 2026
// categorize_test.go

package docker

import (
	"testing"
)

func TestCategorizeImage_Dangling(t *testing.T) {
	img := ImageInfo{
		Repository: "<none>",
		Tag:        "<none>",
		Containers: 0,
		Dangling:   true,
	}
	cat := CategorizeImage(img, nil)
	if cat != CategorySafe {
		t.Errorf("dangling image should be CategorySafe, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeImage_InUse(t *testing.T) {
	img := ImageInfo{
		Repository: "nginx",
		Tag:        "latest",
		Containers: 1,
	}
	cat := CategorizeImage(img, nil)
	if cat != CategoryDoNotTouch {
		t.Errorf("image with containers=1 should be CategoryDoNotTouch, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeImage_MatchesProtection(t *testing.T) {
	img := ImageInfo{
		Repository: "certgames-prod",
		Tag:        "latest",
		Containers: 0,
	}
	cat := CategorizeImage(img, []string{"*certgames*"})
	if cat != CategoryDoNotTouch {
		t.Errorf("image matching *certgames* should be CategoryDoNotTouch, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeImage_Named(t *testing.T) {
	img := ImageInfo{
		Repository: "redis",
		Tag:        "7-alpine",
		Containers: 0,
	}
	cat := CategorizeImage(img, nil)
	if cat != CategoryProbablySafe {
		t.Errorf("named image with containers=0 should be CategoryProbablySafe, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeVolume_Backup(t *testing.T) {
	vol := VolumeInfo{
		Name:  "my-backup-data",
		Size:  1024,
		Links: 0,
	}
	cat := CategorizeVolume(vol, nil)
	if cat != CategoryCheckFirst {
		t.Errorf("volume with 'backup' in name should be CategoryCheckFirst, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeVolume_MatchesProtection(t *testing.T) {
	vol := VolumeInfo{
		Name:  "certgames-mongo_data",
		Size:  5000,
		Links: 0,
	}
	cat := CategorizeVolume(vol, []string{"*mongo*"})
	if cat != CategoryDoNotTouch {
		t.Errorf("volume matching *mongo* should be CategoryDoNotTouch, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeVolume_Attached(t *testing.T) {
	vol := VolumeInfo{
		Name:  "web-data",
		Size:  2048,
		Links: 1,
	}
	cat := CategorizeVolume(vol, nil)
	if cat != CategoryDoNotTouch {
		t.Errorf("volume with links > 0 should be CategoryDoNotTouch, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeVolume_UnusedNamed(t *testing.T) {
	vol := VolumeInfo{
		Name:  "old-project-data",
		Size:  4096,
		Links: 0,
	}
	cat := CategorizeVolume(vol, nil)
	if cat != CategoryProbablySafe {
		t.Errorf("unused named volume with size > 0 should be CategoryProbablySafe, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeVolume_HashSmall(t *testing.T) {
	vol := VolumeInfo{
		Name:  "abcdef0123456789abcdef0123456789abcdef01",
		Size:  0,
		Links: 0,
	}
	cat := CategorizeVolume(vol, nil)
	if cat != CategorySafe {
		t.Errorf("hash-named empty volume should be CategorySafe, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeContainer_Running(t *testing.T) {
	ctr := ContainerInfo{
		Name:    "web-server",
		Image:   "nginx:latest",
		Running: true,
		State:   "running",
	}
	cat := CategorizeContainer(ctr, nil)
	if cat != CategoryDoNotTouch {
		t.Errorf("running container should be CategoryDoNotTouch, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeContainer_Stopped(t *testing.T) {
	ctr := ContainerInfo{
		Name:    "old-worker",
		Image:   "worker:v1",
		Running: false,
		State:   "exited",
	}
	cat := CategorizeContainer(ctr, nil)
	if cat != CategoryProbablySafe {
		t.Errorf("stopped container should be CategoryProbablySafe, got %s", CategoryLabel(cat))
	}
}

func TestCategorizeContainer_StoppedProtected(t *testing.T) {
	ctr := ContainerInfo{
		Name:    "certgames-api",
		Image:   "certgames:latest",
		Running: false,
		State:   "exited",
	}
	cat := CategorizeContainer(ctr, []string{"*certgames*"})
	if cat != CategoryDoNotTouch {
		t.Errorf("stopped protected container should be CategoryDoNotTouch, got %s", CategoryLabel(cat))
	}
}

func TestCategoryLabel(t *testing.T) {
	tests := []struct {
		cat  SafetyCategory
		want string
	}{
		{CategorySafe, "SAFE"},
		{CategoryProbablySafe, "PROBABLY SAFE"},
		{CategoryCheckFirst, "CHECK FIRST"},
		{CategoryDoNotTouch, "DO NOT TOUCH"},
		{SafetyCategory(99), "UNKNOWN"},
	}
	for _, tt := range tests {
		got := CategoryLabel(tt.cat)
		if got != tt.want {
			t.Errorf("CategoryLabel(%d) = %q, want %q", tt.cat, got, tt.want)
		}
	}
}
