package board

import (
	"os"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

var noColor = os.Getenv("NO_COLOR") != ""

const (
	// Boba theme — soft pink / lavender.
	colDone    = "#B1A4DC" // lavender
	colMissed  = "#D27D9B" // muted rose
	colPlanned = "#9694AF" // muted lavender-gray
	colEmpty   = "#4B4C5C" // dark gray
	colCurrent = "#E797BE" // boba pink
	colBorder  = "#69697D" // soft border
	colText    = "#BEBCCD" // soft text
	colMuted   = "#78768C" // dim label
)

const minutesDay = 1440

func colorize(hex, s string) string {
	if noColor {
		return s
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render(s)
}

func state(b store.Block, now int) string {
	switch b.Status(now) {
	case store.StatusDone:
		return "done"
	case store.StatusCurrent:
		return "current"
	case store.StatusMissed:
		return "missed"
	default:
		return "planned"
	}
}

func glyph(s string) string {
	switch s {
	case "done":
		return colorize(colDone, "━")
	case "current":
		return colorize(colCurrent, "━")
	case "missed":
		return colorize(colMissed, "━")
	case "planned":
		return colorize(colPlanned, "┄")
	default:
		return colorize(colEmpty, "·")
	}
}

// hourAxis places each two-character hour label at the exact column it
// falls on, so it always lines up with the timeline below it, at any width.
func hourAxis(width int) string {
	buf := []rune(strings.Repeat(" ", width))
	for hour := 0; hour < 24; hour += 2 {
		pos := hour * 60 * width / minutesDay
		for i, r := range twoDigits(hour) {
			if pos+i < width {
				buf[pos+i] = r
			}
		}
	}
	return string(buf)
}

// Render draws a single card: an hour-labeled timeline, a now-marker,
// a totals summary, and one compact pearl-progress row per block.
// The resolution parameter from earlier callers is kept for compatibility
// but no longer used — the timeline width itself sets the granularity.
func Render(blocks []store.Block, _ int, now int) string {
	const (
		width = 72 // divides 1440 evenly: 20 minutes per column
		pad   = "  "
	)
	cellMinutes := minutesDay / width

	var out strings.Builder


	out.WriteString(colorize(colBorder, strings.Repeat("─", width+len(pad))))
	out.WriteString("\n")

	out.WriteString(pad)
	out.WriteString(colorize(colMuted, hourAxis(width)))
	out.WriteString("\n")

	out.WriteString(pad)
	for i := 0; i < width; i++ {
		minute := i * cellMinutes
		best := "empty"
		for _, b := range blocks {
			if minute < b.Start || minute >= b.End {
				continue
			}
			s := state(b, now)
			if best == "empty" || s == "current" || s == "done" {
				best = s
			}
		}
		out.WriteString(glyph(best))
	}
	out.WriteString("\n")

	if now >= 0 {
		pos := now * width / minutesDay
		if pos >= 0 && pos < width {
			out.WriteString(pad)
			out.WriteString(strings.Repeat(" ", pos))
			out.WriteString(colorize(colCurrent, "▲ now"))
			out.WriteString("\n")
		}
	}

	if len(blocks) > 0 {
		out.WriteString("\n")
		out.WriteString(pad)
		out.WriteString(summaryLine(blocks, now))
		out.WriteString("\n\n")
		for _, b := range blocks {
			out.WriteString(pad)
			out.WriteString(blockRow(b, now))
			out.WriteString("\n")
		}
	}

	out.WriteString("\n")
	out.WriteString(pad)
	out.WriteString(colorize(colDone, "●") + " done    ")
	out.WriteString(colorize(colCurrent, "●") + " current    ")
	out.WriteString(colorize(colPlanned, "●") + " planned    ")
	out.WriteString(colorize(colMissed, "●") + " missed")

	return out.String()
}

// summaryLine totals done/missed/planned time across all of today's blocks.
func summaryLine(blocks []store.Block, now int) string {
	plannedMin, doneMin, missedMin := 0, 0, 0
	for _, b := range blocks {
		d := b.End - b.Start
		plannedMin += d
		switch state(b, now) {
		case "done":
			doneMin += d
		case "missed":
			missedMin += d
		}
	}
	return colorize(colDone, "done "+fmtDur(doneMin)) + "   " +
		colorize(colMissed, "missed "+fmtDur(missedMin)) + "   " +
		colorize(colMuted, "of "+fmtDur(plannedMin)+" planned")
}

// blockRow renders one line: a status dot, the time range, the title,
// and a small pearl bar showing progress through the block.
func blockRow(b store.Block, now int) string {
	const barW = 10

	s := state(b, now)
	var dot, title string
	switch s {
	case "done", "current":
		col := colDone
		if s == "current" {
			col = colCurrent
		}
		dot, title = colorize(col, "●"), colorize(colText, b.Title)
	default:
		col := colPlanned
		if s == "missed" {
			col = colMissed
		}
		dot, title = colorize(col, "○"), colorize(colMuted, b.Title)
	}

	frac := 0.0
	switch {
	case b.Done, now >= b.End:
		frac = 1
	case now >= b.Start && now < b.End:
		frac = float64(now-b.Start) / float64(b.End-b.Start)
	}
	if frac < 0 {
		frac = 0
	}
	if frac > 1 {
		frac = 1
	}
	filled := int(frac*barW + 0.5)

	barCol := colPlanned
	switch s {
	case "done":
		barCol = colDone
	case "current":
		barCol = colCurrent
	case "missed":
		barCol = colMissed
	}
	bar := colorize(barCol, strings.Repeat("●", filled)) +
		colorize(colEmpty, strings.Repeat("○", barW-filled))

	return dot + " " + colorize(colMuted, formatTime(b.Start)+"–"+formatTime(b.End)) +
		"  " + title + "  " + bar
}

func formatTime(minutes int) string {
	return twoDigits(minutes/60) + ":" + twoDigits(minutes%60)
}

func fmtDur(m int) string {
	if m < 60 {
		return itoa(m) + "m"
	}
	h, mm := m/60, m%60
	if mm == 0 {
		return itoa(h) + "h"
	}
	return itoa(h) + "h" + twoDigits(mm) + "m"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func twoDigits(n int) string {
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}
