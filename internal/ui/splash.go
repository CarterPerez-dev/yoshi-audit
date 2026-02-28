// ©AngelaMos | 2026
// splash.go

package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SplashDoneMsg struct{}

type SplashModel struct {
	done bool
}

func NewSplash() SplashModel {
	return SplashModel{}
}

func (s SplashModel) Init() tea.Cmd {
	return tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg {
		return SplashDoneMsg{}
	})
}

func (s SplashModel) Update(msg tea.Msg) (SplashModel, tea.Cmd) {
	switch msg.(type) {
	case SplashDoneMsg:
		s.done = true
	case tea.KeyMsg:
		s.done = true
	}
	return s, nil
}

func (s SplashModel) View() string {
	art := TitleStyle.Render(SplashYoshi)

	lines := strings.Split(art, "\n")
	centered := make([]string, len(lines))
	for i, line := range lines {
		centered[i] = lipgloss.PlaceHorizontal(80, lipgloss.Center, line)
	}

	return lipgloss.JoinVertical(lipgloss.Center, centered...)
}

func (s SplashModel) Done() bool {
	return s.done
}
