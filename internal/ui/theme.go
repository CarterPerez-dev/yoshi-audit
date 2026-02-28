// ©AngelaMos | 2026
// theme.go

package ui

import (
	"github.com/CarterPerez-dev/yoshi-audit/internal/ui/theme"
)

var (
	YoshiGreen = theme.YoshiGreen
	PipeGreen  = theme.PipeGreen
	CoinGold   = theme.CoinGold
	MarioRed   = theme.MarioRed
	OneUpGreen = theme.OneUpGreen
	SnesBg     = theme.SnesBg
	TextWhite  = theme.TextWhite
	TextDim    = theme.TextDim
	BrickBrown = theme.BrickBrown
)

var MarioBorder = theme.MarioBorder

var (
	FrameStyle       = theme.FrameStyle
	TitleStyle       = theme.TitleStyle
	ActiveTabStyle   = theme.ActiveTabStyle
	InactiveTabStyle = theme.InactiveTabStyle
	BarOK            = theme.BarOK
	BarWarn          = theme.BarWarn
	BarCrit          = theme.BarCrit
	StatusOK         = theme.StatusOK
	StatusWarn       = theme.StatusWarn
	StatusCrit       = theme.StatusCrit
	ProtectedStyle   = theme.ProtectedStyle
	DangerStyle      = theme.DangerStyle
	HelpStyle        = theme.HelpStyle
)

func ProgressBar(percent float64, width int) string {
	return theme.ProgressBar(percent, width)
}
