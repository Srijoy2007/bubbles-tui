<p align="center">
  <img src="docs/boba.png" width="150" alt="boba">
</p>

<h1 align="center">boba</h1>

<p align="center">
  <strong>Time blocking for your terminal.</strong><br>
  Watch your day fill up.
</p>

<p align="center">
  <a href="https://github.com/Srijoy2007/bubbles-tui/releases">
    <img alt="release" src="https://img.shields.io/github/v/release/Srijoy2007/bubbles-tui?style=flat-square&labelColor=2F2C4A&color=E797BE">
  </a>
  <a href="go.mod">
    <img alt="go version" src="https://img.shields.io/github/go-mod/go-version/Srijoy2007/bubbles-tui?style=flat-square&labelColor=2F2C4A&color=B1A4DC">
  </a>
  <a href="LICENSE">
    <img alt="license" src="https://img.shields.io/github/license/Srijoy2007/bubbles-tui?style=flat-square&labelColor=2F2C4A&color=9694AF">
  </a>
  <img alt="platform" src="https://img.shields.io/badge/platform-linux%20%7C%20macOS-D27D9B?style=flat-square&labelColor=2F2C4A">
  <img alt="built with Bubble Tea" src="https://img.shields.io/badge/built%20with-Bubble%20Tea-E797BE?style=flat-square&labelColor=2F2C4A">
  <img alt="storage" src="https://img.shields.io/badge/storage-one%20JSON%20file-B1A4DC?style=flat-square&labelColor=2F2C4A">
</p>

<p align="center">
  <img src="docs/demo.gif" alt="boba demo" width="720">
</p>

---

## The idea

Most planners ask **what** you want to get done.

**boba asks when.**

You cut your day into blocks, and boba checks every block against the real clock:

```text
        upcoming
           ↓
        current
           ↓
         done
           ↓
        missed
```

Finished blocks fill in like pearls at the bottom of a cup. Current blocks move with the clock. Missed blocks stay visible until you deal with them.

Everything runs in the terminal and lives in one JSON file.

No account.

No cloud.

No subscription.

No notification nagging you.

Just your day.

---

## Why boba?

Traditional todo lists are good at answering:

> "What do I need to do?"

But they don't answer:

> "What should I be doing right now?"

boba turns your task list into a **time-aware schedule**.

Instead of:

```text
[ ] Study
[ ] Gym
[ ] Deep work
[ ] Read
```

you get:

```text
09:00 ─────── 10:30   Deep work     ●●●●●●●●○○
11:00 ─────── 12:00   Study         ●●●○○○○○○○
14:00 ─────── 15:30   Gym           ○○○○○○○○○○
```

The clock becomes part of the interface.

---

# Features

### 🗓 Time blocking

Plan your day using simple time ranges:

```sh
boba add 9-10:30 "Deep work"
boba add 11-12 "Study"
boba add 14-15:30 "Gym"
```

You can also enter blocks directly from the TUI.

---

### 📊 Day Track

The entire day is represented on one compact timeline.

```text
╭─ DAY ────────────────────────────────────────────────────────────────────────╮
│    00    02    04    06    08    10    12    14    16    18    20    22      │
│    ···························━━━━━·━━━······┄┄┄┄┄·························  │
│                                      ▲ now                                   │
│                                                                              │
│    done 1h30m   missed 0m   of 4h planned                                    │
│                                                                              │
│    ● 09:00–10:30  Deep work  ●●●●●●●●●●                                      │
│    ● 11:00–12:00  Study  ●●●○○○○○○○                                          │
│    ○ 14:00–15:30  Gym  ○○○○○○○○○○                                            │
│                                                                              │
│    ● done    ● current    ● planned    ● missed                              │
╰──────────────────────────────────────────────────────────────────────────────╯
```

The day track gives you a quick answer to:

- What have I planned?
- What am I doing now?
- How much time have I completed?
- What did I miss?
- What's coming next?

---

### ⏱ Focus mode

Press `f` on a block and boba turns it into a focused countdown.

```text
╔═ FOCUS ════════════════════════════════════════════════╗
║   ░░░░░░██░░  ░░██████░░  ░░  ░░░░██░░░░  ██████████   ║
║   ░░░░████░░  ██░░░░░░██  ░░  ░░████░░░░  ░░░░░░░░██   ║
║   ░░██░░██░░  ░░░░░░░░██  ██  ░░░░██░░░░  ░░░░░░██░░   ║
║   ██░░░░██░░  ░░░░░░██░░  ░░  ░░░░██░░░░  ░░░░██░░░░   ║
║   ██████████  ░░░░██░░░░  ██  ░░░░██░░░░  ░░██░░░░░░   ║
║   ░░░░░░██░░  ░░██░░░░░░  ░░  ░░░░██░░░░  ░░██░░░░░░   ║
║   ░░░░░░██░░  ██████████  ░░  ░░██████░░  ░░██░░░░░░   ║
╚════════════════════════════════════════════════════════╝

                  Deep work  09:00–10:30

      ████████████████░░░░░░░░░░░░░░░░░░░░░░░░   40%

                        ▶ running

                 ♪ lofi  ▂▄▆█▇▅▃▂▃▅▇█▆▄▂▁

space pause · + 5min · d done · m music: lofi · esc leave
```

