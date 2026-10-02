package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/Srijoy2007/bubbles-tui/internal/art"
	"github.com/Srijoy2007/bubbles-tui/internal/board"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

var noColor = os.Getenv("NO_COLOR") != ""

const (
	reverseOn         = "\033[7m"
	reverseOff        = "\033[0m"
	splashMinDuration = 3200 * time.Millisecond
)

func highlight(s string) string {
	if noColor {
		return s
	}
	return reverseOn + s + reverseOff
}

type mode int

const (
	modeSplash mode = iota
	modeList
	modeAdd
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
	return tea.Batch(loadCmd, artTimer())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

		return m, artTimer()

	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
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

	progress := float64(time.Since(m.splashStarted)) / float64(splashMinDuration)

	if progress > 1 {
		progress = 1
	}

	barWidth := 24
	filled := int(progress * float64(barWidth))

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

	realToday := time.Now().Format("2006-01-02")

	if m.mode == modeAdd {
		out := "add block (" + m.viewDate + ")\n\n" +
			m.input.View() + "\n"

		if m.addErr != "" {
			out += "\n" + m.addErr + "\n"
		}

		out += "\nenter save · esc cancel\n"

		return altScreen(out)
	}

	blocks := m.store.On(m.viewDate)

	cup := art.RenderSmallFrame(m.artFrame)

	var content strings.Builder

	content.WriteString(centerBlock(cup, 54))
	content.WriteString("\n")

	header := "boba — " + m.viewDate

	if m.viewDate == realToday {
		header += "  (today)"
	}

	content.WriteString(header + "\n")
	content.WriteString("────────────────────────────────────────\n\n")

	if len(blocks) == 0 {
		content.WriteString("nothing planned — press a to add a block\n")
	} else {
		for i, b := range blocks {
			mark := " "

			if b.Done {
				mark = "x"
			}

			cursor := " "

			if i == m.cursor {
				cursor = ">"
			}

			line := fmt.Sprintf(
				"%s[%s] %02d:%02d-%02d:%02d  %s",
				cursor,
				mark,
				b.Start/60,
				b.Start%60,
				b.End/60,
				b.End%60,
				b.Title,
			)

			if i == m.cursor {
				line = highlight(line)
			}

			content.WriteString(line + "\n")
		}
	}

	now := time.Now()

	content.WriteString("\n")
	content.WriteString(
		board.Render(
			blocks,
			5,
			now.Hour()*60+now.Minute(),
		),
	)
	content.WriteString("\n")
	content.WriteString("\n[ / ] day · space done · a add · j/k move · q quit\n")

	return altScreen(content.String())
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
	n := 0
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}

		if inEscape {
			if s[i] >= '@' && s[i] <= '~' {
				inEscape = false
			}
			continue
		}

		n++
	}

	return n
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
