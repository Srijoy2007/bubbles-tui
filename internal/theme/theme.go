package theme

import (
	"os"
	"strings"
	"charm.land/lipgloss/v2"
)

var NoColor = os.Getenv("NO_COLOR") != ""

var (
	Done    = "#B1A4DC"
	Missed  = "#D27D9B"
	Planned = "#9694AF"
	Empty   = "#4B4C5C"
	Current = "#E797BE"
	Border  = "#69697D"
	Text    = "#BEBCCD"
	Muted   = "#78768C"
	Sel     = "#2F2C4A"
	BG      = "" )

type palette struct {
	Done, Missed, Planned, Empty, Current, Border, Text, Muted, Sel, BG string
}

var palettes = map[string]palette{
	"default": {
		Done: "#B1A4DC", Missed: "#D27D9B", Planned: "#9694AF", Empty: "#4B4C5C",
		Current: "#E797BE", Border: "#69697D", Text: "#BEBCCD", Muted: "#78768C", Sel: "#2F2C4A", BG: "",
	},
	"party": {
		Done: "#39FF14", Missed: "#FF5F56", Planned: "#00E5FF", Empty: "#2E003E",
		Current: "#FF2BD6", Border: "#FF8C00", Text: "#FFFFFF", Muted: "#FFD23F", Sel: "#4B1D6B", BG: "",
	},
	"catppuccin": { 
		Done: "#A6E3A1", Missed: "#F38BA8", Planned: "#89B4FA", Empty: "#313244",
		Current: "#F5C2E7", Border: "#585B70", Text: "#CDD6F4", Muted: "#9399B2", Sel: "#45475A", BG: "#1E1E2E",
	},
	"gruvbox": { 		
		Done: "#B8BB26", Missed: "#FB4934", Planned: "#83A598", Empty: "#3C3836",
		Current: "#FABD2F", Border: "#665C54", Text: "#EBDBB2", Muted: "#A89984", Sel: "#504945", BG: "#282828",
	},
	"tokyonight": { 
		Done: "#9ECE6A", Missed: "#F7768E", Planned: "#7AA2F7", Empty: "#24283B",
		Current: "#BB9AF7", Border: "#414868", Text: "#C0CAF5", Muted: "#565F89", Sel: "#292E42", BG: "#1A1B26",
	},
}

func Names() []string {
	return []string{"default", "party", "catppuccin", "gruvbox", "tokyonight"}
}
func SetTheme(name string) {
	p, ok := palettes[name]
	if !ok {
		p = palettes["default"]
	}
	Done, Missed, Planned, Empty = p.Done, p.Missed, p.Planned, p.Empty
	Current, Border, Text, Muted, Sel = p.Current, p.Border, p.Text, p.Muted, p.Sel
	BG = p.BG
}
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

// Box draws a rounded, theme-colored border around content, with an
// optional title embedded in the top edge.
func Box(title, content string) string {
	lines := splitLines(content)
	inner := 0
	for _, l := range lines {
		if w := lipgloss.Width(l); w > inner {
			inner = w
		}
	}
	const padX = 2
	inner += padX * 2

	bc := func(s string) string { return Paint(Border, "", false, s) }

	var sb strings.Builder
	if title != "" {
		label := Paint(Current, "", true, " "+title+" ")
		dashes := inner - lipgloss.Width(label)
		if dashes < 0 {
			dashes = 0
		}
		sb.WriteString(bc("╭─") + label + bc(strings.Repeat("─", dashes)+"╮") + "\n")
	} else {
		sb.WriteString(bc("╭"+strings.Repeat("─", inner)+"╮") + "\n")
	}

	for _, l := range lines {
		pad := inner - padX*2 - lipgloss.Width(l)
		if pad < 0 {
			pad = 0
		}
		sb.WriteString(bc("│") + strings.Repeat(" ", padX) + l +
			strings.Repeat(" ", pad) + strings.Repeat(" ", padX) + bc("│") + "\n")
	}

	sb.WriteString(bc("╰" + strings.Repeat("─", inner) + "╯"))
	return sb.String()
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

func joinLines(lines []string) string {
	out := ""
	for i, l := range lines {
		if i > 0 {
			out += "\n"
		}
		out += l
	}
	return out
}

func overlayLabel(topBorder, label string) string {
	const insertAt = 2
	runes := []rune(topBorder)
	if len(runes) <= insertAt {
		return topBorder
	}
	prefix := string(runes[:insertAt])
	return prefix + label
}
