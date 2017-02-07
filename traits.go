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
	Order                    string
	stepSize, tick, ticktime float32
	right, left              bool

	world      munfall.World
	owner      munfall.Actor
	moveticker *MoveTickTrait
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *MoveOrderTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.Order = parameters["Order"].(string)
	t.stepSize = parameters["StepSize"].(float32)
	t.ticktime = parameters["MoveDelay"].(float32)
}

// NotifyAddedToWorld tries to retrieve a MoveTickTrait from the actor.
func (t *MoveOrderTrait) NotifyAddedToWorld() {
	t.moveticker = t.world.GetTrait(t.owner, (*MoveTickTrait)(nil)).(*MoveTickTrait)
}

// Owner returns the actor this trait is attached on.
func (t *MoveOrderTrait) Owner() munfall.Actor {
	return t.owner
}

// Tick moves the actor when the elapsed ticktime has passed.
func (t *MoveOrderTrait) Tick(du float32) {
	if !munfall.IsNil(t.moveticker) && t.moveticker.blocked {
		return
	}

	t.tick += du
	if t.ticktime <= t.tick {
		move := float32(0)
		if t.right {
			move += t.stepSize
		}

		if t.left {
			move -= t.stepSize
		}

		if move == 0 {
			t.tick = t.ticktime
			return
		}

		path := t.world.WorldMap().GetPath(t.owner, t.owner.Pos(),
			t.owner.Pos().Add(&munfall.WPos{X: move}))

		if !munfall.IsNil(path.Next()) {
			t.world.WorldMap().Move(t.owner, path.Last(), 1)
		}

		t.tick = 0
	}
}

// ResolveOrder moves the actor when it recieves a specific order.
func (t *MoveOrderTrait) ResolveOrder(order *munfall.Order) {
	if order.Order == t.Order {
		munfall.Logger.Info(order)
		switch order.Value.(int) {
		case 2:
			t.right = true
		case 1:
			t.right = false
		case -1:
			t.left = false
		case -2:
			t.left = true
		}
	}
}

// MoveTickTrait bla
type MoveTickTrait struct {
	worldMap       munfall.WorldMap
	owner          munfall.Actor
	tick, ticktime float32

	blocked, rush           bool
	blockedOrder, rushorder string
	move                    *munfall.WPos
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *MoveTickTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.worldMap = world.WorldMap()
	t.owner = owner
	t.ticktime = parameters["TickTime"].(float32)
	t.blockedOrder = parameters["BlockedOrder"].(string)
	t.rushorder = parameters["RushOrder"].(string)
	t.move = parameters["Move"].(*munfall.WPos)
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
	if t.ticktime <= t.tick {
		path := t.worldMap.GetPath(t.owner, t.owner.Pos(), t.owner.Pos().Add(t.move))

		if munfall.IsNil(path.Next()) {
			t.blocked = true

			cbt := t.owner.World().GetTrait(t.owner, (*CellBodyTrait)(nil)).(*CellBodyTrait)
			ys := make(map[uint]interface{}, 0)
			for _, space := range cbt.Space() {
				ys[t.worldMap.ConvertToMPos(space.Offset()).Y] = nil
			}

			order := &munfall.Order{
				Order: t.blockedOrder,
				Value: make([]uint, len(ys)),
			}

			i := 0
			for y := range ys {
				order.Value.([]uint)[i] = y
				i++
			}

			t.owner.World().IssueGlobalOrder(order)
		} else {
			t.worldMap.Move(t.owner, path.Last(), 1)
		}

		t.tick = 0
	}
}

// ResolveOrder moves the actor when it recieves a specific order.
func (t *MoveTickTrait) ResolveOrder(order *munfall.Order) {
	if order.Order == t.rushorder {
		t.ticktime = 0
	}
}

// CellBodyTrait bla
type CellBodyTrait struct {
	world munfall.World
	owner munfall.Actor
	space []munfall.Space

	HalfSize float32
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *CellBodyTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.HalfSize = parameters["HalfSize"].(float32)
	offsets := parameters["Offsets"].([]*munfall.WPos)
	t.space = make([]munfall.Space, len(offsets))
	for i, offset := range offsets {
		t.space[i] = &traits.SpaceCell{
			LocalOffset: &munfall.WPos{
				X: offset.X * t.HalfSize * 2,
				Y: offset.Y * t.HalfSize * 2,
				Z: offset.Z * t.HalfSize * 2,
			},
		}

		t.space[i].Initialize(t)
	}
}

// Owner returns the actor this trait is attached on.
func (t *CellBodyTrait) Owner() munfall.Actor {
	return t.owner
}

// OutOfBounds returns true if the actor is out of bound with the given offset.
func (t *CellBodyTrait) OutOfBounds(offset *munfall.WPos) bool {
	for _, space := range t.space {
		if !t.world.WorldMap().InsideMapWPos(space.Offset().Add(offset)) {
			return true
		}
	}

	return false
}

// Intersects returns true if the two spaces have at least 1 cell in common.
func (t *CellBodyTrait) Intersects(os traits.OccupySpace, offset *munfall.WPos) bool {
	if os.Owner().ActorID() == t.owner.ActorID() {
		return false
	}

	spaces := os.Space()
	for _, space := range spaces {
		for _, tspace := range t.space {
			if tspace.Intersects(space, offset) {
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
		mpos := t.world.WorldMap().ConvertToMPos(space.Offset())
		renderables[i] = CreateSquereRenderable(t.world.WorldMap().ConvertToWPos(mpos), 0xFF0000FF, cbt.HalfSize)
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

	order, moveorder string
	stepsize         float32
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *RowClearer) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.order = parameters["Order"].(string)
	t.moveorder = parameters["MoveOrder"].(string)
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
					pos.X = 0
					actors := make(map[uint]munfall.Actor)
					for {
						munfall.Logger.Info("Run", pos)
						if !wm.InsideMapMPos(pos) {
							munfall.Logger.Info("End")
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
						Order: t.moveorder,
						Value: -1,
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
