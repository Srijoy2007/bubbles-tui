package store

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "data.json")

	s := &Store{Next: 1}
	if _, err := s.Add(Block{Date: "2026-09-30", Start: 540, End: 600, Title: "A"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Add(Block{Date: "2026-09-30", Start: 600, End: 660, Title: "B"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, s) {
		t.Errorf("round trip mismatch:\n got %+v\nwant %+v", got, s)
	}
}

func TestLoadMissing(t *testing.T) {
	s, err := Load(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Next != 1 || len(s.Blocks) != 0 {
		t.Errorf("got %+v, want fresh store with Next=1", s)
	}
}
