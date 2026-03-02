// ©AngelaMos | 2026
// dockertab.go

package dockertab

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/CarterPerez-dev/yoshi-audit/internal/config"
	"github.com/CarterPerez-dev/yoshi-audit/internal/docker"
	"github.com/CarterPerez-dev/yoshi-audit/internal/system"
	"github.com/CarterPerez-dev/yoshi-audit/internal/ui/theme"
)

type SubTab int

const (
	SubTabImages SubTab = iota
	SubTabContainers
	SubTabVolumes
	SubTabNetworks
)

type SelectState int

const (
	Unselected SelectState = iota
	Selected
	Protected
)

type DockerItem struct {
	ID       string
	Name     string
	Size     int64
	Age      string
	Status   string
	Category docker.SafetyCategory
	State    SelectState
	Extra    string
}

type DockerDataMsg struct {
	Images     []DockerItem
	Containers []DockerItem
	Volumes    []DockerItem
	BuildCache int64
	Err        error
}

type DeleteResultMsg struct {
	Deleted int
	Freed   int64
	Err     error
}

type DockerTab struct {
	subTab     SubTab
	images     []DockerItem
	containers []DockerItem
	volumes    []DockerItem
	buildCache int64
	cursor     int
	confirming bool
	confirmMsg string
	protection *docker.ProtectionEngine
	config     config.Config
	err        error
}

func NewDockerTab(cfg config.Config) DockerTab {
	return DockerTab{
		subTab:     SubTabImages,
		protection: docker.NewProtectionEngine(cfg.ProtectionPatterns),
		config:     cfg,
	}
}

func FetchDockerData(cfg config.Config) tea.Msg {
	cli, err := docker.NewClient()
	if err != nil {
		return DockerDataMsg{Err: err}
	}
	defer cli.Close()

	images, containers, volumes, cache, err := cli.GetDiskUsage()
	if err != nil {
		return DockerDataMsg{Err: err}
	}

	pe := docker.NewProtectionEngine(cfg.ProtectionPatterns)

	var imgItems []DockerItem
	for _, img := range images {
		cat := docker.CategorizeImage(img, cfg.ProtectionPatterns)
		name := img.Repository
		if img.Tag != "<none>" {
			name = img.Repository + ":" + img.Tag
		}
		extra := ""
		if img.Dangling {
			extra = "dangling"
		}
		if img.Containers > 0 {
			extra = fmt.Sprintf("%d containers", img.Containers)
		}
		state := Unselected
		fullName := img.Repository
		if img.Tag != "<none>" {
			fullName = img.Repository + ":" + img.Tag
		}
		if pe.IsProtected(fullName) || pe.IsProtected(img.Repository) {
			state = Protected
		}
		imgItems = append(imgItems, DockerItem{
			ID:       img.ID,
			Name:     name,
			Size:     img.Size,
			Age:      formatAge(img.Created),
			Category: cat,
			State:    state,
			Extra:    extra,
		})
	}

	var ctrItems []DockerItem
	for _, ctr := range containers {
		cat := docker.CategorizeContainer(ctr, cfg.ProtectionPatterns)
		state := Unselected
		if pe.IsProtected(ctr.Name) || pe.IsProtected(ctr.Image) {
			state = Protected
		}
		extra := ctr.State
		if ctr.Running {
			extra = "running"
		}
		ctrItems = append(ctrItems, DockerItem{
			ID:       ctr.ID,
			Name:     ctr.Name,
			Size:     ctr.Size,
			Age:      formatAge(ctr.Created),
			Status:   ctr.Status,
			Category: cat,
			State:    state,
			Extra:    extra,
		})
	}

	var volItems []DockerItem
	for _, vol := range volumes {
		cat := docker.CategorizeVolume(vol, cfg.ProtectionPatterns)
		state := Unselected
		if pe.IsProtected(vol.Name) {
			state = Protected
		}
		extra := ""
		if vol.Links > 0 {
			extra = fmt.Sprintf("%d links", vol.Links)
		}
		volItems = append(volItems, DockerItem{
			ID:       vol.Name,
			Name:     vol.Name,
			Size:     vol.Size,
			Age:      formatAge(vol.Created),
			Category: cat,
			State:    state,
			Extra:    extra,
		})
	}

	return DockerDataMsg{
		Images:     imgItems,
		Containers: ctrItems,
		Volumes:    volItems,
		BuildCache: cache.TotalSize,
	}
}

