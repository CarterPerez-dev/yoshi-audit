// ©AngelaMos | 2026
// app.go

package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/CarterPerez-dev/yoshi-audit/internal/ui/dashboard"
)

type Tab int

const (
	TabDashboard Tab = iota
	TabDocker
	TabAudit
)

type TickMsg time.Time

type App struct {
	activeTab Tab
	width     int
	height    int
	paused    bool
	dashboard dashboard.Dashboard
}

func NewApp() App {
	return App{
		activeTab: TabDashboard,
		dashboard: dashboard.NewDashboard(),
	}
}

func doTick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (a App) Init() tea.Cmd {
	return tea.Batch(doTick(), dashboard.FetchStats)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	case dashboard.StatsMsg:
		a.dashboard, _ = a.dashboard.Update(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "1":
			a.activeTab = TabDashboard
		case "2":
			a.activeTab = TabDocker
		case "3":
			a.activeTab = TabAudit
		case "tab":
			a.activeTab = (a.activeTab + 1) % 3
		case "p":
			a.paused = !a.paused
		case "m", "c", "g", "n":
			if a.activeTab == TabDashboard {
				a.dashboard, _ = a.dashboard.Update(msg)
			}
		}
	case TickMsg:
		if !a.paused {
			return a, tea.Batch(doTick(), dashboard.FetchStats)
		}
		return a, doTick()
	}
	return a, nil
}

func (a App) View() string {
	if a.width == 0 {
		return "Loading..."
	}

	tabs := a.renderTabs()
	content := a.renderContent()
	status := a.renderStatusBar()

	body := lipgloss.JoinVertical(lipgloss.Left, tabs, content, status)

	frame := FrameStyle.
		Width(a.width - 2).
		Height(a.height - 2).
		Render(body)

	return frame
}

func (a App) renderTabs() string {
	labels := []struct {
		key  string
		name string
		tab  Tab
	}{
		{"1", "DASHBOARD", TabDashboard},
		{"2", "DOCKER", TabDocker},
		{"3", "AUDIT", TabAudit},
	}

	tabs := make([]string, len(labels))
	for i, l := range labels {
		text := fmt.Sprintf("%s %s-%s %s %s", TabHeader, "1", l.key, l.name, TabHeader)
		if a.activeTab == l.tab {
			tabs[i] = ActiveTabStyle.Render(text)
		} else {
			tabs[i] = InactiveTabStyle.Render(text)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (a App) renderContent() string {
	switch a.activeTab {
	case TabDashboard:
		return a.dashboard.View(a.width-4, a.height-6)
	case TabDocker:
		return "Docker Prune Manager - TODO"
	case TabAudit:
		return "System Audit - TODO"
	default:
		return "Unknown tab"
	}
}

func (a App) renderStatusBar() string {
	if a.paused {
		return HelpStyle.Render("[P]aused [R]efresh [Q]uit [?]Help")
	}
	return HelpStyle.Render("[P]ause [R]efresh [Q]uit [?]Help")
}
