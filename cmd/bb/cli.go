package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Srijoy2007/bubbles-tui/internal/board"
	"github.com/Srijoy2007/bubbles-tui/internal/heatmap"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

const version = "0.1.0"

const helpText = `boba — time blocking for your terminal

usage:
  boba                    launch the interactive app
  boba add <range> <title> [-d date]
                           add a block, e.g. boba add 9-10:30 "Deep work"
  boba done <id>           mark a block done
  boba undone <id>         mark a block not done
  boba rm <id>              remove a block
  boba ls [date]           list a day's blocks (today if omitted)
  boba now                 what's happening right now
  boba board [date]        print the day track for a date
  boba streak [weeks]      show your focus heatmap (default 12 weeks)
  boba theme <name>        set the theme — default, party, catppuccin, gruvbox, tokyonight
  boba path                print the data file location
  boba --help, -h          show this help
  boba --version, -v       print the version

dates: today, tomorrow, yesterday, or YYYY-MM-DD`

func runCLI(args []string) bool {
	if len(args) == 0 {
		return false
	}

	switch args[0] {
	case "--help", "-h", "help":
		fmt.Println(helpText)
		return true

	case "--version", "-v":
		fmt.Println("boba " + version)
		return true

	case "path":
		fmt.Println(store.Path())
		return true

	case "add":
		cliAdd(args[1:])
		return true

	case "done":
		cliSetDone(args[1:], true)
		return true

	case "undone":
		cliSetDone(args[1:], false)
		return true

	case "rm":
		cliRemove(args[1:])
		return true

	case "ls":
		cliLs(args[1:])
		return true

	case "now":
		cliNow()
		return true

	case "board":
		cliBoard(args[1:])
		return true

	case "streak", "heatmap":
		cliStreak(args[1:])
		return true
	}

	return false
}

func mustLoad() *store.Store {
	s, err := store.Load(store.Path())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	return s
}

func cliAdd(args []string) {
	date := time.Now().Format("2006-01-02")
	var rest []string
	for i := 0; i < len(args); i++ {
		if args[i] == "-d" && i+1 < len(args) {
			date = parseDateArg(args[i+1])
			i++
			continue
		}
		rest = append(rest, args[i])
	}
	if len(rest) < 2 {
		fmt.Println("usage: boba add <range> <title> [-d date]")
		os.Exit(1)
	}
	raw := rest[0]
	for _, w := range rest[1:] {
		raw += " " + w
	}
	b, err := parseAdd(date, raw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	s := mustLoad()
	added, err := s.Add(b)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := s.Save(store.Path()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("added #%d %02d:%02d-%02d:%02d %s\n",
		added.ID, added.Start/60, added.Start%60, added.End/60, added.End%60, added.Title)
}

func cliStreak(args []string) {
	weeks := 12
	if len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			weeks = n
		}
	}
	s := mustLoad()
	fmt.Println(heatmap.Render(s.DoneMinutes(), weeks))
}

func cliSetDone(args []string, done bool) {
	if len(args) < 1 {
		fmt.Println("usage: boba done <id>  (or undone)")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: id must be a number")
		os.Exit(1)
	}
	s := mustLoad()
	blk := s.Find(id)
	if blk == nil {
		fmt.Fprintf(os.Stderr, "error: no block #%d\n", id)
		os.Exit(1)
	}
	blk.Done = done
	if err := s.Save(store.Path()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	state := "undone"
	if done {
		state = "done"
	}
	fmt.Printf("#%d marked %s: %s\n", id, state, blk.Title)
}

func cliRemove(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: boba rm <id>")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: id must be a number")
		os.Exit(1)
	}
	s := mustLoad()
	if err := s.Remove(id); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := s.Save(store.Path()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("removed #%d\n", id)
}

func cliLs(args []string) {
	date := time.Now().Format("2006-01-02")
	if len(args) > 0 {
		date = parseDateArg(args[0])
	}
	s := mustLoad()
	blocks := s.On(date)
	if len(blocks) == 0 {
		fmt.Println(date + ": nothing planned")
		return
	}
	fmt.Println(date + ":")
	for _, b := range blocks {
		mark := " "
		if b.Done {
			mark = "x"
		}
		fmt.Printf("  [%s] #%d %02d:%02d-%02d:%02d  %s\n",
			mark, b.ID, b.Start/60, b.Start%60, b.End/60, b.End%60, b.Title)
	}
}

func cliNow() {
	s := mustLoad()
	today := time.Now().Format("2006-01-02")
	now := time.Now()
	nowMin := now.Hour()*60 + now.Minute()
	blocks := s.On(today)
	for _, b := range blocks {
		if b.Status(nowMin) == store.StatusCurrent {
			fmt.Printf("▶ %s (until %02d:%02d)\n", b.Title, b.End/60, b.End%60)
			return
		}
	}
	fmt.Println("nothing scheduled right now")
}

func cliBoard(args []string) {
	date := time.Now().Format("2006-01-02")
	if len(args) > 0 {
		date = parseDateArg(args[0])
	}
	s := mustLoad()
	nowMin := -1
	if date == time.Now().Format("2006-01-02") {
		now := time.Now()
		nowMin = now.Hour()*60 + now.Minute()
	}
	fmt.Println(board.Render(s.On(date), 5, nowMin))
}

func parseDateArg(s string) string {
	now := time.Now()
	switch s {
	case "today":
		return now.Format("2006-01-02")
	case "tomorrow":
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	case "yesterday":
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", s); err == nil {
		return s
	}
	fmt.Fprintf(os.Stderr, "bad date %q, use today/tomorrow/yesterday/YYYY-MM-DD\n", s)
	os.Exit(1)
	return ""
}
