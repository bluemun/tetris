// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris main.go Starts the program.
package main

import (
	"os"
	"time"

	"github.com/bluemun/engine"
	"github.com/bluemun/engine/graphics"
	"github.com/bluemun/engine/graphics/render"
	"github.com/bluemun/engine/logic"
)

var stopped = false

func main() {
	go loop()

	engine.Loop()
}

func loop() {
	window := graphics.CreateWindow()
	camera := new(render.Camera)
	camera.X = 0
	camera.Y = 0
	camera.Width = 20
	camera.Height = 20
	camera.Activate()

	world := logic.CreateWorld()
	renderer := render.CreateRendersTraits2D(world)

	render := time.NewTicker(time.Second / 60)
	update := time.NewTicker(time.Second / 60)

	g := CreateGrid(18, 10)
	world.Traitmanager.AddTrait(g)

	for {
		select {
		case <-render.C:
			window.Clear()
			renderer.Render()
			window.SwapBuffers()
		case <-update.C:
			world.Tick(1 / 60.0)
			window.PollEvents()
			if window.Closed() {
				os.Exit(0)
			}
		}
	}
}
