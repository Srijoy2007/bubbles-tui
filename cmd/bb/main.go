package main

import (
	"fmt"
	"log"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/Srijoy2007/bubbles-tui/internal/board"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

type loadedMsg struct {
	s   *store.Store
	err error
}

type model struct {
	store *store.Store
	today string
	err   error
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
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
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
	nowMin := now.Hour()*60 + now.Minute()
	out += "\n" + board.Render(blocks, 5, nowMin) + "\n"
	out += "\nq to quit\n"

	return tea.NewView(out)
}

func main() {
	m := model{today: time.Now().Format("2006-01-02")}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal(err)
	}
}
