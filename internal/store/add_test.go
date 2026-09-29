package store

import (
	"errors"
	"testing"
)

func TestAdd(t *testing.T) {
	existing := Block{Date: "2026-09-29", Start: 540, End: 600, Title: "A"}

	cases := []struct {
		name    string
		b       Block
		wantErr bool
		isClash bool // expect ErrOverlap specifically
	}{
		{"valid", Block{Date: "2026-09-29", Start: 700, End: 760, Title: "B"}, false, false},
		{"empty title", Block{Date: "2026-09-29", Start: 700, End: 760}, true, false},
		{"start after end", Block{Date: "2026-09-29", Start: 760, End: 700, Title: "B"}, true, false},
		{"end out of range", Block{Date: "2026-09-29", Start: 1400, End: 1500, Title: "B"}, true, false},
		{"overlap", Block{Date: "2026-09-29", Start: 570, End: 630, Title: "B"}, true, true},
		{"back to back ok", Block{Date: "2026-09-29", Start: 600, End: 660, Title: "B"}, false, false},
		{"same time other date ok", Block{Date: "2026-09-30", Start: 540, End: 600, Title: "B"}, false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := &Store{Next: 1}
			if _, err := s.Add(existing); err != nil {
				t.Fatalf("setup: %v", err)
			}
			nextBefore, lenBefore := s.Next, len(s.Blocks)

			got, err := s.Add(c.b)

			if c.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if c.isClash && !errors.Is(err, ErrOverlap) {
					t.Fatalf("err = %v, want ErrOverlap", err)
				}
				if s.Next != nextBefore || len(s.Blocks) != lenBefore {
					t.Error("failed Add changed the store")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != nextBefore {
				t.Errorf("ID = %d, want %d", got.ID, nextBefore)
			}
			if s.Next != nextBefore+1 || len(s.Blocks) != lenBefore+1 {
				t.Error("store not updated correctly")
			}
		})
	}
}