func (dt DockerTab) Update(msg tea.Msg) (DockerTab, tea.Cmd) {
	switch msg := msg.(type) {
	case DockerDataMsg:
		if msg.Err != nil {
			dt.err = msg.Err
			return dt, nil
		}
		dt.images = msg.Images
		dt.containers = msg.Containers
		dt.volumes = msg.Volumes
		dt.buildCache = msg.BuildCache
		dt.err = nil
	case DeleteResultMsg:
		if msg.Err != nil {
			dt.err = msg.Err
		}
		dt.confirming = false
		dt.confirmMsg = ""
		return dt, func() tea.Msg { return FetchDockerData(dt.config) }
	case tea.KeyMsg:
		if dt.confirming {
			switch msg.String() {
			case "y":
				return dt, dt.executeDelete()
			case "escape", "n":
				dt.confirming = false
				dt.confirmMsg = ""
			}
			return dt, nil
		}

		switch msg.String() {
		case "i":
			dt.subTab = SubTabImages
			dt.cursor = 0
		case "c":
			dt.subTab = SubTabContainers
			dt.cursor = 0
		case "v":
			dt.subTab = SubTabVolumes
			dt.cursor = 0
		case "n":
			dt.subTab = SubTabNetworks
			dt.cursor = 0
		case "up", "k":
			if dt.cursor > 0 {
				dt.cursor--
			}
		case "down", "j":
			items := dt.currentItems()
			if dt.cursor < len(items)-1 {
				dt.cursor++
			}
		case " ":
			items := dt.currentItems()
			if dt.cursor < len(items) && items[dt.cursor].State != Protected {
				if items[dt.cursor].State == Selected {
					items[dt.cursor].State = Unselected
				} else {
					items[dt.cursor].State = Selected
				}
				dt.setCurrentItems(items)
			}
		case "a":
			items := dt.currentItems()
			for idx := range items {
				if items[idx].State != Protected {
					items[idx].State = Selected
				}
			}
			dt.setCurrentItems(items)
		case "p":
			items := dt.currentItems()
			if dt.cursor < len(items) {
				if items[dt.cursor].State == Protected {
					items[dt.cursor].State = Unselected
					dt.protection.Unprotect(items[dt.cursor].Name)
				} else {
					items[dt.cursor].State = Protected
					dt.protection.Protect(items[dt.cursor].Name)
				}
				dt.setCurrentItems(items)
			}
		case "d":
			selected := dt.selectedItems()
			if len(selected) > 0 {
				dt.confirming = true
				dt.confirmMsg = dt.buildConfirmMsg(selected)
			}
		case "1":
			dt.applyPreset(0)
		case "2":
			dt.applyPreset(1)
		case "3":
			dt.applyPreset(2)
		case "4":
			dt.applyPreset(3)
		}
	}
	return dt, nil
}

