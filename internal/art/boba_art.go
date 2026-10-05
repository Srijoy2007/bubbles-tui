package art

import (
	"fmt"
	"strconv"
	"strings"
)

// Palette maps one bitmap byte to an ANSI truecolor hex value.
// '.' is transparent. 's' is the animated white shimmer.
type Palette map[byte]string

var reset = "\033[0m"

// BobaCup is a 35x35 pixel-art reconstruction of the supplied cup image.
// It is intentionally transparent around the character so it sits cleanly
// on your TUI's dark background.
var BobaCup = []string{
	".............bbbbbbb..........ff...",
	"...........iihccccccibb......ffffa.",
	"...........bahhhcccccccbjj..jffffa.",
	"...........bdhaaaaacccceaaaaffffaa.",
	"...........bicjbiijcccccccjjffffa..",
	"...........iajachhabiicccccbfff....",
	"..........badhddahhhhhbbhffcjf.....",
	"..........bhhdhadjdaahadbhffcj.....",
	".........bchhhjdcaejjaajfbhfccj....",
	".........bhhhhhajdeaajffffbhhccj...",
	"........bcchhhhhaaajjjffffhbhhcb...",
	"........bhdhhhhhjdddjffffaahjhhhde.",
	".......bcgghbddddjjfffffhggahjhhdd.",
	".......bgggggghjbffffhhgggggabahchj",
	"......bcggggggggggggccgggggeghjchhj",
	"......bgggggggggcghchggcgeggegjchab",
	".....iceggeggggggeeggeccccgeegbahab",
	".....igeegeggegeeeegeghcccceehbhaab",
	"....ejgeeedgeegeeeegegghccceehbhaeb",
	"...ebcedebbeeeeeeeeeeeeghccedciadbb",
	"...eidbbdbbeeedeeeeeeegeehedcijb...",
	"..bdcdbbddddebbeeddeeeeeeeecb......",
	".ideeeddddbbdbbdeddeedededhj.......",
	".iedbbbeedbbddddbbjdebbddgj........",
	"bcddddddbjedddddbbjddbbdcj.........",
	"bcddddddddjedbbddddddddhj..........",
	"bgddddddddddebbdddjbbdhj...........",
	".iedddddddddbdddddjbbej............",
	".igedddddddddbebbdddeb.............",
	"..bdedddddddddbbbdecj..............",
	"...eigedddddddbeddji...............",
	".....iieeddddddgcjd................",
	".......jbgeddegbj..................",
	".......jjedjddejj..................",
	".........iiiiii....................",
}

// BobaPalette is sampled/quantized from the supplied artwork.
var BobaPalette = Palette{
	'a': "#D2C2CF", // cool lavender
	'b': "#896076", // dark mauve
	'c': "#FDF9FB", // near-white highlight
	'd': "#BB8A93", // dusty pink
	'e': "#D9A6A1", // warm pink
	'f': "#8DBAB0", // mint green
	'g': "#F1CBBE", // peach
	'h': "#EEE7EB", // pale lavender-white
	'i': "#6A4958", // deep outline
	'j': "#987B8D", // mid lavender
	's': "#FFFFFF", // animated shimmer
}

// Render turns a bitmap into ANSI truecolor half-block art.
// Two bitmap rows become one terminal row, giving roughly 2x the
// vertical resolution of ordinary one-character-per-pixel ASCII art.
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
		b.WriteByte('\n')
	}

	return b.String()
}

// Width returns the visible terminal width of the bitmap.
func Width() int {
	if len(BobaCup) == 0 {
		return 0
	}
	return len(BobaCup[0])
}

