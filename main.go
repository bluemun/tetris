// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris main.go Starts the program.
package main

import (
	"github.com/bluemun/graphics"
	"os"
	"time"
)

var stopped = false

func main() {
	go loop()

	graphics.Loop()
}

func loop() {
	window := graphics.CreateWindow()
	renderer := window.GetRenderer()
	g := CreateGrid(18, 10)
	render := time.NewTicker(time.Second / 60)
	update := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-render.C:
			renderer.Begin()
			g.Render(renderer)
			renderer.Flush()
			renderer.End()

			window.SwapBuffers()
		case <-update.C:
			window.PollEvents()
			if window.Closed() {
				os.Exit(0)
			}
		}
	}
}
