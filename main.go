// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris main.go Starts the program.
package main

import (
	"github.com/bluemun/engine"
	"github.com/bluemun/engine/graphics"
	"github.com/bluemun/engine/graphics/render"
	"os"
	"time"
)

var stopped = false

func main() {
	go loop()

	engine.Loop()
}

func loop() {
	window := graphics.CreateWindow()
	renderer := render.CreateRenderer2D(10000, 10000)
	g := CreateGrid(18, 10)
	render := time.NewTicker(time.Second / 60)
	update := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-render.C:
			window.Clear()

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
