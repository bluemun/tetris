// Tetris project main.go
package main

import (
	"time"

	"github.com/bluemoon-graphics"
)

var stopped bool = false

func main() {
	window := graphics.CreateWindow()
	renderer := window.GetRenderer()
	g := CreateGrid(18, 10)

	render := time.NewTicker(time.Second / 60)
	update := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-render.C:
			g.String()
		case <-update.C:
		}
	}
}
