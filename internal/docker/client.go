// ©AngelaMos | 2026
// client.go

package docker

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

const noneTag = "<none>"

type ImageInfo struct {
	ID         string
	Repository string
	Tag        string
	Size       int64
	SharedSize int64
	UniqueSize int64
	Containers int64
	Created    time.Time
	Dangling   bool
}

type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Status  string
	State   string
	Size    int64
	Created time.Time
	Running bool
}

type VolumeInfo struct {
	Name    string
	Size    int64
	Links   int
	Created time.Time
}

type NetworkInfo struct {
	ID         string
	Name       string
	Driver     string
	Scope      string
	Containers int
}

type BuildCacheInfo struct {
	TotalSize int64
}

type Client struct {
	cli *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) ListImages() ([]ImageInfo, error) {
	ctx := context.Background()
	summaries, err := c.cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var images []ImageInfo
	for _, s := range summaries {
		repo := noneTag
		tag := noneTag
		if len(s.RepoTags) > 0 {
			parts := strings.SplitN(s.RepoTags[0], ":", 2)
			repo = parts[0]
			if len(parts) > 1 {
				tag = parts[1]
			}
		}

		dangling := repo == noneTag && s.Containers == 0

		images = append(images, ImageInfo{
			ID:         s.ID,
			Repository: repo,
			Tag:        tag,
			Size:       s.Size,
			SharedSize: s.SharedSize,
			UniqueSize: s.Size - s.SharedSize,
			Containers: s.Containers,
			Created:    time.Unix(s.Created, 0),
			Dangling:   dangling,
		})
	}
	return images, nil
}

func (c *Client) ListContainers() ([]ContainerInfo, error) {
	ctx := context.Background()
	summaries, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var containers []ContainerInfo
	for _, s := range summaries {
		name := ""
		if len(s.Names) > 0 {
			name = strings.TrimPrefix(s.Names[0], "/")
		}

		containers = append(containers, ContainerInfo{
			ID:      s.ID,
			Name:    name,
			Image:   s.Image,
			Status:  s.Status,
			State:   s.State,
			Size:    s.SizeRw,
			Created: time.Unix(s.Created, 0),
			Running: s.State == "running",
		})
	}
	return containers, nil
}

func (c *Client) ListVolumes() ([]VolumeInfo, error) {
	ctx := context.Background()
	resp, err := c.cli.VolumeList(ctx, volumeListOptions())
	if err != nil {
		return nil, err
	}

	var volumes []VolumeInfo
	for _, v := range resp.Volumes {
		var size int64
		var links int
		if v.UsageData != nil {
			size = v.UsageData.Size
			links = int(v.UsageData.RefCount)
		}

		var created time.Time
		if v.CreatedAt != "" {
			created, _ = time.Parse(time.RFC3339, v.CreatedAt) //nolint:errcheck
		}

		volumes = append(volumes, VolumeInfo{
			Name:    v.Name,
			Size:    size,
			Links:   links,
			Created: created,
		})
	}
	return volumes, nil
}

func (c *Client) GetDiskUsage() ([]ImageInfo, []ContainerInfo, []VolumeInfo, BuildCacheInfo, error) {
	ctx := context.Background()
	du, err := c.cli.DiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		return nil, nil, nil, BuildCacheInfo{}, err
	}

	var images []ImageInfo
	for _, s := range du.Images {
		if s == nil {
			continue
		}
		repo := noneTag
		tag := noneTag
		if len(s.RepoTags) > 0 {
			parts := strings.SplitN(s.RepoTags[0], ":", 2)
			repo = parts[0]
			if len(parts) > 1 {
				tag = parts[1]
			}
		}
		dangling := repo == noneTag && s.Containers == 0
		images = append(images, ImageInfo{
			ID:         s.ID,
			Repository: repo,
			Tag:        tag,
			Size:       s.Size,
			SharedSize: s.SharedSize,
			UniqueSize: s.Size - s.SharedSize,
			Containers: s.Containers,
			Created:    time.Unix(s.Created, 0),
			Dangling:   dangling,
		})
	}

	var containers []ContainerInfo
	for _, s := range du.Containers {
		if s == nil {
			continue
		}
		name := ""
		if len(s.Names) > 0 {
			name = strings.TrimPrefix(s.Names[0], "/")
		}
		containers = append(containers, ContainerInfo{
			ID:      s.ID,
			Name:    name,
			Image:   s.Image,
			Status:  s.Status,
			State:   s.State,
			Size:    s.SizeRw,
			Created: time.Unix(s.Created, 0),
			Running: s.State == "running",
		})
	}

	var volumes []VolumeInfo
	for _, v := range du.Volumes {
		if v == nil {
			continue
		}
		var size int64
		var links int
		if v.UsageData != nil {
			size = v.UsageData.Size
			links = int(v.UsageData.RefCount)
		}
		var created time.Time
		if v.CreatedAt != "" {
			created, _ = time.Parse(time.RFC3339, v.CreatedAt) //nolint:errcheck
		}
		volumes = append(volumes, VolumeInfo{
			Name:    v.Name,
			Size:    size,
			Links:   links,
			Created: created,
		})
	}

	var cacheSize int64
	for _, bc := range du.BuildCache {
		if bc != nil {
			cacheSize += bc.Size
		}
	}

	return images, containers, volumes, BuildCacheInfo{
		TotalSize: cacheSize,
	}, nil
}

func (c *Client) RemoveImage(id string) error {
	ctx := context.Background()
	_, err := c.cli.ImageRemove(
		ctx,
		id,
		image.RemoveOptions{Force: true, PruneChildren: true},
	)
	return err
}

func (c *Client) RemoveContainer(id string) error {
	ctx := context.Background()
	return c.cli.ContainerRemove(
		ctx,
		id,
		container.RemoveOptions{Force: true, RemoveVolumes: true},
	)
}

func (c *Client) RemoveVolume(name string) error {
	ctx := context.Background()
	return c.cli.VolumeRemove(ctx, name, true)
}

func (c *Client) ListNetworks() ([]NetworkInfo, error) {
	ctx := context.Background()
	nets, err := c.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []NetworkInfo
	for _, n := range nets {
		result = append(result, NetworkInfo{
			ID:         n.ID,
			Name:       n.Name,
			Driver:     n.Driver,
			Scope:      n.Scope,
			Containers: len(n.Containers),
		})
	}
	return result, nil
}

func (c *Client) RemoveNetwork(id string) error {
	ctx := context.Background()
	return c.cli.NetworkRemove(ctx, id)
}

func (c *Client) PruneBuildCache() (int64, error) {
	ctx := context.Background()
	report, err := c.cli.BuildCachePrune(
		ctx,
		types.BuildCachePruneOptions{All: true},
	)
	if err != nil {
		return 0, err
	}
	return int64(report.SpaceReclaimed), nil //nolint:gosec
}