func (dt DockerTab) View(width, height int) string {
	var b strings.Builder

	subTabs := []struct {
		key    string
		label  string
		active bool
	}{
		{"I", "mages", dt.subTab == SubTabImages},
		{"C", "ontainers", dt.subTab == SubTabContainers},
		{"V", "olumes", dt.subTab == SubTabVolumes},
		{"N", "etworks", dt.subTab == SubTabNetworks},
	}

	tabLine := "  "
	for idx, st := range subTabs {
		if st.active {
			tabLine += theme.ActiveTabStyle.Render("["+st.key+"]") + theme.ActiveTabStyle.Render(st.label)
		} else {
			tabLine += theme.HelpStyle.Render("[" + st.key + "]" + st.label)
		}
		if idx < len(subTabs)-1 {
			tabLine += "  "
		}
	}
	b.WriteString(tabLine + "\n\n")

	totalDisk := dt.totalDiskUsage()
	reclaimable := dt.reclaimableSize()
	b.WriteString(fmt.Sprintf("  Total Disk: %s    Reclaimable: %s    Build Cache: %s\n",
		formatSize(totalDisk),
		theme.StatusWarn.Render(formatSize(reclaimable)),
		formatSize(dt.buildCache)))
	b.WriteString("\n")

	if dt.err != nil {
		b.WriteString(fmt.Sprintf("  %s\n", theme.DangerStyle.Render(fmt.Sprintf("Error: %v", dt.err))))
		return b.String()
	}

	items := dt.currentItems()

	if len(items) == 0 {
		b.WriteString("  " + theme.HelpStyle.Render("No items found") + "\n")
	} else {
		headerStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
		nameWidth := width - 60
		if nameWidth < 20 {
			nameWidth = 20
		}
		if nameWidth > 40 {
			nameWidth = 40
		}
		header := fmt.Sprintf("  %-3s %-*s %10s %14s %15s %14s",
			"", nameWidth, "NAME", "SIZE", "AGE", "EXTRA", "CATEGORY")
		b.WriteString(headerStyle.Render(header) + "\n")

		lineWidth := width - 4
		if lineWidth < 40 {
			lineWidth = 40
		}
		b.WriteString("  " + strings.Repeat("\u2500", lineWidth) + "\n")

		maxItems := height - 12
		if maxItems < 1 {
			maxItems = 1
		}
		startIdx := 0
		if dt.cursor >= maxItems {
			startIdx = dt.cursor - maxItems + 1
		}
		endIdx := startIdx + maxItems
		if endIdx > len(items) {
			endIdx = len(items)
		}

		for idx := startIdx; idx < endIdx; idx++ {
			item := items[idx]
			indicator := "[ ]"
			switch item.State {
			case Selected:
				indicator = "[x]"
			case Protected:
				indicator = "[*]"
			}

			name := item.Name
			if len(name) > nameWidth {
				if nameWidth > 3 {
					name = name[:nameWidth-3] + "..."
				} else {
					name = name[:nameWidth]
				}
			}

			catLabel := docker.CategoryLabel(item.Category)
			extra := item.Extra
			if len(extra) > 15 {
				extra = extra[:12] + "..."
			}

			line := fmt.Sprintf("  %s %-*s %10s %14s %15s %14s",
				indicator,
				nameWidth, name,
				formatSize(item.Size),
				item.Age,
				extra,
				catLabel)

			lineStyle := dt.styleForCategory(item.Category)
			if item.State == Protected {
				lineStyle = theme.ProtectedStyle
			}

			cursor := "  "
			if idx == dt.cursor {
				cursor = "> "
				lineStyle = lineStyle.Underline(true)
			}

			b.WriteString(cursor + lineStyle.Render(line) + "\n")
		}

		if len(items) > maxItems {
			paginationStyle := lipgloss.NewStyle().Foreground(theme.CoinGold)
			b.WriteString(fmt.Sprintf("\n  %s\n",
				paginationStyle.Render(fmt.Sprintf("Showing %d-%d of %d items",
					startIdx+1, endIdx, len(items)))))
		}
	}

	if dt.confirming {
		b.WriteString("\n")
		b.WriteString("  " + theme.DangerStyle.Render(dt.confirmMsg) + "\n")
		b.WriteString("  " + theme.DangerStyle.Render("Press [Y] to confirm, [Esc] to cancel") + "\n")
	}

	if len(dt.config.PrunePresets) > 0 {
		b.WriteString("\n  " + theme.TitleStyle.Render("PRESETS") + "\n")
		for idx, preset := range dt.config.PrunePresets {
			if idx >= 4 {
				break
			}
			b.WriteString(fmt.Sprintf("  %s %s\n",
				theme.ActiveTabStyle.Render(fmt.Sprintf("[%d]", idx+1)),
				theme.HelpStyle.Render(preset.Name)))
		}
	}

	b.WriteString("\n")
	b.WriteString("  " + theme.HelpStyle.Render("[Space]Select [A]ll [P]rotect [D]elete [R]efresh [Esc]Cancel"))

	return b.String()
}

func (dt *DockerTab) currentItems() []DockerItem {
	switch dt.subTab {
	case SubTabImages:
		return dt.images
	case SubTabContainers:
		return dt.containers
	case SubTabVolumes:
		return dt.volumes
	default:
		return nil
	}
}

func (dt *DockerTab) setCurrentItems(items []DockerItem) {
	switch dt.subTab {
	case SubTabImages:
		dt.images = items
	case SubTabContainers:
		dt.containers = items
	case SubTabVolumes:
		dt.volumes = items
	}
}

func (dt DockerTab) selectedItems() []DockerItem {
	var selected []DockerItem
	for _, item := range dt.currentItems() {
		if item.State == Selected {
			selected = append(selected, item)
		}
	}
	return selected
}

func (dt DockerTab) buildConfirmMsg(selected []DockerItem) string {
	var totalSize int64
	var names []string
	for _, item := range selected {
		totalSize += item.Size
		names = append(names, item.Name)
	}

	msg := fmt.Sprintf("Delete %d items (%s)?", len(selected), formatSize(totalSize))
	if len(names) <= 5 {
		msg += "\n"
		for _, n := range names {
			msg += fmt.Sprintf("    - %s\n", n)
		}
	}
	return msg
}

