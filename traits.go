// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris block.go Defines blocks and pieces used by the grid.
package main

import (
	"math/rand"

	"github.com/bluemun/munfall"
	"github.com/bluemun/munfall/traits"
)

// MoveOrderTrait moves the actor when the given order is issued.
type MoveOrderTrait struct {
	Order string

	worldMap munfall.WorldMap
	owner    munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *MoveOrderTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.worldMap = world.WorldMap()
	t.owner = owner

	t.Order = parameters["Order"].(string)
}

// Owner returns the actor this trait is attached on.
func (t *MoveOrderTrait) Owner() munfall.Actor {
	return t.owner
}

// ResolveOrder moves the actor when it recieves a specific order.
func (t *MoveOrderTrait) ResolveOrder(order *munfall.Order) {
	if order.Order == t.Order {
		pos1 := t.worldMap.ConvertToMPos(t.owner.Pos())
		pos2 := t.worldMap.ConvertToMPos(t.owner.Pos().Add(order.Value.(*munfall.WPos)))
		path := t.worldMap.GetPath(t.owner, pos1, pos2).Last()
		t.worldMap.Move(t.owner, path, 1)
	}
}

// MoveTickTrait bla
type MoveTickTrait struct {
	worldMap       munfall.WorldMap
	owner          munfall.Actor
	tick, ticktime float32

	blocked      bool
	blockedOrder string
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *MoveTickTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.worldMap = world.WorldMap()
	t.owner = owner
	t.ticktime = parameters["TickTime"].(float32)
	t.blockedOrder = parameters["BlockedOrder"].(string)
}

// Owner returns the actor this trait is attached on.
func (t *MoveTickTrait) Owner() munfall.Actor {
	return t.owner
}

// Tick moves the actor when the elapsed ticktime has passed.
func (t *MoveTickTrait) Tick(du float32) {
	if t.blocked {
		return
	}

	t.tick += du
	if t.ticktime >= t.tick {
		pos1 := t.worldMap.ConvertToMPos(t.owner.Pos())
		pos2 := t.worldMap.ConvertToMPos(t.owner.Pos().Add(&munfall.WPos{X: 0, Y: 1}))
		path := t.worldMap.GetPath(t.owner, pos1, pos2).Last()

		if path.Next() == nil {
			t.blocked = true
			t.owner.World().IssueGlobalOrder(&munfall.Order{Order: t.blockedOrder})
		} else {
			t.worldMap.Move(t.owner, path, 1)
		}
	}
}

// CellBodyTrait bla
type CellBodyTrait struct {
	world munfall.World
	owner munfall.Actor
	space []munfall.Space

	Size float32
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *CellBodyTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.Size = parameters["Size"].(float32)
	offsets := parameters["Offsets"].([]*munfall.WPos)
	t.space = make([]munfall.Space, len(offsets))
	for i, offset := range offsets {
		t.space[i] = &traits.SpaceCell{
			LocalOffset: &munfall.WPos{
				X: offset.X * t.Size,
				Y: offset.Y * t.Size,
				Z: offset.Z * t.Size,
			},
		}

		t.space[i].Initialize(t)
	}
}

// Owner returns the actor this trait is attached on.
func (t *CellBodyTrait) Owner() munfall.Actor {
	return t.owner
}

// Intersects returns true if the two spaces have at least 1 cell in common.
func (t *CellBodyTrait) Intersects(os traits.OccupySpace) bool {
	if os.Owner().ActorID() == t.owner.ActorID() {
		return false
	}

	spaces := os.Space()
	for _, space := range spaces {
		for _, offset := range t.space {
			if *offset.Offset() == *space.Offset() {
				return true
			}
		}
	}

	return false
}

// Space returns the spaces behind this Trait.
func (t *CellBodyTrait) Space() []munfall.Space {
	return t.space
}

// RenderCellBodyTrait bla
type RenderCellBodyTrait struct {
	world munfall.World
	owner munfall.Actor

	color uint32
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *RenderCellBodyTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.color = parameters["Color"].(uint32)
}

// Owner returns the actor this trait is attached on.
func (t *RenderCellBodyTrait) Owner() munfall.Actor {
	return t.owner
}

