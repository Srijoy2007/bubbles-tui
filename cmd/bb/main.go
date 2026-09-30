package main

import (
	"fmt"
	"log"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"

	"github.com/Srijoy2007/bubbles-tui/internal/board"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

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
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.err != nil {
		return tea.NewView(fmt.Sprintf("error loading data: %v\n", m.err))
	}
	if m.store == nil {
		return tea.NewView("loading...\n")
	}

	if m.mode == modeAdd {
		out := "add block (" + m.today + ")\n\n" + m.input.View() + "\n"
		if m.addErr != "" {
			out += "\n" + m.addErr + "\n"
		}
		out += "\nenter save · esc cancel\n"
		return tea.NewView(out)
	}

	out := "boba — " + m.today + "\n\n"
	blocks := m.store.On(m.today)
	if len(blocks) == 0 {
		out += "nothing planned — press a to add a block\n"
	}
	for _, b := range blocks {
		mark := " "
		if b.Done {
			mark = "x"
		}
		out += fmt.Sprintf("[%s] %02d:%02d-%02d:%02d  %s\n",
			mark, b.Start/60, b.Start%60, b.End/60, b.End%60, b.Title)
	}
	now := time.Now()
	out += "\n" + board.Render(blocks, 5, now.Hour()*60+now.Minute()) + "\n"
	out += "\na add · q quit\n"
	return tea.NewView(out)
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
