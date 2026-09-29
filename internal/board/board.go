package board

import "github.com/Srijoy2007/bubbles-tui/internal/store"

func Render(blocks []store.Block, resolution int) string {
	numCells := 1440 / resolution
	cells := make([]string, numCells)
	for i := range cells {
		cells[i] = "empty"
	}

	for _, b := range blocks {
		state := "missed"
		if b.Done {
			state = "focused"
		}
		startCell := b.Start / resolution
		endCell := b.End / resolution
		for i := startCell; i < endCell; i++ {
			cells[i] = state
		}
	}

	glyphs := map[string]string{
		"focused": "●",
		"missed":  "○",
		"empty":   "·",
	}

	result := ""
	for _, c := range cells {
		result += glyphs[c]
	}
	return result
}
