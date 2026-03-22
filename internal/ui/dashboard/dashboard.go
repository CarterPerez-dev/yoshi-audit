// ©AngelaMos | 2026
// dashboard.go

package dashboard

import (
	"fmt"
	"sort"
	"strings"

	"github.com/CarterPerez-dev/yoshi-audit/internal/system"
	"github.com/CarterPerez-dev/yoshi-audit/internal/ui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SortMode int

const (
	SortByMemory SortMode = iota
	SortByCPU
	SortByGPU
	SortByPID
	SortByName
)

type StatsMsg struct {
	CPU      float64
	Memory   system.MemoryInfo
	Disk     system.DiskInfo
	GPU      system.GPUInfo
	GPUProcs []system.GPUProcess
	Procs    []system.ProcessInfo
	Err      error
}

type Dashboard struct {
	cpu      float64
	memory   system.MemoryInfo
	disk     system.DiskInfo
	gpu      system.GPUInfo
	gpuProcs []system.GPUProcess
	procs    []system.ProcessInfo
	sortMode SortMode
	err      error
}

func NewDashboard() Dashboard {
	return Dashboard{sortMode: SortByMemory}
}

func FetchStats() tea.Msg {
	var stats StatsMsg

	cpu, err := system.GetCPUUsage()
	if err != nil {
		stats.Err = err
		return stats
	}
	stats.CPU = cpu

	mem, err := system.GetMemoryInfo()
	if err != nil {
		stats.Err = err
		return stats
	}
	stats.Memory = mem

	disk, err := system.GetDiskUsage("/")
	if err != nil {
		stats.Err = err
		return stats
	}
	stats.Disk = disk

	gpuInfo, err := system.GetGPUInfo()
	if err == nil {
		stats.GPU = gpuInfo
	}

	gpuProcs, err := system.GetGPUProcesses()
	if err == nil {
		stats.GPUProcs = gpuProcs
	}

	procs, err := system.GetProcesses()
	if err != nil {
		stats.Err = err
		return stats
	}
	stats.Procs = procs

	return stats
}

func (d Dashboard) Update(msg tea.Msg) (Dashboard, tea.Cmd) {
	switch msg := msg.(type) {
	case StatsMsg:
		d.cpu = msg.CPU
		d.memory = msg.Memory
		d.disk = msg.Disk
		d.gpu = msg.GPU
		d.gpuProcs = msg.GPUProcs
		d.procs = msg.Procs
		d.err = msg.Err
	case tea.KeyMsg:
		switch msg.String() {
		case "m":
			d.sortMode = SortByMemory
		case "c":
			d.sortMode = SortByCPU
		case "g":
			d.sortMode = SortByGPU
		case "n":
			d.sortMode = SortByName
		}
	}
	return d, nil
}

func (d Dashboard) View(width, height int) string {
	if d.err != nil {
		return fmt.Sprintf("Error: %v", d.err)
	}

	var b strings.Builder

	barWidth := width / 3
	if barWidth < 16 {
		barWidth = 16
	}
	if barWidth > 30 {
		barWidth = 30
	}

	labelStyle := lipgloss.NewStyle().Foreground(theme.MarioBlue).Bold(true)

	b.WriteString(fmt.Sprintf("  %s %s %s\n",
		labelStyle.Render("CPU "),
		theme.ProgressBar(d.cpu, barWidth),
		fmt.Sprintf(" %.1f%%", d.cpu)))

	ramPct := d.memory.RAMPercent()
	b.WriteString(fmt.Sprintf(
		"  %s %s %s\n",
		labelStyle.Render("RAM "),
		theme.ProgressBar(ramPct, barWidth),
		fmt.Sprintf(
			" %s/%s",
			system.FormatBytes(d.memory.UsedRAM),
			system.FormatBytes(d.memory.TotalRAM),
		),
	))

	gpuPct := d.gpu.Utilization
	b.WriteString(fmt.Sprintf("  %s %s %s\n",
		labelStyle.Render("GPU "),
		theme.ProgressBar(gpuPct, barWidth),
		fmt.Sprintf(" %.1f%%", gpuPct)))

	vramPct := d.gpu.VRAMPercent()
	b.WriteString(fmt.Sprintf(
		"  %s %s %s\n",
		labelStyle.Render("VRAM"),
		theme.ProgressBar(vramPct, barWidth),
		fmt.Sprintf(
			" %s/%s",
			system.FormatBytes(d.gpu.UsedVRAM),
			system.FormatBytes(d.gpu.TotalVRAM),
		),
	))

	diskPct := d.disk.Percent()
	b.WriteString(fmt.Sprintf(
		"  %s %s %s\n",
		labelStyle.Render("DISK"),
		theme.ProgressBar(diskPct, barWidth),
		fmt.Sprintf(
			" %s/%s",
			system.FormatBytes(d.disk.Used),
			system.FormatBytes(d.disk.Total),
		),
	))

	swapPct := d.memory.SwapPercent()
	b.WriteString(fmt.Sprintf(
		"  %s %s %s\n",
		labelStyle.Render("SWAP"),
		theme.ProgressBar(swapPct, barWidth),
		fmt.Sprintf(
			" %s/%s",
			system.FormatBytes(d.memory.UsedSwap),
			system.FormatBytes(d.memory.TotalSwap),
		),
	))

	b.WriteString("\n")
	b.WriteString("  " + theme.TitleStyle.Render("TOP PROCESSES") + "\n")

	lineWidth := width - 4
	if lineWidth < 40 {
		lineWidth = 40
	}
	b.WriteString("  " + strings.Repeat("\u2500", lineWidth) + "\n")

	header := fmt.Sprintf("  %-8s %-18s %6s %10s %10s",
		"PID", "NAME", "CPU%", "MEM", "GPU MEM")
	b.WriteString(
		lipgloss.NewStyle().Foreground(theme.TextDim).Render(header) + "\n",
	)

	procs := make([]system.ProcessInfo, len(d.procs))
	copy(procs, d.procs)

	gpuMap := make(map[int]uint64)
	for _, gp := range d.gpuProcs {
		gpuMap[gp.PID] = gp.UsedVRAM
	}

	switch d.sortMode {
	case SortByMemory:
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].RSS > procs[j].RSS
		})
	case SortByCPU:
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].CPUPerc > procs[j].CPUPerc
		})
	case SortByGPU:
		sort.Slice(procs, func(i, j int) bool {
			return gpuMap[procs[i].PID] > gpuMap[procs[j].PID]
		})
	case SortByPID:
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].PID < procs[j].PID
		})
	case SortByName:
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].Name < procs[j].Name
		})
	}

	maxProcs := 15
	availableRows := height - 12
	if availableRows < maxProcs && availableRows > 0 {
		maxProcs = availableRows
	}
	if maxProcs > len(procs) {
		maxProcs = len(procs)
	}

	for i := 0; i < maxProcs; i++ {
		p := procs[i]
		gpuMem := "\u2014"
		if vram, ok := gpuMap[p.PID]; ok && vram > 0 {
			gpuMem = system.FormatBytes(vram)
		}

		line := fmt.Sprintf("  %-8d %-18s %6.1f %10s %10s",
			p.PID,
			system.FormatProcessName(p.Name, 18),
			p.CPUPerc,
			system.FormatBytes(p.RSS),
			gpuMem)
		b.WriteString(line + "\n")
	}

	b.WriteString("\n")

	sortLabels := []struct {
		key    string
		label  string
		active bool
	}{
		{"M", "emory", d.sortMode == SortByMemory},
		{"C", "PU", d.sortMode == SortByCPU},
		{"G", "PU", d.sortMode == SortByGPU},
		{"N", "ame", d.sortMode == SortByName},
	}

	sortLine := "  Sort: "
	for i, sl := range sortLabels {
		if sl.active {
			sortLine += theme.ActiveTabStyle.Render(
				"["+sl.key+"]",
			) + theme.ActiveTabStyle.Render(
				sl.label,
			)
		} else {
			sortLine += theme.HelpStyle.Render("[" + sl.key + "]" + sl.label)
		}
		if i < len(sortLabels)-1 {
			sortLine += "  "
		}
	}
	b.WriteString(sortLine)

	return b.String()
}
