package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/Srijoy2007/bubbles-tui/internal/art"
	"github.com/Srijoy2007/bubbles-tui/internal/audio"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
	"github.com/Srijoy2007/bubbles-tui/internal/theme"
)

const (
	dimPink = "#C77FA3" // scanline shade of Current
	dimRose = "#B3687F" // scanline shade of Missed
)

var digitFont = map[rune][]string{
	'0': {"01110", "10001", "10011", "10101", "11001", "10001", "01110"},
	'1': {"00100", "01100", "00100", "00100", "00100", "00100", "01110"},
	'2': {"01110", "10001", "00001", "00010", "00100", "01000", "11111"},
	'3': {"11110", "00001", "00001", "01110", "00001", "00001", "11110"},
	'4': {"00010", "00110", "01010", "10010", "11111", "00010", "00010"},
	'5': {"11111", "10000", "11110", "00001", "00001", "10001", "01110"},
	'6': {"00110", "01000", "10000", "11110", "10001", "10001", "01110"},
	'7': {"11111", "00001", "00010", "00100", "01000", "01000", "01000"},
	'8': {"01110", "10001", "10001", "01110", "10001", "10001", "01110"},
	'9': {"01110", "10001", "10001", "01111", "00001", "00010", "01100"},
	':': {"0", "0", "1", "0", "1", "0", "0"},
}

// ---- state transitions ----

func (m *model) startFocus() {
	if m.store == nil || m.viewDate != time.Now().Format("2006-01-02") {
		return
	}
	blocks := m.store.On(m.viewDate)
	if m.cursor < 0 || m.cursor >= len(blocks) {
		return
	}
	b := blocks[m.cursor]
	if b.Done {
		return
	}

	now := time.Now()
	total := time.Duration(b.End-b.Start) * time.Minute
	dur := total
	nowMin := now.Hour()*60 + now.Minute()
	if nowMin >= b.Start && nowMin < b.End {
		end := time.Date(now.Year(), now.Month(), now.Day(), b.End/60, b.End%60, 0, 0, now.Location())
		dur = end.Sub(now)
	}

	m.focusID = b.ID
	m.focusEnd = now.Add(dur)
	m.focusTotal = total
	m.focusPaused = false
	m.focusDone = false
	m.mode = modeFocus
}

func (m model) focusRemaining() time.Duration {
	if m.focusPaused {
		return m.focusLeft
	}
	return max(0, time.Until(m.focusEnd))
}

func (m *model) tickFocus() {
	if m.mode != modeFocus || m.focusDone || m.focusPaused {
		return
	}
	if m.focusRemaining() <= 0 {
		m.completeFocus()
	}
}

func (m *model) completeFocus() {
	if blk := m.store.Find(m.focusID); blk != nil {
		blk.Done = true
		_ = m.store.Save(store.Path())
	}
	m.focusDone = true
}

func (m model) focusKey(k string) (tea.Model, tea.Cmd) {
	if k == "ctrl+c" {
		return m, tea.Quit
	}
	if m.focusDone { // any key leaves the "complete" screen
		m.audio.Stop()
		m.mode = modeList
		return m, nil
	}
	switch k {
	case "space":
		if m.focusPaused {
			m.focusEnd = time.Now().Add(m.focusLeft)
			m.focusPaused = false
		} else {
			m.focusLeft = m.focusRemaining()
			m.focusPaused = true
		}
	case "+", "=":
		if m.focusPaused {
			m.focusLeft += 5 * time.Minute
		} else {
			m.focusEnd = m.focusEnd.Add(5 * time.Minute)
		}
		m.focusTotal += 5 * time.Minute
	case "d":
		m.completeFocus()
	case "m":
		m.audioLevel = (m.audioLevel + 1) % len(audio.Tracks)
		if m.audioLevel == 0 {
			m.audio.Stop()
		} else if err := m.audio.Start(audio.Tracks[m.audioLevel].URL); err != nil {
			m.audioLevel = 0 // fall back to "off" silently rather than crash
		}
	case "esc", "q":
		m.audio.Stop()
		m.mode = modeList
	}
	return m, nil
}

// ---- rendering ----

func clockText(d time.Duration) string {
	secs := int((d + time.Second - 1) / time.Second) // round up
	h, mn, s := secs/3600, (secs%3600)/60, secs%60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, mn, s)
	}
	return fmt.Sprintf("%02d:%02d", mn, s)
}

func bigClock(text, c1, c2 string) []string {
	rows := make([]string, 7)
	for r := 0; r < 7; r++ {
		col := c1
		if r%2 == 1 {
			col = c2
		}
		lit := theme.Paint(col, "", false, "██")
		off := theme.Paint(theme.Empty, "", false, "░░")

		var sb strings.Builder
		for i, ch := range text {
			for _, px := range digitFont[ch][r] {
				if px == '1' {
					sb.WriteString(lit)
				} else {
					sb.WriteString(off)
				}
			}
			if i < len(text)-1 {
				sb.WriteString("  ")
			}
		}
		rows[r] = sb.String()
	}
	return rows
}

