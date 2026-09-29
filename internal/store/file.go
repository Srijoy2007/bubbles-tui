package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Path returns the data file location: $TB_FILE, else ~/.timeblock/data.json.
func Path() string {
	if p := os.Getenv("TB_FILE"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".timeblock", "data.json")
}

// Load reads the store from path. A missing file yields a fresh store.
func Load(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Store{Next: 1}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", path, err)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("load %s: %w", path, err)
	}
	if s.Next < 1 {
		s.Next = 1
	}
	return &s, nil
}

// Save writes the store atomically (temp file + rename).
func (s *Store) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("save %s: %w", path, err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("save %s: %w", path, err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("save %s: %w", path, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("save %s: %w", path, err)
	}
	return nil
}