func (dt DockerTab) executeDelete() tea.Cmd {
	selected := dt.selectedItems()
	subTab := dt.subTab

	return func() tea.Msg {
		cli, err := docker.NewClient()
		if err != nil {
			return DeleteResultMsg{Err: err}
		}
		defer cli.Close()

		var deleted int
		var freed int64

		for _, item := range selected {
			switch subTab {
			case SubTabImages:
				err = cli.RemoveImage(item.ID)
			case SubTabContainers:
				err = cli.RemoveContainer(item.ID)
			case SubTabVolumes:
				err = cli.RemoveVolume(item.ID)
			default:
				continue
			}
			if err == nil {
				deleted++
				freed += item.Size
			}
		}

		return DeleteResultMsg{
			Deleted: deleted,
			Freed:   freed,
		}
	}
}

func (dt *DockerTab) applyPreset(index int) {
	if index >= len(dt.config.PrunePresets) {
		return
	}

	preset := dt.config.PrunePresets[index]

	var imgInfos []docker.ImageInfo
	for _, item := range dt.images {
		repo := item.Name
		tag := "<none>"
		if parts := strings.SplitN(item.Name, ":", 2); len(parts) == 2 {
			repo = parts[0]
			tag = parts[1]
		}
		imgInfos = append(imgInfos, docker.ImageInfo{
			ID:         item.ID,
			Repository: repo,
			Tag:        tag,
			Size:       item.Size,
			Dangling:   item.Extra == "dangling",
		})
	}

	var volInfos []docker.VolumeInfo
	for _, item := range dt.volumes {
		volInfos = append(volInfos, docker.VolumeInfo{
			Name: item.Name,
			Size: item.Size,
		})
	}

	imageIDs, volumeNames := docker.ApplyPreset(preset, imgInfos, volInfos, dt.protection)

	imageSet := make(map[string]bool, len(imageIDs))
	for _, id := range imageIDs {
		imageSet[id] = true
	}
	for idx := range dt.images {
		if dt.images[idx].State != Protected && imageSet[dt.images[idx].ID] {
			dt.images[idx].State = Selected
		}
	}

	volumeSet := make(map[string]bool, len(volumeNames))
	for _, name := range volumeNames {
		volumeSet[name] = true
	}
	for idx := range dt.volumes {
		if dt.volumes[idx].State != Protected && volumeSet[dt.volumes[idx].Name] {
			dt.volumes[idx].State = Selected
		}
	}
}

func (dt DockerTab) totalDiskUsage() int64 {
	var total int64
	for _, item := range dt.images {
		total += item.Size
	}
	for _, item := range dt.containers {
		total += item.Size
	}
	for _, item := range dt.volumes {
		total += item.Size
	}
	total += dt.buildCache
	return total
}

func (dt DockerTab) reclaimableSize() int64 {
	var total int64
	for _, item := range dt.images {
		if item.State == Selected {
			total += item.Size
		}
	}
	for _, item := range dt.containers {
		if item.State == Selected {
			total += item.Size
		}
	}
	for _, item := range dt.volumes {
		if item.State == Selected {
			total += item.Size
		}
	}
	return total
}

func (dt DockerTab) styleForCategory(cat docker.SafetyCategory) lipgloss.Style {
	switch cat {
	case docker.CategorySafe:
		return lipgloss.NewStyle().Foreground(theme.OneUpGreen)
	case docker.CategoryProbablySafe:
		return lipgloss.NewStyle().Foreground(theme.MarioBlue)
	case docker.CategoryCheckFirst:
		return lipgloss.NewStyle().Foreground(theme.BrickBrown)
	case docker.CategoryDoNotTouch:
		return theme.DangerStyle
	default:
		return lipgloss.NewStyle().Foreground(theme.TextWhite)
	}
}

func formatAge(created time.Time) string {
	if created.IsZero() {
		return "unknown"
	}
	dur := time.Since(created)
	days := int(dur.Hours() / 24)
	if days < 1 {
		hours := int(dur.Hours())
		if hours < 1 {
			return "just now"
		}
		return fmt.Sprintf("%dh ago", hours)
	}
	if days < 7 {
		return fmt.Sprintf("%dd ago", days)
	}
	weeks := days / 7
	if weeks < 5 {
		return fmt.Sprintf("%dw ago", weeks)
	}
	months := days / 30
	if months < 12 {
		return fmt.Sprintf("%dmo ago", months)
	}
	years := days / 365
	return fmt.Sprintf("%dy ago", years)
}

func formatSize(b int64) string {
	if b < 0 {
		return "0 B"
	}
	return system.FormatBytes(uint64(b))
}