func retroBox(title string, rows []string) string {
	inner := 0
	for _, r := range rows {
		inner = max(inner, ansiVisibleLen(r))
	}
	inner += 6

	bc := func(s string) string { return theme.Paint(theme.Border, "", false, s) }
	label := " " + title + " "

	var sb strings.Builder
	sb.WriteString(bc("╔═") + theme.Paint(theme.Current, "", true, label) +
		bc(strings.Repeat("═", max(0, inner-1-len([]rune(label))))+"╗") + "\n")
	for _, r := range rows {
		pad := inner - 3 - ansiVisibleLen(r)
		sb.WriteString(bc("║") + "   " + r + strings.Repeat(" ", max(0, pad)) + bc("║") + "\n")
	}
	sb.WriteString(bc("╚" + strings.Repeat("═", inner) + "╝"))
	return sb.String()
}

func retroBar(frac float64, width int) string {
	frac = max(0, min(1, frac))
	filled := int(frac*float64(width) + 0.5)
	return theme.Paint(theme.Current, "", false, strings.Repeat("█", filled)) +
		theme.Paint(theme.Empty, "", false, strings.Repeat("░", width-filled))
}

var barLevels = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

// musicViz draws a row of bars that ripple smoothly over time, each bar
// phase-offset from its neighbor so they move like a real equalizer
// rather than all pulsing in lockstep. Driven by the same artFrame
// ticker that already runs every 90ms — no extra timer needed.
func musicViz(frame, bars int) string {
	var sb strings.Builder
	for i := 0; i < bars; i++ {
		phase := float64(frame)*0.25 + float64(i)*0.9
		h := (math.Sin(phase) + 1) / 2 // normalize to 0..1
		lvl := int(h * float64(len(barLevels)-1))
		sb.WriteString(theme.Paint(theme.Current, "", false, barLevels[lvl]))
	}
	return sb.String()
}

func (m model) focusView() tea.View {
	w := m.width
	if w <= 0 {
		w = 80
	}

	title, rng := "focus", ""
	if blk := m.store.Find(m.focusID); blk != nil {
		title = blk.Title
		rng = fmt.Sprintf("%02d:%02d–%02d:%02d", blk.Start/60, blk.Start%60, blk.End/60, blk.End%60)
	}

	var body string
	if m.focusDone {
		cup := strings.TrimRight(art.RenderSmallFrame(m.artFrame), "\n")
		body = cup + "\n\n" +
			theme.Paint(theme.Current, "", true, "BLOCK COMPLETE") + "\n" +
			theme.Paint(theme.Muted, "", false, title) + "\n\n" +
			theme.Paint(theme.Muted, "", false, "press any key")
	} else {
		rem := m.focusRemaining()
		c1, c2 := theme.Current, dimPink
		switch {
		case m.focusPaused:
			c1, c2 = theme.Planned, theme.Muted
		case rem <= time.Minute:
			c1, c2 = theme.Missed, dimRose
		}

		text := clockText(rem)
		rows := bigClock(text, c1, c2)
		if m.width > 0 && ansiVisibleLen(rows[0])+8 > m.width { // narrow terminal
			rows = []string{theme.Paint(c1, "", true, text)}
		}

		frac := 0.0
		if m.focusTotal > 0 {
			frac = 1 - float64(rem)/float64(m.focusTotal)
		}

		status := theme.Paint(theme.Done, "", false, "▶ running")
		if m.focusPaused {
			status = strings.Repeat(" ", 8)
			if (m.artFrame/6)%2 == 0 {
				status = theme.Paint(theme.Missed, "", true, "‖ paused")
			}
		}

		lines := []string{
			retroBox("FOCUS", rows),
			"",
			theme.Paint(theme.Text, "", true, title) + "  " + theme.Paint(theme.Muted, "", false, rng),
			"",
			retroBar(frac, 40) + theme.Paint(theme.Muted, "", false, fmt.Sprintf("  %3d%%", int(frac*100))),
			"",
			status,
		}

		if m.audioLevel > 0 {
			trackName := theme.Paint(theme.Muted, "", false, "♪ "+audio.Tracks[m.audioLevel].Name+"  ")
			viz := trackName
			if !m.focusPaused {
				viz += musicViz(m.artFrame, 16)
			} else {
				viz += theme.Paint(theme.Muted, "", false, strings.Repeat("▁", 16))
			}
			lines = append(lines, "", viz)
		}

		lines = append(lines, "",
			theme.Paint(theme.Muted, "", false,
				"space pause · + 5min · d done · m music: "+audio.Tracks[m.audioLevel].Name+" · esc leave"),
		)

		body = strings.Join(lines, "\n")
	}

	top := 0
	if m.height > 0 {
		top = max(0, (m.height-(strings.Count(body, "\n")+1))/2)
	}
	return altScreen(strings.Repeat("\n", top) + centerBlock(body, w))
}
