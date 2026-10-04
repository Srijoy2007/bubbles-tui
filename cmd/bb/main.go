package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Srijoy2007/bubbles-tui/internal/art"
	"github.com/Srijoy2007/bubbles-tui/internal/board"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
	"github.com/Srijoy2007/bubbles-tui/internal/theme"
)

type secondTickMsg struct{}

var noColor = os.Getenv("NO_COLOR") != ""

const splashMinDuration = 3200 * time.Millisecond

type mode int

const (
	modeSplash mode = iota
	modeList
	modeAdd
	modeFocus
)

type loadedMsg struct {
	s   *store.Store
	err error
}

type artTickMsg struct{}

type model struct {
	store         *store.Store
	viewDate      string
	err           error
	mode          mode
	input         textinput.Model
	addErr        string
	cursor        int
	artFrame      int
	splashFrame   int
	splashStarted time.Time	
	height,width int
	focusID     int
	focusEnd    time.Time
	focusTotal  time.Duration
	focusLeft   time.Duration // frozen remaining while paused
	focusPaused bool
	focusDone   bool
}

func secondTicker() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return secondTickMsg{}
	})
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "9-10:30 Deep work"
	ti.CharLimit = 100

	return model{
		viewDate:      time.Now().Format("2006-01-02"),
		input:         ti,
		mode:          modeSplash,
		splashStarted: time.Now(),
	}
}

func artTimer() tea.Cmd {
	return tea.Tick(90*time.Millisecond, func(time.Time) tea.Msg {
		return artTickMsg{}
	})
}

func loadCmd() tea.Msg {
	s, err := store.Load(store.Path())
	return loadedMsg{s: s, err: err}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(loadCmd, artTimer(), secondTicker())
}
func (m model) showCup(nBlocks int) bool {
	if m.width == 0 { // size not known yet
		return true
	}
	return m.width >= 70 && m.height >= 28+max(0, nBlocks-10)
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	
	case secondTickMsg:
		return m, secondTicker()
	case tea.WindowSizeMsg:
		m.width,m.height = msg.Width,msg.Height
		return m,nil

	case loadedMsg:
		m.store = msg.s
		m.err = msg.err

		if time.Since(m.splashStarted) >= splashMinDuration {
			m.finishSplash()
		}

		return m, nil

	case artTickMsg:
		m.artFrame++

		if m.mode == modeSplash {
			m.splashFrame++

			if m.store != nil && time.Since(m.splashStarted) >= splashMinDuration {
				m.finishSplash()
			}
		}
		m.tickFocus()
		return m, artTimer()

	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.mode == modeFocus {
			return m.focusKey(msg.String())
		}

		if m.mode == modeSplash {
			if m.store != nil && (msg.String() == "enter" || msg.String() == "space") {
				m.finishSplash()
			}
			return m, nil
		}

		if m.store == nil {
			return m, nil
		}

		if m.mode == modeAdd {
			switch msg.String() {
			case "esc":
				m.mode = modeList
				m.input.Reset()
				m.addErr = ""
				return m, nil

			case "enter":
				b, err := parseAdd(m.viewDate, m.input.Value())

				if err == nil {
					_, err = m.store.Add(b)
				}

				if err == nil {
					err = m.store.Save(store.Path())
				}

				if err != nil {
					m.addErr = err.Error()
					return m, nil
				}

				m.mode = modeList
				m.input.Reset()
				m.addErr = ""

				return m, nil
			}

			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)

			return m, cmd
		}

		switch msg.String() {
		case "[":
			t, _ := time.Parse("2006-01-02", m.viewDate)
			m.viewDate = t.AddDate(0, 0, -1).Format("2006-01-02")
			m.cursor = 0
			return m, nil

		case "]":
			t, _ := time.Parse("2006-01-02", m.viewDate)
			m.viewDate = t.AddDate(0, 0, 1).Format("2006-01-02")
			m.cursor = 0
			return m, nil

		case "a":
			m.mode = modeAdd
			m.input.Focus()
			return m, textinput.Blink
		case "f":
			m.startFocus()
			return m, nil
		case "j", "down":
			blocks := m.store.On(m.viewDate)

			if m.cursor < len(blocks)-1 {
				m.cursor++
			}

			return m, nil

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

			return m, nil
        case "r":
		nowMin := -1
		if m.viewDate == time.Now().Format("2006-01-02") {
			n := time.Now()
			nowMin = n.Hour()*60 + n.Minute()
		}
		blocks := m.store.On(m.viewDate)
		if m.cursor >= 0 && m.cursor < len(blocks) {
			b := blocks[m.cursor]
			if b.Status(nowMin) == store.StatusMissed {
				t, _ := time.Parse("2006-01-02", m.viewDate)
				tomorrow := t.AddDate(0, 0, 1).Format("2006-01-02")
				if err := m.store.Reschedule(b.ID, tomorrow); err == nil {
					_ = m.store.Save(store.Path())
					if m.cursor > 0 && m.cursor >= len(blocks)-1 {
						m.cursor--
				}
			}
		}
	}
		return m, nil
		case "space":
			blocks := m.store.On(m.viewDate)

			if m.cursor >= 0 && m.cursor < len(blocks) {
				b := blocks[m.cursor]
				blk := m.store.Find(b.ID)

				if blk != nil {
					blk.Done = !blk.Done
					_ = m.store.Save(store.Path())
				}
			}

			return m, nil
		}
	}

	return m, nil
}

