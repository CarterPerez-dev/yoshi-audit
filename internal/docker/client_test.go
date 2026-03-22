// ©AngelaMos | 2026
// client_test.go

package docker

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer c.Close() //nolint:errcheck
}

func TestGetDiskUsage(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer c.Close() //nolint:errcheck

	images, containers, volumes, cache, err := c.GetDiskUsage()
	if err != nil {
		t.Fatalf("GetDiskUsage failed: %v", err)
	}
	if len(images) == 0 {
		t.Log("Warning: no images found")
	}
	_ = containers
	_ = volumes
	_ = cache
}

func TestListImages(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer c.Close() //nolint:errcheck

	images, err := c.ListImages()
	if err != nil {
		t.Fatalf("ListImages failed: %v", err)
	}
	for _, img := range images {
		if img.ID == "" {
			t.Error("image ID should not be empty")
		}
	}
}

func TestListContainers(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer c.Close() //nolint:errcheck

	containers, err := c.ListContainers()
	if err != nil {
		t.Fatalf("ListContainers failed: %v", err)
	}
	for _, ctr := range containers {
		if ctr.ID == "" {
			t.Error("container ID should not be empty")
		}
	}
}

func TestListNetworks(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer c.Close() //nolint:errcheck

	networks, err := c.ListNetworks()
	if err != nil {
		t.Fatalf("ListNetworks failed: %v", err)
	}
	if len(networks) < 3 {
		t.Errorf(
			"expected at least 3 default networks (bridge/host/none), got %d",
			len(networks),
		)
	}
	for _, n := range networks {
		if n.ID == "" {
			t.Error("network ID should not be empty")
		}
		if n.Name == "" {
			t.Error("network name should not be empty")
		}
	}
}

func TestListVolumes(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer c.Close() //nolint:errcheck

	volumes, err := c.ListVolumes()
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	for _, vol := range volumes {
		if vol.Name == "" {
			t.Error("volume name should not be empty")
		}
	}
}
