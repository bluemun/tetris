// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris main.go Starts the program.
package main

import (
	"github.com/bluemun/engine/game"
	"github.com/bluemun/engine/input"
	"github.com/bluemun/engine/logic"
)

var stopped = false

const framrate int64 = 120

var theGame *game.Game

func main() {
	theGame = &game.Game{}
	theGame.Initialize()

	theGame.Camera.X = 0
	theGame.Camera.Y = 0
	theGame.Camera.Width = 20
	theGame.Camera.Height = 20

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
	theGame.SetOrderGenerator(og)

	ar := theGame.ActorRegistry()
	ar.RegisterTrait("Grid", (*grid)(nil))
	ar.RegisterTrait("Piece", (*piece)(nil))

	ad := logic.CreateActorDefinition("Grid")
	ad.AddTrait("Grid")
	ad.AddParameter("Grid", "width", 10)
	ad.AddParameter("Grid", "height", 18)
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("Piece")
	ad.AddTrait("Piece")
	ar.RegisterActor(ad)

	ar.CreateActor("Grid", nil, theGame.World())

	theGame.Start(framrate)
}