Focus mode understands the block's actual schedule.

If the block is already underway, the countdown runs until its scheduled end.

If it hasn't started yet, the timer represents the full block duration.

When the timer reaches zero, the block is marked complete.

You can also:

- pause/resume
- add five minutes
- finish early
- leave focus mode
- toggle music

---

### 🎧 Ambient music

boba can optionally play background music while you focus.

Available modes:

```text
music: off
music: lofi
music: ambient
```

The interface includes a small equalizer that reacts while music is playing and goes flat when paused.

Music playback uses an external player such as `mpv` or `ffplay`.

The included streams are from [SomaFM](https://somafm.com/):

- Groove Salad
- Drone Zone

SomaFM is listener-supported. If you enjoy the streams, consider [supporting them](https://somafm.com/support/).

---

### 🔥 Focus history

At the end of the week, see how consistently you showed up.

```text
      Mon Tue Wed Thu Fri Sat Sun
W1     ██  ██  ░░  ██  ██  ░░  ██
W2     ██  ░░  ██  ██  ░░  ██  ██
W3     ░░  ██  ██  ██  ██  ░░  ██
W4     ██  ██  ██  ░░  ██  ██  ██
```

Run:

```sh
boba streak 12
```

to display up to twelve weeks of focused hours as a terminal heatmap.

---

### 🖥 Status bar

boba doesn't need to run in the background.

`boba now` reads your schedule and prints one line:

```text
▶ Deep work · until 10:30
```

That makes it work nicely with terminal multiplexers and desktop status bars.

#### tmux

```tmux
set -g status-interval 30
set -g status-right '#(boba now)'
```

#### Waybar

```json
"custom/boba": {
  "exec": "boba now",
  "interval": 30
}
```

The same schedule powers both the TUI and the status bar.

---

## Install

boba is a single Go binary.

### Go install

If you already have Go installed:

```sh
go install github.com/Srijoy2007/bubbles-tui/cmd/boba@latest
```

Make sure your Go binary directory is in your `PATH`.

You can verify the installation with:

```sh
boba --version
```

---

### Build from source

```sh
git clone https://github.com/Srijoy2007/bubbles-tui.git
cd bubbles-tui

go build -o boba ./cmd/boba

mkdir -p ~/.local/bin
cp boba ~/.local/bin/
```

Then:

```sh
boba
```

---

### Requirements

boba currently supports:

- Linux
- macOS

A truecolor-capable terminal is recommended.

Good options include:

- Kitty
- Ghostty
- WezTerm
- Alacritty
- foot
- iTerm2

Music playback additionally requires either:

```sh
mpv
```

or:

```sh
ffplay
```

on your `PATH`.

Music is optional. The rest of boba works without it.

---

## Use it

Start boba:

```sh
boba
```

### Main view

```text
j / k          move between blocks
↑ / ↓          move between blocks

[ / ]          previous / next day

a              add a block
space          toggle done
f              focus selected block
r              reschedule missed block

q              quit
```

### Focus mode

```text
space          pause / resume
+              add five minutes
d              mark done
m              cycle music
esc            leave focus mode
```

### Adding a block

Type a time range followed by a title:

```text
9-10:30 Deep work
14:00-15:30 Review
18-19 Gym
```

boba understands both compact and explicit times.

---

# CLI

Every core action is also available from the shell.

This means you can use boba without opening the interactive interface.

### Add blocks

```console
$ boba add 9-10:30 "Deep work"
added #1 09:00-10:30 Deep work

$ boba add 11-12 "Study"
added #2 11:00-12:00 Study

$ boba add 14-15:30 "Gym"
added #3 14:00-15:30 Gym
```

### List the day

```console
$ boba ls

2026-10-08:
  [ ] #1 09:00-10:30  Deep work
  [ ] #2 11:00-12:00  Study
  [ ] #3 14:00-15:30  Gym
```

### See what's happening now

```console
$ boba now

▶ Deep work (until 10:30)
```

### Mark a block done

```console
$ boba done 1

#1 marked done: Deep work
```

### Undo completion

```console
$ boba undone 1
```

### Remove a block

```console
$ boba rm 3
```

### Print the day track

```console
$ boba board
```

### Show your focus history

```console
$ boba streak
```

---

# Commands

```text
boba                           open the app

boba add <range> <title> [-d date]
                               add a block

boba done <id>                 mark a block done
boba undone <id>               mark a block not done
boba rm <id>                   remove a block

boba ls [date]                 list a day's blocks
boba now                       show what's happening now
boba board [date]              print the day track
boba streak [weeks]            print the focus heatmap

boba theme <name>              set the theme
boba path                      print the data file location

boba --help, -h                show help
boba --version, -v             print the version
```

Dates can be:

```text
today
tomorrow
yesterday
YYYY-MM-DD
```

For example:

```sh
boba ls tomorrow
boba add 10-11 "Algorithms" -d tomorrow
boba board 2026-10-15
```

---

# Themes

boba comes with several built-in themes.

```sh
boba theme tokyonight
```

| Theme | Style |
|---|---|
| `default` | boba pink on lavender |
| `catppuccin` | Catppuccin-inspired dark palette |
| `gruvbox` | warm and retro |
| `tokyonight` | deep blue night |
| `party` | neon and loud |

Use:

```sh
boba theme default
```

to return to the default theme.

For environments where ANSI color is undesirable:

```sh
NO_COLOR=1 boba
```

---

# Your data

boba is deliberately local-first.

Your schedule lives in:

```text
~/.boba/data.json
```

Preferences such as your theme live in:

```text
~/.boba/prefs.json
```

You can inspect the location at any time:

```sh
boba path
```

The data format is intentionally simple:

```json
{
  "next": 4,
  "blocks": [
    {
      "id": 1,
      "date": "2026-10-08",
      "start": 540,
      "end": 630,
      "title": "Deep work",
      "done": true
    }
  ]
}
```

Times are stored as minutes since midnight.

That means your data is:

- human-readable
- easy to back up
- easy to script
- easy to inspect
- independent of a database

You can also point boba at another data file:

```sh
BOBA_FILE=~/work/boba.json boba
```

This can be useful if you want separate schedules for different projects or environments.

---

# Local-first by design

boba intentionally doesn't have:

- user accounts
- cloud sync
- mandatory internet access
- a web dashboard
- subscriptions
- background daemons
- notification spam

Your schedule belongs to you.

The terminal is the interface.

The JSON file is the database.

The binary is the application.

---

# Shell scripting

Because boba exposes its core operations through the CLI, it can be composed with other Unix tools.

For example:

```sh
boba ls
```

can be piped into:

```sh
grep "Deep work"
```

or used from scripts.

You can also create quick shell aliases:

```sh
alias today='boba ls today'
alias tomorrow='boba ls tomorrow'
alias now='boba now'
```

---

# Development

Clone the repository:

```sh
git clone https://github.com/Srijoy2007/bubbles-tui.git
cd bubbles-tui
```

Run directly:

```sh
go run ./cmd/boba
```

Run the test suite:

```sh
go test ./...
```

Build:

```sh
go build ./cmd/boba
```

The project is written in Go and uses Charm's terminal UI ecosystem:

- Bubble Tea
- Bubbles
- Lip Gloss

---

# Project structure

The codebase is organized around the core parts of the application:

```text
.
├── cmd/
│   └── boba/             # CLI entry point
│
├── internal/
│   ├── art/              # terminal art and animations
│   ├── audio/            # music playback
│   ├── board/            # day track and progress rendering
│   ├── heatmap/          # focus history
│   ├── prefs/            # user preferences and themes
│   ├── store/            # schedule persistence
│   └── theme/            # visual themes
│
├── docs/
│   ├── boba.png
│   └── demo.gif
│
├── go.mod
├── go.sum
└── LICENSE
```

The exact internal structure may change as the project evolves.

---

# Philosophy

boba is built around a few simple ideas.

### Time is a resource

A task list can grow forever.

A day can't.

Giving something a time block forces you to decide where it actually belongs.

### Progress should be visible

A completed task disappearing from a list isn't very satisfying.

Seeing a block gradually fill as time passes makes progress tangible.

### Your tools should stay out of the way

boba is intentionally small.

It shouldn't become another thing you have to manage.

Open it.

Plan.

Focus.

Close it.

### Local-first is enough

You shouldn't need an account to decide what you're doing at 9 AM.

---

# Roadmap

boba is still evolving.

Potential future work includes:

- [ ] Improved weekly planning
- [ ] Better rescheduling workflows
- [ ] More detailed focus statistics
- [ ] Custom block colors
- [ ] Recurring blocks
- [ ] Daily/weekly goals
- [ ] More status-bar integrations
- [ ] Configurable keybindings
- [ ] More terminal themes
- [ ] Better accessibility for small terminals
- [ ] Optional schedule sharing
- [ ] Calendar import/export
- [ ] Desktop/widget companion
- [ ] Friend/accountability features

The goal is not to turn boba into another giant productivity suite.

The goal is to make **time blocking feel native to the terminal**.

---

# Contributing

Contributions are welcome.

If you find a bug, have an idea, or want to improve the interface:

1. Fork the repository.
2. Create a branch.

```sh
git checkout -b feature/my-feature
```

3. Make your changes.
4. Run the tests.

```sh
go test ./...
```

5. Commit your changes.

```sh
git commit -m "add my feature"
```

6. Push the branch.

```sh
git push origin feature/my-feature
```

7. Open a pull request.

For larger changes, opening an issue first is recommended so the approach can be discussed before implementation.

---

# Credits

Built with the excellent terminal tooling from [Charm](https://charm.land/):

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles)

Music streams are provided by [SomaFM](https://somafm.com/).

---

# License

boba is released under the MIT License.

See [`LICENSE`](LICENSE) for the full license text.

---

<p align="center">
  Made for people who spend too much time in terminals.
</p>

<p align="center">
  <strong>small steps · big dreams</strong>
</p>
