package audio

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

type Track struct {
	Name string
	URL  string
}

// Tracks[0] is always "off".
var Tracks = []Track{
	{Name: "off"},
	{Name: "lofi", URL: "http://ice1.somafm.com/groovesalad-128-mp3"},
	{Name: "ambient", URL: "http://ice1.somafm.com/dronezone-128-mp3"},
}

var ErrNoPlayer = errors.New("no audio player found, install mpv or ffplay")

// Player manages one background playback process at a time.
type Player struct {
	cmd *exec.Cmd
	log *os.File
}

func (p *Player) Start(url string) error {
	p.Stop()
	bin, args, ok := findPlayer(url)
	if !ok {
		return ErrNoPlayer
	}
	cmd := exec.Command(bin, args...)
	if f := openLog(); f != nil {
		cmd.Stdout = f
		cmd.Stderr = f
		p.log = f
	}
	if err := cmd.Start(); err != nil {
		p.closeLog()
		return err
	}
	p.cmd = cmd
	return nil
}

func (p *Player) Stop() {
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		_ = p.cmd.Wait()
	}
	p.cmd = nil
	p.closeLog()
}

func (p *Player) closeLog() {
	if p.log != nil {
		_ = p.log.Close()
		p.log = nil
	}
}

func (p *Player) Playing() bool { return p.cmd != nil }

func openLog() *os.File {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	dir := filepath.Join(home, ".boba")
	if os.MkdirAll(dir, 0o755) != nil {
		return nil
	}
	f, err := os.OpenFile(filepath.Join(dir, "audio.log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil
	}
	return f
}

func findPlayer(url string) (string, []string, bool) {
	candidates := []struct {
		bin  string
		args func(string) []string
	}{
		{"mpv", func(u string) []string { return []string{u, "--no-video", "--really-quiet"} }},
		{"ffplay", func(u string) []string { return []string{"-nodisp", "-autoexit", "-loglevel", "quiet", u} }},
	}
	for _, c := range candidates {
		if path, err := exec.LookPath(c.bin); err == nil {
			return path, c.args(url), true
		}
	}
	return "", nil, false
}
