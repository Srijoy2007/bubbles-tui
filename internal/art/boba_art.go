package art

import (
	"fmt"
	"strconv"
	"strings"
)

type Palette map[byte]string

var reset = "\033[0m"

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

var BobaPalette = Palette{
	'a': "#D2C2CF",
	'b': "#896076",
	'c': "#FDF9FB",
	'd': "#BB8A93",
	'e': "#D9A6A1",
	'f': "#8DBAB0",
	'g': "#F1CBBE",
	'h': "#EEE7EB",
	'i': "#6A4958",
	'j': "#987B8D",
	's': "#FFFFFF",
}
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


func Width() int {
	if len(BobaCup) == 0 {
		return 0
	}
	return len(BobaCup[0])
}
func RenderSmallFrame(frame int) string {
	const smallW = 24
	const smallH = 24

	bitmap := resizeNearest(BobaCup, smallW, smallH)

	dx := []int{0, 0, 0, 1, 1, 0, 0, -1, -1, 0}[frame%10]
	dy := []int{0, 0, -1, -1, 0, 0, 1, 1, 0, 0}[frame%10]
	bitmap = shifted(bitmap, dx, dy)

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

	shimmer := [][2]int{
		{15, 14}, {18, 13}, {21, 12}, {24, 11},
		{27, 10}, {24, 13}, {21, 14}, {18, 15},
	}

	sx, sy := shimmer[frame%n][0]+dx, shimmer[frame%n][1]+dy
	putIfPainted(bitmap, sx, sy, 's')

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
const mascotNative = 24
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
var (
	BobaBear    = mirror(bearHalf)
	BobaBearBig = mirror(resizeNearest(bearHalf, 18, 36))
	BobaCat     = mirror(catHalf)
	BobaCatBig  = mirror(resizeNearest(catHalf, 18, 36))
)

var BearPalette = Palette{
	'k': "#2B1B17",
	'f': "#A9683F",
	'h': "#D9A56C", 
	'm': "#E8C18A", 
	'n': "#241714", 
	'e': "#241714", 
	's': "#FFF1D2",
}

var CatPalette = Palette{
	'k': "#35252A", 
	'f': "#F0D7C9", 
	'p': "#D98D9C", 
	'h': "#E7BBA9", 
	'm': "#F4E4DB", 
	'n': "#6D3541", 
	'e': "#35252A", 
	's': "#FFF7F2", 
}


var (
	bearGlints = [][2]int{{7, 6}, {10, 5}, {13, 6}, {16, 7}, {13, 8}, {10, 7}}
	catGlints  = [][2]int{{7, 6}, {9, 6}, {11, 6}, {13, 7}, {11, 7}, {9, 7}}
)

func mascotFrame(sprite []string, p Palette, glints [][2]int, frame int) string {
	dy := []int{0, 0, 1, 1, 0, 0, -1, -1}[frame%8]

	bitmap := shifted(sprite, 0, dy)

	size := len(sprite)
	g := glints[frame%len(glints)]
	putIfPainted(bitmap, g[0]*size/mascotNative, g[1]*size/mascotNative+dy, 's')

	return Render(bitmap, p)
}

func BearFrame(frame int) string {
	return mascotFrame(BobaBearBig, BearPalette, bearGlints, frame)
}


func BearSmallFrame(frame int) string {
	return mascotFrame(BobaBear, BearPalette, bearGlints, frame)
}

func CatFrame(frame int) string {
	return mascotFrame(BobaCatBig, CatPalette, catGlints, frame)
}

func CatSmallFrame(frame int) string {
	return mascotFrame(BobaCat, CatPalette, catGlints, frame)
}