func (m *model) finishSplash() {
	if m.store == nil {
		return
	}

	m.mode = modeList
}

var splashFont = map[rune][]string{
	'B': {
		"11110",
		"10001",
		"10001",
		"11110",
		"10001",
		"10001",
		"11110",
	},
	'O': {
		"01110",
		"10001",
		"10001",
		"10001",
		"10001",
		"10001",
		"01110",
	},
	'A': {
		"01110",
		"10001",
		"10001",
		"11111",
		"10001",
		"10001",
		"10001",
	},
	'T': {
		"11111",
		"00100",
		"00100",
		"00100",
		"00100",
		"00100",
		"00100",
	},
	' ': {
		"00000",
		"00000",
		"00000",
		"00000",
		"00000",
		"00000",
		"00000",
	},
}

func renderSplashTitle(frame int) string {
	const scale = 1

	text := "BOBA T"
	var b strings.Builder

	totalWidth := len(text)*6 - 1
	reveal := frame - 2

	if reveal < 0 {
		reveal = 0
	}

	for row := 0; row < 7; row++ {
		b.WriteString(strings.Repeat(" ", 0))

		for i, r := range text {
			glyph := splashFont[r]

			for col := 0; col < len(glyph[row]); col++ {
				visibleAt := i*5 + col

				if visibleAt > reveal {
					b.WriteString(" ")
					continue
				}

				if glyph[row][col] == '1' {
					if noColor {
						b.WriteString("█")
					} else {
						b.WriteString(splashPixelColor(frame, visibleAt))
						b.WriteString(strings.Repeat("█", scale))
						b.WriteString("\033[0m")
					}
				} else {
					b.WriteString(" ")
				}
			}

			if i < len(text)-1 {
				b.WriteByte(' ')
			}
		}

		b.WriteByte('\n')
	}

	_ = totalWidth

	return b.String()
}

func splashPixelColor(frame, x int) string {
	phase := (frame + x) % 16

	if phase < 8 {
		return "\033[38;2;231;151;190m"
	}

	return "\033[38;2;177;164;220m"
}

func splashView(m model) tea.View {
	cup := art.RenderFrame(m.artFrame)

	title := renderSplashTitle(m.splashFrame)

	spinners := []string{"◐", "◓", "◑", "◒"}
	spin := spinners[m.splashFrame%len(spinners)]

	dots := strings.Repeat("·", (m.splashFrame/3)%4)

	status := fmt.Sprintf(
		"%s  brewing your day %s",
		spin,
		dots,
	)

	pct := float64(time.Since(m.splashStarted)) / float64(splashMinDuration)

	if pct > 1 {
		pct = 1
	}

	barWidth := 24
	filled := int(pct * float64(barWidth))

	bar := strings.Repeat("━", filled) +
		strings.Repeat("─", barWidth-filled)

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(centerBlock(cup, 54))
	b.WriteString("\n")
	b.WriteString(centerBlock(title, 54))
	b.WriteString("\n")
	b.WriteString(centerBlock(status, 54))
	b.WriteString("\n")
	b.WriteString(centerBlock(bar, 54))
	b.WriteString("\n\n")
	b.WriteString(centerBlock("small steps · big dreams", 54))
	b.WriteString("\n")
	b.WriteString(centerBlock("enter / space  skip", 54))
	b.WriteString("\n")

	return altScreen(b.String())
}

