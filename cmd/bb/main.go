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
	"github.com/Srijoy2007/bubbles-tui/internal/audio"
	"github.com/Srijoy2007/bubbles-tui/internal/board"
	"github.com/Srijoy2007/bubbles-tui/internal/heatmap"
	"github.com/Srijoy2007/bubbles-tui/internal/prefs"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
	"github.com/Srijoy2007/bubbles-tui/internal/theme"
)

type secondTickMsg struct{}
type artTickMsg struct{}

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
	height        int
	width         int

	focusID     int
	focusEnd    time.Time
	focusTotal  time.Duration
	focusLeft   time.Duration
	focusPaused bool
	focusDone   bool

	mascot string

	audio      audio.Player
	audioLevel int
}

func secondTicker() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return secondTickMsg{}
	})
}

func artTimer() tea.Cmd {
	return tea.Tick(90*time.Millisecond, func(time.Time) tea.Msg {
		return artTickMsg{}
	})
}

func newModel(p prefs.Prefs) model {
	ti := textinput.New()
	ti.Placeholder = "9-10:30 Deep work"
	ti.CharLimit = 100

	return model{
		viewDate:      time.Now().Format("2006-01-02"),
		input:         ti,
		mode:          modeSplash,
		splashStarted: time.Now(),
		mascot:        p.Mascot,
	}
}

func loadCmd() tea.Msg {
	s, err := store.Load(store.Path())
	return loadedMsg{
		s:   s,
		err: err,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		loadCmd,
		artTimer(),
		secondTicker(),
	)
}