// RenderFrame adds a subtle "alive" animation:
//   - the cup gently floats/bobs,
//   - a bright shimmer travels across the cup,
//   - a second glint periodically hits the mint leaf.
//
// The animation never changes the underlying pixel-art shape.
// RenderSmallFrame renders a compact version for the persistent dashboard header.
// It uses nearest-neighbour downsampling so the pixel-art character remains crisp.
func RenderSmallFrame(frame int) string {
	const smallW = 24
	const smallH = 24

	bitmap := resizeNearest(BobaCup, smallW, smallH)

	// Reuse the same animation idea, but with a gentler motion for the dashboard.
	dx := []int{0, 0, 0, 1, 1, 0, 0, -1, -1, 0}[frame%10]
	dy := []int{0, 0, -1, -1, 0, 0, 1, 1, 0, 0}[frame%10]
	bitmap = shifted(bitmap, dx, dy)

	// Small travelling highlight.
	points := [][2]int{{9, 9}, {11, 8}, {13, 8}, {15, 9}, {13, 10}, {11, 10}}
	pt := points[frame%len(points)]
	putIfPainted(bitmap, pt[0]+dx, pt[1]+dy, 's')

	return Render(bitmap, BobaPalette)
}
func resizeNearest(src []string, newW, newH int) []string {
	if len(src) == 0 || newW <= 0 || newH <= 0 {
		return nil
	}
	oldH := len(src)
	oldW := len(src[0])
	out := make([]string, newH)
	for y := 0; y < newH; y++ {
		row := make([]byte, newW)
		sy := y * oldH / newH
		if sy >= oldH {
			sy = oldH - 1
		}
		srcRow := src[sy]
		for x := 0; x < newW; x++ {
			sx := x * oldW / newW
			if sx >= len(srcRow) {
				row[x] = '.'
				continue
			}
			row[x] = srcRow[sx]
		}
		out[y] = string(row)
	}
	return out
}
func RenderFrame(frame int) string {
	const n = 8

	dx := []int{0, 0, 1, 1, 0, -1, -1, 0}[frame%n]
	dy := []int{0, -1, -1, 0, 0, 1, 1, 0}[frame%n]

	bitmap := shifted(BobaCup, dx, dy)

	// Moving highlight across the upper body.
	shimmer := [][2]int{
		{15, 14}, {18, 13}, {21, 12}, {24, 11},
		{27, 10}, {24, 13}, {21, 14}, {18, 15},
	}

	sx, sy := shimmer[frame%n][0]+dx, shimmer[frame%n][1]+dy
	putIfPainted(bitmap, sx, sy, 's')

	// Occasional tiny glint on the mint leaf.
	if frame%8 == 2 || frame%8 == 3 {
		putIfPainted(bitmap, 29+dx, 4+dy, 's')
	}

	return Render(bitmap, BobaPalette)
}

func shifted(src []string, dx, dy int) []string {
	h := len(src)
	if h == 0 {
		return nil
	}
	w := len(src[0])

	out := make([][]byte, h)
	for y := range out {
		out[y] = []byte(strings.Repeat(".", w))
	}

	for y, row := range src {
		for x := 0; x < len(row); x++ {
			if row[x] == '.' {
				continue
			}

			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < w && ny >= 0 && ny < h {
				out[ny][nx] = row[x]
			}
		}
	}

	result := make([]string, h)
	for y := range out {
		result[y] = string(out[y])
	}
	return result
}

