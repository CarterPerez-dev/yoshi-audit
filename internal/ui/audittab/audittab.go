// ©AngelaMos | 2026
// audittab.go

package audittab

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/CarterPerez-dev/yoshi-audit/internal/audit"
	"github.com/CarterPerez-dev/yoshi-audit/internal/config"
	"github.com/CarterPerez-dev/yoshi-audit/internal/system"
	"github.com/CarterPerez-dev/yoshi-audit/internal/ui/theme"
)

const (
	statusOKFace   = "~(^.^)~"
	statusWarnFace = ">(o.O)>"
	statusCritFace = "X(x.x)X"
)

type SummaryCount struct {
	OK   int
	Warn int
	Crit int
}

type AuditDataMsg struct {
	Findings []audit.Finding
	Summary  map[audit.FindingType]SummaryCount
	Err      error
}

type AuditTab struct {
	findings   []audit.Finding
	summary    map[audit.FindingType]SummaryCount
	cursor     int
	scanner    *audit.Scanner
	baseline   *audit.Baseline
	rssHistory map[int][]uint64
	lastScan   time.Time
	scanning   bool
	showDetail int
	err        error
}

func NewAuditTab() AuditTab {
	bl := audit.NewBaseline()
	bl.Load(config.DefaultBaselinePath())
	return AuditTab{
		findings:   nil,
		scanner:    audit.NewScanner(bl),
		baseline:   bl,
		rssHistory: make(map[int][]uint64),
		showDetail: -1,
	}
}

func (at AuditTab) Scanner() *audit.Scanner {
	return at.scanner
}

func (at AuditTab) RSSHistory() map[int][]uint64 {
	return at.rssHistory
}

func FetchAuditData(scanner *audit.Scanner, rssHistory map[int][]uint64) tea.Msg {
	procs, err := system.GetProcesses()
	if err != nil {
		return AuditDataMsg{Err: err}
	}

	gpuProcs, _ := system.GetGPUProcesses()

	for _, p := range procs {
		rssHistory[p.PID] = append(rssHistory[p.PID], p.RSS)
	}

	findings := scanner.ScanAll(procs, gpuProcs)
	findings = append(findings, scanner.ScanMemoryLeaks(procs, rssHistory)...)

	summary := make(map[audit.FindingType]SummaryCount)
	allTypes := []audit.FindingType{
		audit.FindingZombie,
		audit.FindingOrphan,
		audit.FindingDaemon,
		audit.FindingKernelThread,
		audit.FindingGPUShadow,
		audit.FindingMemoryLeak,
		audit.FindingUnknownSvc,
	}
	for _, ft := range allTypes {
		summary[ft] = SummaryCount{}
	}
	for _, f := range findings {
		sc := summary[f.Type]
		switch f.Severity {
		case audit.SeverityOK:
			sc.OK++
		case audit.SeverityWarn:
			sc.Warn++
		case audit.SeverityCrit:
			sc.Crit++
		}
		summary[f.Type] = sc
	}

	return AuditDataMsg{
		Findings: findings,
		Summary:  summary,
	}
}

func (at AuditTab) Update(msg tea.Msg) (AuditTab, tea.Cmd) {
	switch msg := msg.(type) {
	case AuditDataMsg:
		if msg.Err != nil {
			at.err = msg.Err
			at.scanning = false
			return at, nil
		}
		at.findings = msg.Findings
		at.summary = msg.Summary
		at.lastScan = time.Now()
		at.scanning = false
		at.err = nil
	case tea.KeyMsg:
		filtered := at.filteredFindings()
		switch msg.String() {
		case "r":
			at.scanning = true
			scanner := at.scanner
			rssHistory := at.rssHistory
			return at, func() tea.Msg {
				return FetchAuditData(scanner, rssHistory)
			}
		case "up", "k":
			if at.cursor > 0 {
				at.cursor--
			}
		case "down", "j":
			if at.cursor < len(filtered)-1 {
				at.cursor++
			}
		case "a":
			if at.cursor < len(filtered) && at.baseline != nil {
				f := filtered[at.cursor]
				at.baseline.Add(f.Name)
				at.baseline.Save(config.DefaultBaselinePath())
				at.scanning = true
				scanner := at.scanner
				rssHistory := at.rssHistory
				return at, func() tea.Msg {
					return FetchAuditData(scanner, rssHistory)
				}
			}
		case "i":
			if at.cursor < len(filtered) {
				target := filtered[at.cursor]
				var newFindings []audit.Finding
				removed := false
				for _, f := range at.findings {
					if !removed && f.PID == target.PID && f.Type == target.Type && f.Name == target.Name {
						removed = true
						continue
					}
					newFindings = append(newFindings, f)
				}
				at.findings = newFindings
				if at.cursor >= len(at.filteredFindings()) && at.cursor > 0 {
					at.cursor--
				}
			}
		case "enter":
			if at.showDetail == at.cursor {
				at.showDetail = -1
			} else {
				at.showDetail = at.cursor
			}
		}
	}
	return at, nil
}