// Render2D renders the actor according to the CellBodyTrait.
func (t *RenderCellBodyTrait) Render2D() []munfall.Renderable {
	cbt := t.world.GetTrait(t.owner, (*CellBodyTrait)(nil)).(*CellBodyTrait)
	spaces := cbt.Space()
	renderables := make([]munfall.Renderable, len(spaces))

	for i, space := range spaces {
		renderables[i] = CreateSquereRenderable(space.Offset(), 0xFF0000FF, cbt.Size)
	}

	return renderables
}

// ActorSpawner spawns an actor at a specific spawnpoint when it recieves a specific order.
type ActorSpawner struct {
	world munfall.World
	owner munfall.Actor

	spawnpoint *munfall.WPos
	actors     []string
	order      string
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *ActorSpawner) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.spawnpoint = parameters["SpawnPoint"].(*munfall.WPos)
	t.actors = parameters["Actors"].([]string)
	t.order = parameters["Order"].(string)
}

// Owner returns the actor this trait is attached on.
func (t *ActorSpawner) Owner() munfall.Actor {
	return t.owner
}

// ResolveOrder spawns a specified actor when a specified order is resolved.
func (t *ActorSpawner) ResolveOrder(order *munfall.Order) {
	if order.Order == t.order {
		a := t.actors[rand.Intn(len(t.actors))]
		w := t.world
		w.AddFrameEndTask(func() {
			a := theGame.ActorRegistry().CreateActor(a, nil, w, false)
			a.SetPos(t.spawnpoint)
			w.AddToWorld(a)
		})
	}
}

// RowClearer clears full rows and moves cells down a block when it recives a specific order.
type RowClearer struct {
	world munfall.World
	owner munfall.Actor

	order    string
	stepsize float32
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *RowClearer) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.order = parameters["Order"].(string)
	t.stepsize = parameters["StepSize"].(float32)
}

// Owner returns the actor this trait is attached on.
func (t *RowClearer) Owner() munfall.Actor {
	return t.owner
}

// ResolveOrder checks for full rows in the map, clears them and moves all the
// actors that contain the MoveControlTrait and are disabled down 1 cell if a full
// row was found.
func (t *RowClearer) ResolveOrder(order *munfall.Order) {
	if order.Order == t.order {
		rows := order.Value.([]uint)
		w := t.world
		w.AddFrameEndTask(func() {
			wm := w.WorldMap()

			pos := &munfall.MPos{}
			for _, row := range rows {
				pos.X = 0
				pos.Y = row

				full := true
				rowsize := 0
				for {
					if !wm.InsideMapMPos(pos) {
						break
					}

					cell := wm.CellAt(pos)
					if len(cell.Space()) == 0 {
						full = false
						break
					}

					pos.X++
					rowsize++
				}

				if full {
					actors := make(map[uint]munfall.Actor)
					for {
						if !wm.InsideMapMPos(pos) {
							break
						}

						for _, space := range wm.CellAt(pos).Space() {
							actor := space.Trait().Owner()
							actors[actor.ActorID()] = actor
						}

						pos.X++
					}

					for _, actor := range actors {
						wm.Deregister(actor)
					}

					w.AddFrameEndTask(func() {
						for _, actor := range actors {
							theGame.ActorRegistry().DisposeActor(actor, w)
						}
					})

					actors = make(map[uint]munfall.Actor)
				RowClearerOuter:
					for {
						empty := false
						for {
							// If we are out of X bounds.
							if !wm.InsideMapMPos(pos) {
								pos.X = 0
								pos.Y++
								// If we are out of Y bounds.
								if !wm.InsideMapMPos(pos) {
									break RowClearerOuter
								}

								break
							}

							// Get all actors in this cell.
							cell := wm.CellAt(pos)
							if len(cell.Space()) != 0 {
								empty = false
								for _, space := range cell.Space() {
									actor := space.Trait().Owner()
									actors[actor.ActorID()] = actor
								}
							}

							pos.X++
						}

						// We hit an empty row so we can exit early.
						if empty {
							break
						}
					}

					order := &munfall.Order{
						Order: "move",
						Value: &munfall.WPos{X: 0, Y: -t.stepsize, Z: 0},
					}

					// Issue a move order to every actor above the row on the y axis.
					for _, actor := range actors {
						w.IssueOrder(actor, order)
					}
				}
			}
		})
	}
}
