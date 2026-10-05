package prefs

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Prefs struct {
	Mascot string `json:"mascot"` // "cup" | "bear"
	Theme  string `json:"theme"`  // "default" | "party"
}

func Path() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".boba", "prefs.json")
}

func Load() Prefs {
	p := Prefs{Mascot: "cup", Theme: "default"}
	data, err := os.ReadFile(Path())
	if err != nil {
		return p
	}
	_ = json.Unmarshal(data, &p)
	if p.Mascot == "" {
		p.Mascot = "cup"
	}
	if p.Theme == "" {
		p.Theme = "default"
	}
	return p
}

func Save(p Prefs) error {
	if err := os.MkdirAll(filepath.Dir(Path()), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0o644)
}
