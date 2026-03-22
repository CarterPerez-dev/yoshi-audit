// ©AngelaMos | 2026
// app_test.go

package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app.activeTab != TabDashboard {
		t.Error("should start on dashboard")
	}
}

func TestAppInit(t *testing.T) {
	app := NewApp()
	cmd := app.Init()
	if cmd == nil {
		t.Error("Init should return a command")
	}
}

func TestAppTabSwitch(t *testing.T) {
	app := NewApp()
	app.showSplash = false

	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	a, ok := model.(App)
	if !ok {
		t.Fatal("Update did not return App")
	}
	if a.activeTab != TabDocker {
		t.Error("should be on docker tab")
	}

	model, _ = a.Update(tea.KeyMsg{Type: tea.KeyTab})
	a, ok = model.(App)
	if !ok {
		t.Fatal("Update did not return App")
	}
	if a.activeTab != TabAudit {
		t.Error("should be on audit tab")
	}
}

func TestAppPause(t *testing.T) {
	app := NewApp()
	app.showSplash = false

	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	a, ok := model.(App)
	if !ok {
		t.Fatal("Update did not return App")
	}
	if !a.paused {
		t.Error("should be paused")
	}
}

func TestAppWindowSize(t *testing.T) {
	app := NewApp()
	app.showSplash = false

	model, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	a, ok := model.(App)
	if !ok {
		t.Fatal("Update did not return App")
	}
	if a.width != 120 || a.height != 40 {
		t.Error("should store dimensions")
	}
}
