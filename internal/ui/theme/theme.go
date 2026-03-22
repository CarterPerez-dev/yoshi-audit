// ©AngelaMos | 2026 "yoshi"
// theme.go

package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	YoshiGreen = lipgloss.Color("#66BB6A")
	PipeGreen  = lipgloss.Color("#2E7D32")
	MarioBlue  = lipgloss.Color("#4A90D9")
	CoinGold   = lipgloss.Color("#FFD600")
	MarioRed   = lipgloss.Color("#E53935")
	OneUpGreen = lipgloss.Color("#4CAF50")
	SnesBg     = lipgloss.Color("#1A1A2E")
	TextWhite  = lipgloss.Color("#FAFAFA")
	TextDim    = lipgloss.Color("#888888")
	BrickBrown = lipgloss.Color("#8B4513")
)

var MarioBorder = lipgloss.Border{
	Top:         "▀",
	Bottom:      "▄",
	Left:        "█",
	Right:       "█",
	TopLeft:     "▀",
	TopRight:    "▀",
	BottomLeft:  "▄",
	BottomRight: "▄",
}

var (
	FrameStyle = lipgloss.NewStyle().
			Border(MarioBorder).
			BorderForeground(PipeGreen)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(YoshiGreen)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(CoinGold)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(TextDim)

	BarOK = lipgloss.NewStyle().
		Foreground(OneUpGreen)

	BarWarn = lipgloss.NewStyle().
		Foreground(MarioBlue)

	BarCrit = lipgloss.NewStyle().
		Foreground(MarioRed)

	StatusOK = lipgloss.NewStyle().
			Bold(true).
			Foreground(OneUpGreen)

	StatusWarn = lipgloss.NewStyle().
			Bold(true).
			Foreground(MarioBlue)

	StatusCrit = lipgloss.NewStyle().
			Bold(true).
			Foreground(MarioRed)

	ProtectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(CoinGold)

	DangerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(MarioRed)

	HelpStyle = lipgloss.NewStyle().
			Foreground(TextDim)
)

func repeatChar(ch rune, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(string(ch), n)
}

func ProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int(percent / 100 * float64(width))
	empty := width - filled

	filledStr := repeatChar('▓', filled)
	emptyStr := repeatChar('░', empty)

	var style lipgloss.Style
	switch {
	case percent > 85:
		style = BarCrit
	case percent >= 65:
		style = BarWarn
	default:
		style = BarOK
	}

	return style.Render(filledStr) + HelpStyle.Render(emptyStr)
}
