package board 

import (
	"testing"
	"github.com/Srijoy2007/bubbles-tui/internal/store"
)

func TestRenderSmoke(t *testing.T) {
	b := []store.Block{{Date: "2026-09-29", Start: 540, End: 600, Done: true}}
	out := Render(b, 5, 1000) // now = 16:40, well after this block
	t.Log(out)
}
