// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris main.go Starts the program.
package main

import (
	"github.com/bluemun/engine"
	"github.com/bluemun/engine/game"
	"github.com/bluemun/engine/input"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var stopped = false

func main() {
	game := &game.Game{}
	game.Initialize()

	game.Camera.X = 0
	game.Camera.Y = 0
	game.Camera.Width = 20
	game.Camera.Height = 20

	og := input.CreateScriptableOrderGenerator()
	og.AddKeyScript(int(glfw.KeyA), "move", &move{x: 1, y: 0})
	game.SetOrderGenerator(og)

	game.World().CreateActor(func() engine.Trait {
		return CreateGrid(game.World(), 18, 10)
	})

	game.Start()
}
