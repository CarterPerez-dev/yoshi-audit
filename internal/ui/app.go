// ©AngelaMos | 2026
// app.go

package ui

import (
	"fmt"
	"time"

	"github.com/CarterPerez-dev/yoshi-audit/internal/config"
	"github.com/CarterPerez-dev/yoshi-audit/internal/ui/audittab"
	"github.com/CarterPerez-dev/yoshi-audit/internal/ui/dashboard"
	"github.com/CarterPerez-dev/yoshi-audit/internal/ui/dockertab"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab int

const (
	TabDashboard Tab = iota
	TabDocker
	TabAudit
)

const (
	keyCtrlC = "ctrl+c"
	keyTab   = "tab"
)

type TickMsg time.Time

type App struct {
	activeTab  Tab
	width      int
	height     int
	paused     bool
	showSplash bool
	splash     SplashModel
	dashboard  dashboard.Dashboard
	dockerTab  dockertab.DockerTab
	auditTab   audittab.AuditTab
	cfg        config.Config
}

func NewApp() App {
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		cfg = config.Default()
	}
	return App{
		activeTab:  TabDashboard,
		showSplash: true,
		splash:     NewSplash(),
		dashboard:  dashboard.NewDashboard(),
		dockerTab:  dockertab.NewDockerTab(cfg),
		auditTab:   audittab.NewAuditTab(),
		cfg:        cfg,
	}
}

func (a App) doTick() tea.Cmd {
	interval := time.Duration(a.cfg.RefreshInterval) * time.Second
	if interval < time.Second {
		interval = time.Second
	}
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (a App) Init() tea.Cmd {
	if a.showSplash {
		return a.splash.Init()
	}
	return a.initDataFetches()
}

func (a App) initDataFetches() tea.Cmd {
	cfg := a.cfg
	auditScanner := a.auditTab.Scanner()
	auditRSS := a.auditTab.RSSHistory()
	return tea.Batch(a.doTick(), dashboard.FetchStats, func() tea.Msg {
		return dockertab.FetchDockerData(cfg)
	}, func() tea.Msg {
		return audittab.FetchAuditData(auditScanner, auditRSS)
	})
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if a.showSplash {
		if wsm, ok := msg.(tea.WindowSizeMsg); ok {
			a.width = wsm.Width
			a.height = wsm.Height
		}
		var cmd tea.Cmd
		a.splash, cmd = a.splash.Update(msg)
		if a.splash.Done() {
			a.showSplash = false
			return a, a.initDataFetches()
		}
		return a, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	case dashboard.StatsMsg:
		a.dashboard, _ = a.dashboard.Update(msg)
	case dockertab.DockerDataMsg:
		a.dockerTab, _ = a.dockerTab.Update(msg)
	case audittab.AuditDataMsg:
		a.auditTab, _ = a.auditTab.Update(msg)
	case audittab.ExportResultMsg:
		a.auditTab, _ = a.auditTab.Update(msg)
	case dockertab.DeleteResultMsg:
		var cmd tea.Cmd
		a.dockerTab, cmd = a.dockerTab.Update(msg)
		return a, cmd
	case tea.KeyMsg:
		if a.activeTab == TabDocker {
			switch msg.String() {
			case "q", keyCtrlC:
				return a, tea.Quit
			case keyTab:
				a.activeTab = (a.activeTab + 1) % 3
			case "r":
				cfg := a.cfg
				return a, func() tea.Msg { return dockertab.FetchDockerData(cfg) }
			default:
				var cmd tea.Cmd
				a.dockerTab, cmd = a.dockerTab.Update(msg)
				return a, cmd
			}
			return a, nil
		}

		if a.activeTab == TabAudit {
			switch msg.String() {
			case "q", keyCtrlC:
				return a, tea.Quit
			case keyTab:
				a.activeTab = (a.activeTab + 1) % 3
			default:
				var cmd tea.Cmd
				a.auditTab, cmd = a.auditTab.Update(msg)
				return a, cmd
			}
			return a, nil
		}

		switch msg.String() {
		case "q", keyCtrlC:
			return a, tea.Quit
		case "1":
			a.activeTab = TabDashboard
		case "2":
			a.activeTab = TabDocker
		case "3":
			a.activeTab = TabAudit
		case keyTab:
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
			return a, tea.Batch(a.doTick(), dashboard.FetchStats)
		}
		return a, a.doTick()
	}
	return a, nil
}

func (a App) View() string {
	if a.showSplash {
		return a.splash.View()
	}

	if a.width == 0 {
		return "Loading..."
	}

	tabs := a.renderTabs()
	content := a.renderContent()
	status := a.renderStatusBar()

	body := lipgloss.JoinVertical(lipgloss.Left, tabs, content, status)

	borderColor := PipeGreen
	switch a.activeTab {
	case TabDashboard:
		borderColor = PipeGreen
	case TabDocker:
		borderColor = MarioBlue
	case TabAudit:
		borderColor = MarioRed
	}

	frame := FrameStyle.
		BorderForeground(borderColor).
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
		text := fmt.Sprintf(
			"%s [%s] %s %s",
			TabHeader,
			l.key,
			l.name,
			TabHeader,
		)
		if a.activeTab == l.tab {
			tabs[i] = ActiveTabStyle.Render(text)
		} else {
			tabs[i] = InactiveTabStyle.Render(text)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (a App) renderContent() string {
	artWidth := 42
	showArt := a.width > 110

	contentWidth := a.width - 4
	if showArt {
		contentWidth = a.width - 4 - artWidth
	}

	var content string
	var art string
	var artColor lipgloss.Color

	switch a.activeTab {
	case TabDashboard:
		content = a.dashboard.View(contentWidth, a.height-6)
		art = StarPower
		artColor = CoinGold
	case TabDocker:
		content = a.dockerTab.View(contentWidth, a.height-6)
		art = ShyGuy
		artColor = MarioBlue
	case TabAudit:
		content = a.auditTab.View(contentWidth, a.height-6)
		art = KoopaShell
		artColor = MarioRed
	default:
		return "Unknown tab"
	}

	if showArt {
		contentBox := lipgloss.NewStyle().
			Width(contentWidth).
			MaxWidth(contentWidth).
			Render(content)
		artStyled := lipgloss.NewStyle().
			Foreground(artColor).
			Width(artWidth).
			PaddingTop(2).
			Render(art)
		return lipgloss.JoinHorizontal(lipgloss.Top, contentBox, artStyled)
	}
	return content
}

func (a App) renderStatusBar() string {
	if a.paused {
		return HelpStyle.Render("[P]aused [R]efresh [Q]uit [?]Help")
	}
	return HelpStyle.Render("[P]ause [R]efresh [Q]uit [?]Help")
}
