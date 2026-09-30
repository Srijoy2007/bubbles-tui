// cmd/logotest/main.go
package main

import (
	"fmt"

	"github.com/Srijoy2007/bubbles-tui/internal/art"
)

func main() {
	fmt.Print(art.Render(art.BobaCup, art.BobaPalette))
}
