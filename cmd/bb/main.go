package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"

	"github.com/Srijoy2007/bubbles-tui/internal/board"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

var noColor = os.Getenv("NO_COLOR") != ""

const (
	reverseOn  = "\033[7m"
	reverseOff = "\033[0m"
)

func highlight(s string) string {
	if noColor {
		return s
	}
	return reverseOn + s + reverseOff
}

type loadedMsg struct {
	s   *store.Store
	err error
}

type mode int

const (
	modeList mode = iota
	modeAdd
)

type model struct {
	store  *store.Store
	today  string
	err    error
	mode   mode
	input  textinput.Model
	addErr string
	cursor int
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "9-10:30 Deep work"
	ti.CharLimit = 100
	return model{today: time.Now().Format("2006-01-02"), input: ti}
}

func loadCmd() tea.Msg {
	s, err := store.Load(store.Path())
	return loadedMsg{s: s, err: err}
}

func (m model) Init() tea.Cmd {
	return loadCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case loadedMsg:
		m.store = msg.s
		m.err = msg.err
		return m, nil

	case tea.KeyPressMsg:
		if m.mode == modeAdd {
			switch msg.String() {
			case "esc":
				m.mode = modeList
				m.input.Reset()
				m.addErr = ""
				return m, nil
			case "enter":
				b, err := parseAdd(m.today, m.input.Value())
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
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			m.mode = modeAdd
			m.input.Focus()
			return m, textinput.Blink
		case "j", "down":
			blocks := m.store.On(m.today)
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
			blocks := m.store.On(m.today)
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

func (m model) View() tea.View {
	if m.err != nil {
		return altScreen(fmt.Sprintf("error loading data: %v\n", m.err))
	}
	if m.store == nil {
		return altScreen("loading...\n")
	}

	if m.mode == modeAdd {
		out := "add block (" + m.today + ")\n\n" + m.input.View() + "\n"
		if m.addErr != "" {
			out += "\n" + m.addErr + "\n"
		}
		out += "\nenter save · esc cancel\n"
		return altScreen(out)
	}

	out := "boba — " + m.today + "\n\n"
	blocks := m.store.On(m.today)
	if len(blocks) == 0 {
		out += "nothing planned — press a to add a block\n"
	}
	for i, b := range blocks {
		mark := " "
		if b.Done {
			mark = "x"
		}
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		line := fmt.Sprintf("%s[%s] %02d:%02d-%02d:%02d  %s",
			cursor, mark, b.Start/60, b.Start%60, b.End/60, b.End%60, b.Title)
		if i == m.cursor {
			line = highlight(line)
		}
		out += line + "\n"
	}
	now := time.Now()
	out += "\n" + board.Render(blocks, 5, now.Hour()*60+now.Minute()) + "\n"
	out += "\nspace done · a add · j/k move · q quit\n"
	return altScreen(out)
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
