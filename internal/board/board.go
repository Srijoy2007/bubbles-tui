package board

import (


	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

// color hex values from the design
const (
	colGreen    = "38;2;92;242;165"  // #5CF2A5 focused
	colDimGreen = "38;2;46;125;91"   // #2E7D5B partial
	colGray     = "38;2;122;122;122" // #7A7A7A empty / planned
	colRed      = "38;2;229;72;77"   // #E5484D missed
)

func colorize(code, glyph string) string {
	return "\033[" + code + "m" + glyph + "\033[0m"
}


func Render(blocks []store.Block, resolution int, now int) string {
	numCells := 1440 / resolution
	cells := make([]string, numCells)
	for i := range cells {
		cells[i] = "empty"
	}

	for _, b := range blocks {
		startCell := b.Start / resolution
		endCell := b.End / resolution
		for i := startCell; i < endCell; i++ {
			cellTime := i * resolution
			switch {
			case b.Done:
				cells[i] = "focused"
			case now >= 0 && cellTime >= now:
				cells[i] = "planned" // hasn't happened yet
			default:
				cells[i] = "missed" // was supposed to happen, wasn't done
			}
		}
	}

	glyphs := map[string]string{
		"focused": colorize(colGreen, "●"),
		"missed":  colorize(colRed, "○"),
		"planned": colorize(colGray, "◌"),
		"empty":   colorize(colGray, "·"),
	}

	result := ""
	for _, c := range cells {
		result += glyphs[c]
	}
	return result
}
