package audio

import (
	"errors"
	"os"
	"os/exec"
)

type Track struct {
	Name string
	URL  string
}


var Tracks = []Track{
	{Name: "off"},
	{Name: "lofi", URL: "http://ice1.somafm.com/groovesalad-128-mp3"},
	{Name: "ambient", URL: "http://ice1.somafm.com/dronezone-128-mp3"},
}

var ErrNoPlayer = errors.New("no audio player found — install mpv or ffplay.")


type Player struct {
	cmd *exec.Cmd
}

func (p *Player) Start(url string) error {
	p.Stop()
	bin, args, ok := findPlayer(url)
	if !ok {
		return ErrNoPlayer
	}
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
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
}

func (p *Player) Playing() bool {
	return p.cmd != nil
}

func findPlayer(url string) (string, []string, bool) {
	candidates := []struct {
		bin  string
		args func(string) []string
	}{
		{"mpv", func(u string) []string { return []string{u, "--no-video", "--really-quiet"} }},
		{"ffplay", func(u string) []string { return []string{"-nodisp", "-autoexit", "-loglevel", "quiet", u} }},
		{"afplay", func(u string) []string { return []string{u} }}, // macOS fallback
	}
	for _, c := range candidates {
		if path, err := exec.LookPath(c.bin); err == nil {
			return path, c.args(url), true
		}
	}
	return "", nil, false
}
