package store

import "errors"

type Block struct {
	ID    int    `json:"id"`
	Date  string `json:"date"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
type Store struct{
	Next int `json:"next"`
	Blocks []Block `json:"blocks"`


}

var ErrInvalid = errors.New("start must be before end, within 0-1440")
var ErrOverlap = errors.New("overlaps an existing block")

func (s *Store) Add(b Block) (Block, error){
	if err := b.Validate(); err != nil{
		return Block{},err 
	}
	for _,e := range s.Blocks{
		if b.Overlaps(e){
			return Block{},ErrOverlap

		}
	}
	b.ID = s.Next
	s.Next++
	s.Blocks = append(s.Blocks,b)

	return b,nil

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
