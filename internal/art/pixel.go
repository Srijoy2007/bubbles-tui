package art

import (
	"fmt"
	"strconv"
	"strings"
)

// Palette maps a bitmap character to a hex color. '.' is always transparent.
type Palette map[byte]string

var reset = "\033[0m"

// BobaCup is a 14×20 pixel bitmap: cup, lid, straw and two rows of pearls.
var BobaCup = []string{
	strings.Repeat(".", 14),
	strings.Repeat(".", 8) + "p" + strings.Repeat(".", 5),
	strings.Repeat(".", 7) + "pp" + strings.Repeat(".", 5),
	strings.Repeat(".", 6) + "pp" + strings.Repeat(".", 6),
	strings.Repeat(".", 6) + "p" + strings.Repeat(".", 7),
	strings.Repeat(".", 2) + strings.Repeat("k", 10) + strings.Repeat(".", 2),
	".k" + strings.Repeat("c", 10) + "k.",
	".k" + strings.Repeat("m", 10) + "k.",
	".k" + strings.Repeat("m", 10) + "k.",
	".k" + strings.Repeat("m", 10) + "k.",
	".k" + strings.Repeat("m", 10) + "k.",
	strings.Repeat(".", 2) + "k" + strings.Repeat("m", 8) + "k" + strings.Repeat(".", 2),
	strings.Repeat(".", 2) + "k" + strings.Repeat("m", 8) + "k" + strings.Repeat(".", 2),
	strings.Repeat(".", 2) + "k" + "mdmdmdmd" + "k" + strings.Repeat(".", 2),
	strings.Repeat(".", 2) + "k" + "dmdmdmdm" + "k" + strings.Repeat(".", 2),
	strings.Repeat(".", 3) + "k" + strings.Repeat("m", 6) + "k" + strings.Repeat(".", 3),
	strings.Repeat(".", 3) + "k" + strings.Repeat("m", 6) + "k" + strings.Repeat(".", 3),
	strings.Repeat(".", 3) + "k" + strings.Repeat("m", 6) + "k" + strings.Repeat(".", 3),
	strings.Repeat(".", 4) + "k" + strings.Repeat("m", 4) + "k" + strings.Repeat(".", 4),
	strings.Repeat(".", 4) + strings.Repeat("k", 6) + strings.Repeat(".", 4),
}

var BobaPalette = Palette{
	'k': "#3B2A22", // espresso outline
	'c': "#F3E9D7", // lid
	'm': "#7A553A", // tea
	'd': "#241A14", // pearls
	'p': "#E0A6A6", // straw
}

// Render turns a bitmap into ANSI truecolor half-block art. Every two
// source rows become one terminal row: ▀'s foreground paints the top
// pixel, its background paints the bottom.
func Render(bitmap []string, p Palette) string {
	var b strings.Builder
	for y := 0; y < len(bitmap); y += 2 {
		top := bitmap[y]
		bottom := strings.Repeat(".", len(top))
		if y+1 < len(bitmap) {
			bottom = bitmap[y+1]
		}
		for x := 0; x < len(top); x++ {
			b.WriteString(cell(top[x], bottom[x], p))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func cell(top, bottom byte, p Palette) string {
	tc, tOK := p[top]
	bc, bOK := p[bottom]
	switch {
	case !tOK && !bOK:
		return " "
	case tOK && !bOK:
		return fg(tc) + "▀" + reset
	case !tOK && bOK:
		return fg(bc) + "▄" + reset
	default:
		return fg(tc) + bg(bc) + "▀" + reset
	}
}

func fg(hex string) string { r, g, b := hexRGB(hex); return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b) }
func bg(hex string) string { r, g, b := hexRGB(hex); return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b) }

func hexRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	r, _ := strconv.ParseInt(hex[0:2], 16, 32)
	g, _ := strconv.ParseInt(hex[2:4], 16, 32)
	b, _ := strconv.ParseInt(hex[4:6], 16, 32)
	return int(r), int(g), int(b)
}