func (m model) showCup(nBlocks int) bool {
	if m.width == 0 {
		return true
	}

	return m.width >= 70 &&
		m.height >= 28+max(0, nBlocks-10)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case secondTickMsg:
		return m, secondTicker()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

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

			if m.store != nil &&
				time.Since(m.splashStarted) >= splashMinDuration {
				m.finishSplash()
			}
		}

		m.tickFocus()

		return m, artTimer()

	case tea.KeyPressMsg:
		key := msg.String()

		if key == "q" || key == "ctrl+c" {
			m.audio.Stop()
			return m, tea.Quit
		}

		if m.mode == modeFocus {
			return m.focusKey(key)
		}

		if m.mode == modeSplash {
			if m.store != nil &&
				(key == "enter" || key == "space") {
				m.finishSplash()
			}
			return m, nil
		}

		if m.store == nil {
			return m, nil
		}

		if m.mode == modeAdd {
			switch key {

			case "esc":
				m.mode = modeList
				m.input.Reset()
				m.addErr = ""
				return m, nil

			case "enter":
				b, err := parseAdd(
					m.viewDate,
					m.input.Value(),
				)

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

		switch key {

		case "[":
			t, _ := time.Parse(
				"2006-01-02",
				m.viewDate,
			)

			m.viewDate = t.
				AddDate(0, 0, -1).
				Format("2006-01-02")

			m.cursor = 0
			return m, nil

		case "]":
			t, _ := time.Parse(
				"2006-01-02",
				m.viewDate,
			)

			m.viewDate = t.
				AddDate(0, 0, 1).
				Format("2006-01-02")

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
				now := time.Now()
				nowMin = now.Hour()*60 + now.Minute()
			}

			blocks := m.store.On(m.viewDate)

			if m.cursor >= 0 && m.cursor < len(blocks) {
				b := blocks[m.cursor]

				if b.Status(nowMin) == store.StatusMissed {
					t, _ := time.Parse(
						"2006-01-02",
						m.viewDate,
					)

					tomorrow := t.
						AddDate(0, 0, 1).
						Format("2006-01-02")

					if err := m.store.Reschedule(
						b.ID,
						tomorrow,
					); err == nil {
						_ = m.store.Save(store.Path())

						if m.cursor > 0 &&
							m.cursor >= len(blocks)-1 {
							m.cursor--
						}
					}
				}
			}

			return m, nil

		case "space":
			blocks := m.store.On(m.viewDate)

			if m.cursor >= 0 &&
				m.cursor < len(blocks) {

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

// ============================================================
// SPLASH
// ============================================================

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

	reveal := frame - 2

	if reveal < 0 {
		reveal = 0
	}

	for row := 0; row < 7; row++ {
		for i, r := range text {
			glyph := splashFont[r]

			for col := 0; col < len(glyph[row]); col++ {
				visibleAt := i*5 + col

				if visibleAt > reveal {
					b.WriteByte(' ')
					continue
				}

				if glyph[row][col] == '1' {
					if noColor {
						b.WriteString("█")
					} else {
						b.WriteString(
							splashPixelColor(
								frame,
								visibleAt,
							),
						)

						b.WriteString(
							strings.Repeat("█", scale),
						)

						b.WriteString("\033[0m")
					}
				} else {
					b.WriteByte(' ')
				}
			}

			if i < len(text)-1 {
				b.WriteByte(' ')
			}
		}

		b.WriteByte('\n')
	}

	return b.String()
}

func splashPixelColor(frame, x int) string {
	phase := (frame + x) % 16

	if phase < 8 {
		return "\033[38;2;231;151;190m"
	}

	return "\033[38;2;177;164;220m"
}

// ============================================================
// MASCOTS
// ============================================================

func (m model) cupArt() string {
	switch m.mascot {
	case "bear":
		return art.BearSmallFrame(m.artFrame)

	case "cat":
		return art.CatSmallFrame(m.artFrame)

	default:
		return art.RenderSmallFrame(m.artFrame)
	}
}

func (m model) cupArtBig() string {
	switch m.mascot {
	case "bear":
		return art.BearFrame(m.artFrame)

	case "cat":
		return art.CatFrame(m.artFrame)

	default:
		return art.RenderFrame(m.artFrame)
	}
}

func splashView(m model) tea.View {
	mascot := m.cupArtBig()
	title := renderSplashTitle(m.splashFrame)

	spinners := []string{
		"◐",
		"◓",
		"◑",
		"◒",
	}

	spin := spinners[m.splashFrame%len(spinners)]

	dots := strings.Repeat(
		"·",
		(m.splashFrame/3)%4,
	)

	status := fmt.Sprintf(
		"%s  brewing your day %s",
		spin,
		dots,
	)

	pct := float64(
		time.Since(m.splashStarted),
	) / float64(splashMinDuration)

	if pct > 1 {
		pct = 1
	}

	const barWidth = 24

	filled := int(
		pct * float64(barWidth),
	)

	bar := strings.Repeat(
		"━",
		filled,
	) + strings.Repeat(
		"─",
		barWidth-filled,
	)

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(centerBlock(mascot, 54))
	b.WriteString("\n")
	b.WriteString(centerBlock(title, 54))
	b.WriteString("\n")
	b.WriteString(centerBlock(status, 54))
	b.WriteString("\n")
	b.WriteString(centerBlock(bar, 54))
	b.WriteString("\n\n")
	b.WriteString(
		centerBlock(
			"small steps · big dreams",
			54,
		),
	)
	b.WriteString("\n")
	b.WriteString(
		centerBlock(
			"enter / space  skip",
			54,
		),
	)
	b.WriteString("\n")

	return altScreen(b.String())
}

// ============================================================
// MAIN VIEW
// ============================================================

func (m model) View() tea.View {
	if m.mode == modeSplash {
		return splashView(m)
	}

	if m.err != nil {
		return altScreen(
			fmt.Sprintf(
				"boba\n\n"+
					"could not load today's blocks\n\n"+
					"error: %v\n\n"+
					"q quit\n",
				m.err,
			),
		)
	}

	if m.store == nil {
		return altScreen("loading...\n")
	}

	if m.mode == modeFocus {
		return m.focusView()
	}

	realToday := time.Now().Format("2006-01-02")

	if m.mode == modeAdd {
		out := "add block (" +
			m.viewDate +
			")\n\n" +
			m.input.View() +
			"\n"

		if m.addErr != "" {
			out += "\n" +
				theme.Paint(
					theme.Missed,
					"",
					false,
					m.addErr,
				) +
				"\n"
		}

		out += "\n" +
			theme.Paint(
				theme.Muted,
				"",
				false,
				"enter save · esc cancel",
			) +
			"\n"

		return altScreen(out)
	}

	blocks := m.store.On(m.viewDate)

	nowMin := -1

	if m.viewDate == realToday {
		now := time.Now()
		nowMin = now.Hour()*60 + now.Minute()
	}

	// ========================================================
	// LEFT: SCHEDULE
	// ========================================================

	var left strings.Builder

	header :=
		theme.Paint(
			theme.Current,
			"",
			true,
			"boba",
		) +
			theme.Paint(
				theme.Muted,
				"",
				false,
				"  "+m.viewDate,
			)

	if m.viewDate == realToday {
		header += theme.Paint(
			theme.Done,
			"",
			false,
			"  · today",
		)
	}

	left.WriteString(header)
	left.WriteString("\n\n")

	if len(blocks) == 0 {
		left.WriteString(
			theme.Paint(
				theme.Muted,
				"",
				false,
				"nothing planned — press a to add a block",
			),
		)
		left.WriteByte('\n')
	} else {
		titleW := 10

		for _, b := range blocks {
			if n := len([]rune(b.Title)); n > titleW {
				titleW = n
			}
		}

		for i, b := range blocks {
			left.WriteString(
				renderRow(
					b,
					i == m.cursor,
					nowMin,
					titleW,
				),
			)
			left.WriteByte('\n')
		}
	}

	leftPanel := theme.Box(
		"",
		strings.TrimRight(
			left.String(),
			"\n",
		),
	)

	if m.showCup(len(blocks)) {
		mascot := strings.TrimRight(
			m.cupArt(),
			"\n",
		)

		leftPanel = lipgloss.JoinVertical(
			lipgloss.Left,
			mascot,
			"",
			leftPanel,
		)
	}

	// ========================================================
	// RIGHT: STATS
	// ========================================================

	doneMin := 0
	missedMin := 0
	plannedMin := 0

	for _, b := range blocks {
		duration := b.End - b.Start
		plannedMin += duration

		switch b.Status(nowMin) {
		case store.StatusDone:
			doneMin += duration

		case store.StatusMissed:
			missedMin += duration
		}
	}

	var stats strings.Builder

	stats.WriteString(
		theme.Paint(
			theme.Text,
			"",
			true,
			"today",
		),
	)
	stats.WriteString("\n\n")

	stats.WriteString(
		theme.Paint(
			theme.Muted,
			"",
			false,
			"done     ",
		),
	)

	stats.WriteString(
		theme.Paint(
			theme.Done,
			"",
			true,
			fmtDur(doneMin),
		),
	)

	stats.WriteByte('\n')

	stats.WriteString(
		theme.Paint(
			theme.Muted,
			"",
			false,
			"missed   ",
		),
	)

	stats.WriteString(
		theme.Paint(
			theme.Missed,
			"",
			true,
			fmtDur(missedMin),
		),
	)

	stats.WriteByte('\n')

	stats.WriteString(
		theme.Paint(
			theme.Muted,
			"",
			false,
			"planned  ",
		),
	)

	stats.WriteString(
		theme.Paint(
			theme.Text,
			"",
			false,
			fmtDur(plannedMin),
		),
	)

	rightCol := theme.Box(
		"STATS",
		strings.TrimRight(
			stats.String(),
			"\n",
		),
	)

	// ========================================================
	// DAY + RHYTHM
	//
	// These intentionally live inside ONE box.
	// ========================================================

	dayTrack := board.Render(
		blocks,
		5,
		nowMin,
	)

	dayContent := strings.TrimRight(
		dayTrack,
		"\n",
	)

	// Heatmap stays directly underneath DAY.
	if m.width == 0 || m.width >= 90 {
		rhythm := heatmap.Render(
			m.store.DoneMinutes(),
			12,
		)

		dayContent += "\n\n"

		dayContent += theme.Paint(
			theme.Current,
			"",
			true,
			"RHYTHM · 12 WEEKS",
		)

		dayContent += "\n"

		dayContent += rhythm
	}

	dayPanel := theme.Box(
		"DAY",
		dayContent,
	)

	// ========================================================
	// LAYOUT
	// ========================================================

	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		"  ",
		rightCol,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,

		top,

		"",

		dayPanel,

		"",

		theme.Paint(
			theme.Muted,
			"",
			false,
			"j/k navigate   space complete   f focus   a add   r reschedule   [ / ] day   q quit",
		),
	)

	return altScreen(content)
}

// ============================================================
// HELPERS
// ============================================================

func rowState(
	b store.Block,
	nowMin int,
) string {
	switch {

	case b.Done:
		return "done"

	case nowMin < 0:
		return "planned"

	case nowMin >= b.Start &&
		nowMin < b.End:
		return "current"

	case nowMin >= b.End:
		return "missed"
	}

	return "planned"
}

func fmtDur(minutes int) string {
	switch {

	case minutes < 60:
		return fmt.Sprintf(
			"%dm",
			minutes,
		)

	case minutes%60 == 0:
		return fmt.Sprintf(
			"%dh",
			minutes/60,
		)
	}

	return fmt.Sprintf(
		"%dh%02dm",
		minutes/60,
		minutes%60,
	)
}

func renderRow(
	b store.Block,
	selected bool,
	nowMin int,
	titleW int,
) string {
	var (
		dot    string
		dotCol string
		txtCol string
		bold   bool
	)

	switch rowState(b, nowMin) {

	case "done":
		dot = "●"
		dotCol = theme.Done
		txtCol = theme.Muted

	case "current":
		dot = "●"
		dotCol = theme.Current
		txtCol = theme.Text
		bold = true

	case "missed":
		dot = "○"
		dotCol = theme.Missed
		txtCol = theme.Muted

	default:
		dot = "○"
		dotCol = theme.Planned
		txtCol = theme.Text
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

	sp := theme.Paint(
		"",
		bg,
		false,
		" ",
	)

	timeRange := fmt.Sprintf(
		"%02d:%02d–%02d:%02d",
		b.Start/60,
		b.Start%60,
		b.End/60,
		b.End%60,
	)

	return theme.Paint(
		theme.Current,
		bg,
		false,
		marker,
	) +
		sp +

		theme.Paint(
			dotCol,
			bg,
			false,
			dot,
		) +
		sp +

		theme.Paint(
			theme.Muted,
			bg,
			false,
			timeRange,
		) +
		sp +
		sp +

		theme.Paint(
			txtCol,
			bg,
			bold,
			fmt.Sprintf(
				"%-*s",
				titleW,
				b.Title,
			),
		) +
		sp +
		sp +

		theme.Paint(
			theme.Muted,
			bg,
			false,
			fmt.Sprintf(
				"%-6s",
				fmtDur(b.End-b.Start),
			),
		) +
		sp
}

func centerBlock(
	s string,
	width int,
) string {
	lines := strings.Split(
		strings.TrimRight(s, "\n"),
		"\n",
	)

	var b strings.Builder

	for _, line := range lines {
		visible := ansiVisibleLen(line)

		pad := (width - visible) / 2

		if pad < 0 {
			pad = 0
		}

		b.WriteString(
			strings.Repeat(
				" ",
				pad,
			),
		)

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

// ============================================================
// MAIN
// ============================================================

func main() {
	if len(os.Args) > 1 {
		if runCLI(os.Args[1:]) {
			return
		}

		switch os.Args[1] {

		case "bobo":
			p := prefs.Load()

			switch p.Mascot {

			case "bear":
				p.Mascot = "cat"
				fmt.Println(
					"bobo went to sleep... the cat has arrived 🐱",
				)

			case "cat":
				p.Mascot = "cup"
				fmt.Println(
					"the cat has gone back into hiding 🧋",
				)

			default:
				p.Mascot = "bear"
				fmt.Println(
					"shh... bobo the bear has taken over your cup 🐻",
				)
			}

			_ = prefs.Save(p)
			return

		case "theme":
			if len(os.Args) < 3 {
				fmt.Println(
					"usage: boba theme <default|party>",
				)
				return
			}

			p := prefs.Load()
			p.Theme = os.Args[2]

			_ = prefs.Save(p)

			fmt.Println(
				"theme set to",
				p.Theme,
			)

			return
		}
	}

	p := prefs.Load()

	theme.SetTheme(p.Theme)

	if _, err := tea.NewProgram(
		newModel(p),
	).Run(); err != nil {
		log.Fatal(err)
	}
}
