package store

import "errors"

type Block struct {
	ID    int
	Date  string
	Start int
	End   int
	Title string
	Done  bool
}

func (b Block) Duration() int {
	return b.End - b.Start
}

func (b Block) Overlaps(o Block) bool {
	if b.Date != o.Date {
		return false
	}
	return b.Start < o.End && o.Start < b.End
}

func (b Block) Validate() error {
	if b.Title == "" {
		return errors.New("title cannot be empty")
	}
	if b.Start < 0 || b.End > 1440 {
		return errors.New("time is out of range")
	}
	if b.End <= b.Start {
		return errors.New("end must be after start")
	}
	return nil
}