func putIfPainted(bitmap []string, x, y int, ch byte) {
	if y < 0 || y >= len(bitmap) || x < 0 || x >= len(bitmap[y]) {
		return
	}
	if bitmap[y][x] == '.' {
		return
	}

	row := []byte(bitmap[y])
	row[x] = ch
	bitmap[y] = string(row)
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

func fg(hex string) string {
	r, g, b := hexRGB(hex)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

func bg(hex string) string {
	r, g, b := hexRGB(hex)
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

func hexRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 32)
	g, _ := strconv.ParseInt(hex[2:4], 16, 32)
	b, _ := strconv.ParseInt(hex[4:6], 16, 32)

	return int(r), int(g), int(b)
}

// ============================================================
// MASCOTS: BOBO THE BEAR and the BOBA CAT
// ============================================================
//
// Both sprites are drawn as the LEFT HALF only (12 columns) and mirrored
// at init, so they are perfectly symmetric and every row has the same
// width by construction (hand-typed full rows drift by a pixel or two).
//
//	'.' transparent   k outline   f fur       h belly / chest highlight
//	m muzzle          n nose      e eyes      p inner ear (cat)
//	s animated shimmer
//
// Each mascot has two sizes with the same footprint as the cup:
//
//	*SmallFrame  24x24 px  (24 cols x 12 rows, dashboard header)
//	*Frame       36x36 px  (36 cols x 18 rows, splash screen)

const mascotNative = 24 // size the glint coordinates below are written in

// mirror pads every half-row to the widest one and appends its reverse.
func mirror(half []string) []string {
	w := 0
	for _, r := range half {
		if len(r) > w {
			w = len(r)
		}
	}
	out := make([]string, len(half))
	for i, r := range half {
		r += strings.Repeat(".", w-len(r))
		rev := make([]byte, w)
		for x := 0; x < w; x++ {
			rev[x] = r[w-1-x]
		}
		out[i] = r + string(rev)
	}
	return out
}

var bearHalf = []string{
	"............",
	"....kkkk....",
	"...kffffk...",
	"..kfffffkkkk",
	".kffffffffff",
	".kffffffffff",
	".kffffffffff",
	".kffffffffff",
	".kfffffeffff",
	".kfffffeffff",
	".kfffffffmmm",
	".kffffffmmmm",
	".kffffffmmmn",
	".kffffffmmmm",
	"..kfffffffmm",
	"..kfffffffff",
	"...kkkkkkkkk",
	"....kfffffff",
	"...kffffhhhh",
	"...kffffhhhh",
	"...kffffhhhh",
	"...kfffffhhh",
	"....kkkkkkkk",
	"............",
}

var catHalf = []string{
	"............",
	"...kk.......",
	"..kppk......",
	"..kpppk.....",
	".kfpppfk....",
	".kffffffkkkk",
	".kffffffffff",
	".kffffffffff",
	".kffffffffff",
	".kfffeefffff",
	".kffffffmmmm",
	".kffffffmmmn",
	".kffffffmmmm",
	".kffffffffff",
	"..kfffffffff",
	"...kkkkkkkkk",
	"....kfffffff",
	"...kffffhhhh",
	"...kffffhhhh",
	"...kfffffhhh",
	"...kffffffff",
	"....kfffffff",
	".....kkkkkkk",
	"............",
}

// BobaBear / BobaCat are the 24x24 sprites; the *Big versions are the same
// art scaled 1.5x. The half is scaled BEFORE mirroring so symmetry survives.
var (
	BobaBear    = mirror(bearHalf)
	BobaBearBig = mirror(resizeNearest(bearHalf, 18, 36))
	BobaCat     = mirror(catHalf)
	BobaCatBig  = mirror(resizeNearest(catHalf, 18, 36))
)

var BearPalette = Palette{
	'k': "#2B1B17", // dark outline
	'f': "#A9683F", // teddy fur
	'h': "#D9A56C", // belly
	'm': "#E8C18A", // muzzle
	'n': "#241714", // nose
	'e': "#241714", // eyes
	's': "#FFF1D2", // shimmer
}

var CatPalette = Palette{
	'k': "#35252A", // dark outline
	'f': "#F0D7C9", // pale cream fur
	'p': "#D98D9C", // pink inner ears
	'h': "#E7BBA9", // chest highlight
	'm': "#F4E4DB", // muzzle
	'n': "#6D3541", // nose
	'e': "#35252A", // eyes
	's': "#FFF7F2", // shimmer
}

// Glint paths, in 24x24 coordinates (all land on fur).
var (
	bearGlints = [][2]int{{7, 6}, {10, 5}, {13, 6}, {16, 7}, {13, 8}, {10, 7}}
	catGlints  = [][2]int{{7, 6}, {9, 6}, {11, 6}, {13, 7}, {11, 7}, {9, 7}}
)

// mascotFrame draws a sprite with a gentle 1px bob and a travelling glint.
// The sprites keep their first and last rows empty, so the bob never clips.
func mascotFrame(sprite []string, p Palette, glints [][2]int, frame int) string {
	dy := []int{0, 0, 1, 1, 0, 0, -1, -1}[frame%8]

	bitmap := shifted(sprite, 0, dy)

	size := len(sprite)
	g := glints[frame%len(glints)]
	putIfPainted(bitmap, g[0]*size/mascotNative, g[1]*size/mascotNative+dy, 's')

	return Render(bitmap, p)
}

// BearFrame renders Bobo at splash size.
func BearFrame(frame int) string {
	return mascotFrame(BobaBearBig, BearPalette, bearGlints, frame)
}

// BearSmallFrame renders Bobo at dashboard size.
func BearSmallFrame(frame int) string {
	return mascotFrame(BobaBear, BearPalette, bearGlints, frame)
}

// CatFrame renders the cat at splash size.
func CatFrame(frame int) string {
	return mascotFrame(BobaCatBig, CatPalette, catGlints, frame)
}

// CatSmallFrame renders the cat at dashboard size.
func CatSmallFrame(frame int) string {
	return mascotFrame(BobaCat, CatPalette, catGlints, frame)
}
