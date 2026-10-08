package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func Path() string {
	if p := os.Getenv("BOBA_FILE"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".boba", "data.json")
}

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

func (s *Store) On(date string) []Block{
	var out []Block
	for _,b := range s.Blocks{
		if b.Date == date {
			out = append(out,b)
		}
	}
	sort.Slice(out, func(i,j int) bool {return out[i].Start < out[j].Start})
	return out
}

func (s *Store) Reschedule(id int, newDate string) error {
	blk := s.Find(id)
	if blk == nil {
		return errors.New("block not found")
	}
	candidate := *blk
	candidate.Date = newDate
	for _, e := range s.Blocks {
		if e.ID == id {
			continue
		}
		if candidate.Overlaps(e) {
			return ErrOverlap
		}
	}
	blk.Date = newDate
	return nil
}
func (s *Store) Find(id int) *Block {
	for i := range s.Blocks {
		if s.Blocks[i].ID == id {
			return &s.Blocks[i]
		}
	}
	return nil
}
func (s *Store) DoneMinutes() map[string]int {
	out := map[string]int{}
	for _, b := range s.Blocks {
		if !b.Done {
			continue
		}
		out[b.Date] += b.End - b.Start
	}
	return out
}
func (s *Store) Remove(id int) error {
	for i, b := range s.Blocks {
		if b.ID == id {
			s.Blocks = append(s.Blocks[:i], s.Blocks[i+1:]...)
			return nil
		}
	}
	return errors.New("block not found")
}
