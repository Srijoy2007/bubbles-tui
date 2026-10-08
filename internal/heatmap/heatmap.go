package heatmap

import (
	"strings"
	"time"

	"github.com/Srijoy2007/bubbles-tui/internal/theme"
)

var dayLabels = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}


func Render(minutesByDate map[string]int, weeks int) string {
	today := time.Now()
	thisWeekSunday := today.AddDate(0, 0, -int(today.Weekday()))
	start := thisWeekSunday.AddDate(0, 0, -7*(weeks-1))

	max := 0
	for d := start; !d.After(today); d = d.AddDate(0, 0, 1) {
		if m := minutesByDate[d.Format("2006-01-02")]; m > max {
			max = m
		}
	}
	if max == 0 {
		max = 1
	}

	var sb strings.Builder
	for row := 0; row < 7; row++ {
		sb.WriteString(theme.Paint(theme.Muted, "", false, padLabel(dayLabels[row])))
		for w := 0; w < weeks; w++ {
			d := start.AddDate(0, 0, w*7+row)
			if d.After(today) {
				sb.WriteString("  ")
				continue
			}
			mins := minutesByDate[d.Format("2006-01-02")]
			sb.WriteString(levelGlyph(levelFor(mins, max)))
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n" + theme.Paint(theme.Muted, "", false, "less "))
	for i := 0; i <= 4; i++ {
		sb.WriteString(levelGlyph(i) + " ")
	}
	sb.WriteString(theme.Paint(theme.Muted, "", false, "more"))

	return sb.String()
}

func padLabel(s string) string {
	for len([]rune(s)) < 4 {
		s += " "
	}
	return s
}

func levelFor(mins, max int) int {
	if mins <= 0 {
		return 0
	}
	frac := float64(mins) / float64(max)
	switch {
	case frac < 0.25:
		return 1
	case frac < 0.5:
		return 2
	case frac < 0.75:
		return 3
	default:
		return 4
	}
}

func levelGlyph(level int) string {
	switch level {
	case 0:
		return theme.Paint(theme.Empty, "", false, "■")
	case 1:
		return theme.Paint("#4A3F63", "", false, "■")
	case 2:
		return theme.Paint("#6B5590", "", false, "■")
	case 3:
		return theme.Paint("#9B7FC9", "", false, "■")
	default:
		return theme.Paint(theme.Current, "", false, "■")
	}
}