func (at AuditTab) filteredFindings() []audit.Finding {
	var filtered []audit.Finding
	for _, f := range at.findings {
		if f.Severity == audit.SeverityWarn || f.Severity == audit.SeverityCrit {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func (at AuditTab) View(width, height int) string {
	var b strings.Builder

	scanStatus := "Complete"
	if at.scanning {
		scanStatus = "Scanning..."
	}
	if at.err != nil {
		scanStatus = "Error"
	}

	lastScanStr := "never"
	if !at.lastScan.IsZero() {
		dur := time.Since(at.lastScan)
		if dur < time.Minute {
			lastScanStr = "just now"
		} else {
			lastScanStr = fmt.Sprintf("%d min ago", int(dur.Minutes()))
		}
	}

	statusLine := fmt.Sprintf("  SCAN STATUS: %s  |  Last: %s  |  [R]escan",
		scanStatus, lastScanStr)
	b.WriteString(theme.TitleStyle.Render(statusLine) + "\n\n")

	if at.err != nil {
		b.WriteString(fmt.Sprintf("  %s\n", theme.StatusCrit.Render(fmt.Sprintf("Error: %v", at.err))))
		return b.String()
	}

	headerStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
	header := fmt.Sprintf("  %-16s %4s %5s %5s", "FINDINGS", "OK", "WARN", "CRIT")
	b.WriteString(headerStyle.Render(header) + "\n")

	typeLabels := []struct {
		ft    audit.FindingType
		label string
	}{
		{audit.FindingZombie, "Zombies"},
		{audit.FindingOrphan, "Orphans"},
		{audit.FindingDaemon, "Daemons"},
		{audit.FindingKernelThread, "Kernel Threads"},
		{audit.FindingGPUShadow, "GPU Shadows"},
		{audit.FindingMemoryLeak, "Memory Leaks"},
		{audit.FindingUnknownSvc, "Unknown Svcs"},
	}

	for _, tl := range typeLabels {
		sc := at.summary[tl.ft]
		okStr := theme.StatusOK.Render(fmt.Sprintf("%4d", sc.OK))
		warnStr := theme.StatusWarn.Render(fmt.Sprintf("%5d", sc.Warn))
		critStr := theme.StatusCrit.Render(fmt.Sprintf("%5d", sc.Crit))
		b.WriteString(fmt.Sprintf("  %-16s %s %s %s\n", tl.label, okStr, warnStr, critStr))
	}

	b.WriteString("\n")

	totalWarn := 0
	totalCrit := 0
	for _, sc := range at.summary {
		totalWarn += sc.Warn
		totalCrit += sc.Crit
	}

	if totalCrit > 0 {
		face := theme.StatusCrit.Render(fmt.Sprintf("  %s  GAME OVER: %d CRITICAL", statusCritFace, totalCrit))
		b.WriteString(face + "\n")
	} else if totalWarn > 0 {
		face := theme.StatusWarn.Render(fmt.Sprintf("  %s  FIRE FLOWER: %d WARNINGS", statusWarnFace, totalWarn))
		b.WriteString(face + "\n")
	} else {
		face := theme.StatusOK.Render(fmt.Sprintf("  %s  1UP: ALL CLEAR", statusOKFace))
		b.WriteString(face + "\n")
	}

	b.WriteString("\n")

	filtered := at.filteredFindings()

	if totalWarn+totalCrit > 0 {
		lineWidth := width - 4
		if lineWidth < 40 {
			lineWidth = 40
		}
		warnHeader := fmt.Sprintf("  [!] %d WARNINGS ", totalWarn+totalCrit)
		remaining := lineWidth - len(warnHeader)
		if remaining < 0 {
			remaining = 0
		}
		b.WriteString(theme.StatusWarn.Render(warnHeader+strings.Repeat("\u2500", remaining)) + "\n\n")
	}

	maxItems := height - 22
	if maxItems < 1 {
		maxItems = 1
	}
	startIdx := 0
	if at.cursor >= maxItems {
		startIdx = at.cursor - maxItems + 1
	}
	endIdx := startIdx + maxItems
	if endIdx > len(filtered) {
		endIdx = len(filtered)
	}

	for idx := startIdx; idx < endIdx; idx++ {
		f := filtered[idx]
		sevStyle := theme.StatusWarn
		if f.Severity == audit.SeverityCrit {
			sevStyle = theme.StatusCrit
		}

		nameStr := f.Name
		if len(nameStr) > 20 {
			nameStr = nameStr[:17] + "..."
		}

		line := fmt.Sprintf("  %-4s  %-10s  \"%s\" PID %d - %s",
			f.Severity, f.Type, nameStr, f.PID, f.Message)

		cursor := "  "
		if idx == at.cursor {
			cursor = "> "
			line = sevStyle.Underline(true).Render(line)
		} else {
			line = sevStyle.Render(line)
		}

		b.WriteString(cursor + line + "\n")

		if at.showDetail == idx && f.Detail != "" {
			detailStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
			b.WriteString("      " + detailStyle.Render(f.Detail) + "\n")
		}
	}

	if len(filtered) > maxItems {
		b.WriteString(fmt.Sprintf("\n  %s\n",
			theme.HelpStyle.Render(fmt.Sprintf("Showing %d-%d of %d findings",
				startIdx+1, endIdx, len(filtered)))))
	}

	b.WriteString("\n")
	b.WriteString("  " + theme.HelpStyle.Render("[Enter] Inspect  [A]dd to baseline  [I]gnore  [R]escan"))

	return b.String()
}
