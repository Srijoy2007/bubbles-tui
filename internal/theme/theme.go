package theme

import (
	"os"

	"charm.land/lipgloss/v2"
)

var NoColor = os.Getenv("NO_COLOR") != ""

const (
	Done    = "#B1A4DC"
	Missed  = "#D27D9B"
	Planned = "#9694AF"
	Empty   = "#4B4C5C"
	Current = "#E797BE"
	Border  = "#69697D"
	Text    = "#BEBCCD"
	Muted   = "#78768C"
	Sel     = "#2F2C4A" // selected-row background
)

// Paint colors s. Empty fg/bg means "leave unset".
func Paint(fg, bg string, bold bool, s string) string {
	if NoColor {
		return s
	}
	st := lipgloss.NewStyle().Bold(bold)
	if fg != "" {
		st = st.Foreground(lipgloss.Color(fg))
	}
	if bg != "" {
		st = st.Background(lipgloss.Color(bg))
	}
	return st.Render(s)
}
