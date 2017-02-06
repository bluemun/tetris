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
	wmap := gridworldmap.CreateGridWorldMap(10, 20, 25, 25)

	theGame = &game.Game{}
	theGame.Initialize(wmap)

	theGame.Camera.X = 5
	theGame.Camera.Y = 9
	theGame.Camera.Width = 20
	theGame.Camera.Height = 20

	og := input.CreateScriptableOrderGenerator()
	/*og.AddKeyScript(39, true, "rush", nil)

	og.AddKeyScript(38, true, "move", &orderpack{left: true, enabled: true})
	og.AddKeyScript(38, false, "move", &orderpack{left: true, enabled: false})
	og.AddKeyScript(40, true, "move", &orderpack{left: false, enabled: true})
	og.AddKeyScript(40, false, "move", &orderpack{left: false, enabled: false})

	og.AddKeyScript(24, true, "rotate", &orderpack{left: true, enabled: true})
	og.AddKeyScript(24, false, "rotate", &orderpack{left: true, enabled: false})
	og.AddKeyScript(26, true, "rotate", &orderpack{left: false, enabled: true})
	og.AddKeyScript(26, false, "rotate", &orderpack{left: false, enabled: false})*/

	theGame.SetOrderGenerator(og)

	ar := theGame.ActorRegistry()
	ar.RegisterTrait("ActorSpawner", (*ActorSpawner)(nil))
	ar.RegisterTrait("CellBodyTrait", (*CellBodyTrait)(nil))
	ar.RegisterTrait("MoveOrderTrait", (*MoveOrderTrait)(nil))
	ar.RegisterTrait("MoveTickTrait", (*MoveTickTrait)(nil))
	ar.RegisterTrait("RenderCellBodyTrait", (*RenderCellBodyTrait)(nil))
	ar.RegisterTrait("RowClearer", (*CellBodyTrait)(nil))

	ad := logic.CreateActorDefinition("LPiece")
	ad.AddTrait("MoveOrderTrait")
	ad.AddParameter("MoveOrderTrait", "Order", "move")
	ad.AddTrait("MoveTickTrait")
	ad.AddParameter("MoveOrderTrait", "BlockedOrder", "finished")
	ad.AddParameter("MoveTickTrait", "TickTime", 0.5)
	ad.AddTrait("CellBodyTrait")
	ad.AddParameter("CellBodyTrait", "Size", 25)
	ad.AddParameter("CellBodyTrait", "Offsets", []*munfall.WPos{
		&munfall.WPos{X: 0, Y: 0, Z: 0},
	})
	ad.AddTrait("RenderCellBodyTrait")
	ad.AddParameter("RenderCellBodyTrait", "Color", 0xFF0000FF)
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("ReverseLPiece")
	ad.AddTrait("MoveOrderTrait")
	ad.AddParameter("MoveOrderTrait", "Order", "move")
	ad.AddTrait("MoveTickTrait")
	ad.AddParameter("MoveOrderTrait", "BlockedOrder", "finished")
	ad.AddParameter("MoveTickTrait", "TickTime", 0.5)
	ad.AddTrait("CellBodyTrait")
	ad.AddParameter("CellBodyTrait", "Size", 25)
	ad.AddParameter("CellBodyTrait", "Offsets", []*munfall.WPos{
		&munfall.WPos{X: 0, Y: 0, Z: 0},
	})
	ad.AddTrait("RenderCellBodyTrait")
	ad.AddParameter("RenderCellBodyTrait", "Color", 0xFF0000FF)
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("LinePiece")
	ad.AddTrait("MoveOrderTrait")
	ad.AddParameter("MoveOrderTrait", "Order", "move")
	ad.AddTrait("MoveTickTrait")
	ad.AddParameter("MoveOrderTrait", "BlockedOrder", "finished")
	ad.AddParameter("MoveTickTrait", "TickTime", 0.5)
	ad.AddTrait("CellBodyTrait")
	ad.AddParameter("CellBodyTrait", "Size", 25)
	ad.AddParameter("CellBodyTrait", "Offsets", []*munfall.WPos{
		&munfall.WPos{X: 0, Y: 0, Z: 0},
	})
	ad.AddTrait("RenderCellBodyTrait")
	ad.AddParameter("RenderCellBodyTrait", "Color", 0xFF0000FF)
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("SquerePiece")
	ad.AddTrait("MoveOrderTrait")
	ad.AddParameter("MoveOrderTrait", "Order", "move")
	ad.AddTrait("MoveTickTrait")
	ad.AddParameter("MoveOrderTrait", "BlockedOrder", "finished")
	ad.AddParameter("MoveTickTrait", "TickTime", 0.5)
	ad.AddTrait("CellBodyTrait")
	ad.AddParameter("CellBodyTrait", "Size", 25)
	ad.AddParameter("CellBodyTrait", "Offsets", []*munfall.WPos{
		&munfall.WPos{X: 0, Y: 0, Z: 0},
	})
	ad.AddTrait("RenderCellBodyTrait")
	ad.AddParameter("RenderCellBodyTrait", "Color", 0xFF0000FF)
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("SquiglyPiece")
	ad.AddTrait("MoveOrderTrait")
	ad.AddParameter("MoveOrderTrait", "Order", "move")
	ad.AddTrait("MoveTickTrait")
	ad.AddParameter("MoveOrderTrait", "BlockedOrder", "finished")
	ad.AddParameter("MoveTickTrait", "TickTime", 0.5)
	ad.AddTrait("CellBodyTrait")
	ad.AddParameter("CellBodyTrait", "Size", 25)
	ad.AddParameter("CellBodyTrait", "Offsets", []*munfall.WPos{
		&munfall.WPos{X: 0, Y: 0, Z: 0},
	})
	ad.AddTrait("RenderCellBodyTrait")
	ad.AddParameter("RenderCellBodyTrait", "Color", 0xFF0000FF)
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("ReverseSquiglyPiece")
	ad.AddTrait("MoveOrderTrait")
	ad.AddParameter("MoveOrderTrait", "Order", "move")
	ad.AddTrait("MoveTickTrait")
	ad.AddParameter("MoveOrderTrait", "BlockedOrder", "finished")
	ad.AddParameter("MoveTickTrait", "TickTime", 0.5)
	ad.AddTrait("CellBodyTrait")
	ad.AddParameter("CellBodyTrait", "Size", 25)
	ad.AddParameter("CellBodyTrait", "Offsets", []*munfall.WPos{
		&munfall.WPos{X: 0, Y: 0, Z: 0},
	})
	ad.AddTrait("RenderCellBodyTrait")
	ad.AddParameter("RenderCellBodyTrait", "Color", 0xFF0000FF)
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("TPiece")
	ad.AddTrait("MoveOrderTrait")
	ad.AddParameter("MoveOrderTrait", "Order", "move")
	ad.AddTrait("MoveTickTrait")
	ad.AddParameter("MoveOrderTrait", "BlockedOrder", "finished")
	ad.AddParameter("MoveTickTrait", "TickTime", 0.5)
	ad.AddTrait("CellBodyTrait")
	ad.AddParameter("CellBodyTrait", "Size", 25)
	ad.AddParameter("CellBodyTrait", "Offsets", []*munfall.WPos{
		&munfall.WPos{X: 0, Y: 0, Z: 0},
	})
	ad.AddTrait("RenderCellBodyTrait")
	ad.AddParameter("RenderCellBodyTrait", "Color", 0xFF0000FF)
	ar.RegisterActor(ad)

	ad = logic.CreateActorDefinition("Manager")
	ad.AddTrait("ActorSpawner")
	ad.AddParameter("ActorSpawner", "Actors", []string{"TPiece"})
	ad.AddParameter("ActorSpawner", "Order", "finished")
	ad.AddParameter("ActorSpawner", "SpawnPoint", &munfall.WPos{X: 5, Y: 18, Z: 0})
	ad.AddTrait("RowClearer")
	ad.AddParameter("RowClearer", "Order", "finished")
	ad.AddParameter("RowClearer", "StepSize", 25)
	ar.RegisterActor(ad)

	manager := ar.CreateActor("Manager", nil, theGame.World(), true)
	theGame.World().IssueOrder(manager, &munfall.Order{Order: "finished"})
	theGame.Start(framrate)
}
