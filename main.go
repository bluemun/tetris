// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris main.go Starts the program.
package main

import (
	"github.com/bluemun/munfall"
	"github.com/bluemun/munfall/game"
	"github.com/bluemun/munfall/gridworldmap"
	"github.com/bluemun/munfall/input"
	"github.com/bluemun/munfall/logic"
)

var stopped = false

const framrate int64 = 120

var theGame *game.Game

func main() {
	halfsize := float32(25)
	size := halfsize * 2
	gravTime := float32(0.5)
	gravMove := &munfall.WPos{X: 0, Y: -size, Z: 0}
	moveDelay := float32(0.08)

	wmap := gridworldmap.CreateGridWorldMap(10, 21, size, size)

	theGame = &game.Game{}
	theGame.Initialize(wmap)

	theGame.Camera.X = 5 * size
	theGame.Camera.Y = 9 * size
	theGame.Camera.Width = 20 * size
	theGame.Camera.Height = 20 * size

	og := input.CreateScriptableOrderGenerator()
	og.AddKeyScript(65, true, "rush", nil)

	og.AddKeyScript(38, true, "move", -2)
	og.AddKeyScript(38, false, "move", -1)
	og.AddKeyScript(40, true, "move", 2)
	og.AddKeyScript(40, false, "move", 1)

	/*og.AddKeyScript(24, true, "rotate", &orderpack{left: true, enabled: true})
	og.AddKeyScript(24, false, "rotate", &orderpack{left: true, enabled: false})
	og.AddKeyScript(26, true, "rotate", &orderpack{left: false, enabled: true})
	og.AddKeyScript(26, false, "rotate", &orderpack{left: false, enabled: false})*/

	theGame.SetOrderGenerator(og)

	ar := theGame.ActorRegistry()
	ar.RegisterTrait("SpawnActorFollowerTrait", (*SpawnActorFollowerTrait)(nil))
	ar.RegisterTrait("SpawnActorOrderTrait", (*SpawnActorOrderTrait)(nil))
	ar.RegisterTrait("CellBodyTrait", (*CellBodyTrait)(nil))
	ar.RegisterTrait("MoveOrderTrait", (*MoveOrderTrait)(nil))
	ar.RegisterTrait("MoveTickTrait", (*MoveTickTrait)(nil))
	ar.RegisterTrait("RenderCellBodyTrait", (*RenderCellBodyTrait)(nil))
	ar.RegisterTrait("ClearRowTrait", (*ClearRowTrait)(nil))

	ad := logic.CreateActorDefinition("Block")
	ad.AddTrait(logic.CreateTraitDefinition("CellBodyTrait").
		AddParameter("HalfSize", halfsize).
		AddParameter("Offsets", []*munfall.WPos{&munfall.WPos{X: 0, Y: 0, Z: 0}}))
	ad.AddTrait(logic.CreateTraitDefinition("RenderCellBodyTrait").
		AddParameter("Color", uint32(0xFF0000FF)))
	ar.RegisterActor(ad)

	moveordertrait := logic.CreateTraitDefinition("MoveOrderTrait").
		AddParameter("Order", "move").
		AddParameter("StepSize", size).
		AddParameter("MoveDelay", moveDelay)

	moveticktrait := logic.CreateTraitDefinition("MoveTickTrait").
		AddParameter("BlockedOrder", "finished").
		AddParameter("RushOrder", "rush").
		AddParameter("TickTime", gravTime).
		AddParameter("Move", gravMove)

	ad = logic.CreateActorDefinition("LPiece")
	ad.AddTrait(moveordertrait)
	ad.AddTrait(moveticktrait)
	ad.AddTrait(logic.CreateTraitDefinition("SpawnActorFollowerTrait").
		AddParameter("Actor", "Block").
		AddParameter("Offset", &munfall.WPos{}))
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("ReverseLPiece")
	ad.AddTrait(moveordertrait)
	ad.AddTrait(moveticktrait)
	ad.AddTrait(logic.CreateTraitDefinition("SpawnActorFollowerTrait").
		AddParameter("Actor", "Block").
		AddParameter("Offset", &munfall.WPos{}))
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("LinePiece")
	ad.AddTrait(moveordertrait)
	ad.AddTrait(moveticktrait)
	ad.AddTrait(logic.CreateTraitDefinition("SpawnActorFollowerTrait").
		AddParameter("Actor", "Block").
		AddParameter("Offset", &munfall.WPos{}))
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("SquerePiece")
	ad.AddTrait(moveordertrait)
	ad.AddTrait(moveticktrait)
	ad.AddTrait(logic.CreateTraitDefinition("SpawnActorFollowerTrait").
		AddParameter("Actor", "Block").
		AddParameter("Offset", &munfall.WPos{}))
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("SquiglyPiece")
	ad.AddTrait(moveordertrait)
	ad.AddTrait(moveticktrait)
	ad.AddTrait(logic.CreateTraitDefinition("SpawnActorFollowerTrait").
		AddParameter("Actor", "Block").
		AddParameter("Offset", &munfall.WPos{}))
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("ReverseSquiglyPiece")
	ad.AddTrait(moveordertrait)
	ad.AddTrait(moveticktrait)
	ad.AddTrait(logic.CreateTraitDefinition("SpawnActorFollowerTrait").
		AddParameter("Actor", "Block").
		AddParameter("Offset", &munfall.WPos{}))
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("TPiece")
	ad.AddTrait(moveordertrait)
	ad.AddTrait(moveticktrait)
	ad.AddTrait(logic.CreateTraitDefinition("SpawnActorFollowerTrait").
		AddParameter("Actor", "Block").
		AddParameter("Offset", &munfall.WPos{}))
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("Manager")
	ad.AddTrait(logic.CreateTraitDefinition("SpawnActorOrderTrait").
		AddParameter("Actors", []string{"TPiece", "ReverseSquiglyPiece", "SquiglyPiece", "LinePiece", "LPiece", "ReverseLPiece", "SquerePiece"}).
		AddParameter("Order", "finished").
		AddParameter("SpawnPoint", &munfall.WPos{X: 5 * size, Y: 18 * size, Z: 0}))
	ad.AddTrait(logic.CreateTraitDefinition("ClearRowTrait").
		AddParameter("MoveOrder", "move").
		AddParameter("Order", "finished").
		AddParameter("StepSize", size))
	ar.RegisterActor(ad)

	manager := ar.CreateActor("Manager", nil, theGame.World(), true)
	theGame.World().IssueOrder(manager, &munfall.Order{Order: "finished", Value: []uint{}})
	theGame.Start(framrate)
}
