// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris main.go Starts the program.
package main

import (
	"github.com/bluemun/engine"
	"github.com/bluemun/engine/game"
	"github.com/bluemun/engine/input"
)

var stopped = false

const framrate int64 = 120

func main() {
	game := &game.Game{}
	game.Initialize()

	game.Camera.X = 0
	game.Camera.Y = 0
	game.Camera.Width = 20
	game.Camera.Height = 20

	og := input.CreateScriptableOrderGenerator()
	og.AddKeyScript(39, true, "rush", nil)

	og.AddKeyScript(38, true, "move", &orderpack{left: true, enabled: true})
	og.AddKeyScript(38, false, "move", &orderpack{left: true, enabled: false})
	og.AddKeyScript(40, true, "move", &orderpack{left: false, enabled: true})
	og.AddKeyScript(40, false, "move", &orderpack{left: false, enabled: false})

	og.AddKeyScript(24, true, "rotate", &orderpack{left: true, enabled: true})
	og.AddKeyScript(24, false, "rotate", &orderpack{left: true, enabled: false})
	og.AddKeyScript(26, true, "rotate", &orderpack{left: false, enabled: true})
	og.AddKeyScript(26, false, "rotate", &orderpack{left: false, enabled: false})
	game.SetOrderGenerator(og)

	game.World().CreateActor(func() engine.Trait {
		return CreateGrid(game.World(), 18, 10)
	})

	game.Start(framrate)
}