func (m model) View() tea.View {
	if m.mode == modeSplash {
		return splashView(m)
	}

	if m.err != nil {
		return altScreen(fmt.Sprintf(
			"boba\n\ncould not load today's blocks\n\nerror: %v\n\nq quit\n",
			m.err,
		))
	}

	if m.store == nil {
		return altScreen("loading...\n")
	}
	if m.mode == modeFocus {
		return m.focusView()
	}

	realToday := time.Now().Format("2006-01-02")

	if m.mode == modeAdd {
		out := "add block (" + m.viewDate + ")\n\n" +
			m.input.View() + "\n"

		if m.addErr != "" {
			out += "\n" + theme.Paint(theme.Missed, "", false, m.addErr) + "\n"
		}

		out += "\n" + theme.Paint(theme.Muted, "", false, "enter save · esc cancel") + "\n"

		return altScreen(out)
	}

	blocks := m.store.On(m.viewDate)

	// nowMin is -1 for any day other than today, so past/future days
	// never show a live "current" or "missed" block.
	nowMin := -1
	if m.viewDate == realToday {
		now := time.Now()
		nowMin = now.Hour()*60 + now.Minute()
	}

	// Header + block list.
	var text strings.Builder

	header := theme.Paint(theme.Current, "", true, "boba") +
		theme.Paint(theme.Muted, "", false, " — "+m.viewDate)
	if m.viewDate == realToday {
		header += theme.Paint(theme.Done, "", false, "  (today)")
	}
	text.WriteString(header + "\n")
	text.WriteString(theme.Paint(theme.Border, "", false, strings.Repeat("─", 40)) + "\n\n")

	if len(blocks) == 0 {
		text.WriteString(theme.Paint(theme.Muted, "", false, "nothing planned — press a to add a block") + "\n")
	}

	titleW := 10
	for _, b := range blocks {
		if n := len([]rune(b.Title)); n > titleW {
			titleW = n
		}
	}
	for i, b := range blocks {
		text.WriteString(renderRow(b, i == m.cursor, nowMin, titleW) + "\n")
	}

	top := strings.TrimRight(text.String(), "\n")

	// Cup beside the header when there's room, hidden otherwise.
	if m.showCup(len(blocks)) {
		cup := strings.TrimRight(art.RenderSmallFrame(m.artFrame), "\n")
		top = lipgloss.JoinHorizontal(lipgloss.Center, cup, "   ", top)
	}

	var content strings.Builder
	content.WriteString(top + "\n\n")
	content.WriteString(board.Render(blocks, 5, nowMin))
	content.WriteString("\n\n")
	content.WriteString(theme.Paint(theme.Muted, "", false, "[ / ] day · space done · f focus . a add · j/k move · q quit") + "\n")

	return altScreen(content.String())
}

func rowState(b store.Block, nowMin int) string {
	switch {
	case b.Done:
		return "done"
	case nowMin < 0:
		return "planned"
	case nowMin >= b.Start && nowMin < b.End:
		return "current"
	case nowMin >= b.End:
		return "missed"
	}
	return "planned"
}

func fmtDur(m int) string {
	switch {
	case m < 60:
		return fmt.Sprintf("%dm", m)
	case m%60 == 0:
		return fmt.Sprintf("%dh", m/60)
	}
	return fmt.Sprintf("%dh%02dm", m/60, m%60)
}

func renderRow(b store.Block, selected bool, nowMin, titleW int) string {
	var dot, dotCol, txtCol string
	bold := false
	switch rowState(b, nowMin) {
	case "done":
		dot, dotCol, txtCol = "●", theme.Done, theme.Muted
	case "current":
		dot, dotCol, txtCol, bold = "●", theme.Current, theme.Text, true
	case "missed":
		dot, dotCol, txtCol = "○", theme.Missed, theme.Muted
	default:
		dot, dotCol, txtCol = "○", theme.Planned, theme.Text
	}

	bg := ""
	marker := " "
	if selected {
		bg = theme.Sel
		marker = "▌"
		if theme.NoColor {
			marker = ">"
		}
	}

	sp := theme.Paint("", bg, false, " ")
	timeRange := fmt.Sprintf("%02d:%02d–%02d:%02d", b.Start/60, b.Start%60, b.End/60, b.End%60)

	return theme.Paint(theme.Current, bg, false, marker) + sp +
		theme.Paint(dotCol, bg, false, dot) + sp +
		theme.Paint(theme.Muted, bg, false, timeRange) + sp + sp +
		theme.Paint(txtCol, bg, bold, fmt.Sprintf("%-*s", titleW, b.Title)) + sp + sp +
		theme.Paint(theme.Muted, bg, false, fmt.Sprintf("%-6s", fmtDur(b.End-b.Start))) + sp
}

func centerBlock(s string, width int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")

	var b strings.Builder

	for _, line := range lines {
		visible := ansiVisibleLen(line)

		pad := (width - visible) / 2

		if pad < 0 {
			pad = 0
		}

		b.WriteString(strings.Repeat(" ", pad))
		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String()
}
func ansiVisibleLen(s string) int {
	return lipgloss.Width(s)
}
func altScreen(s string) tea.View {
	v := tea.NewView(s)
	v.AltScreen = true

	return v
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
