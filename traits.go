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
	world          munfall.World
	owner          munfall.Actor
	tick, ticktime float32

	blocked, rush           bool
	blockedOrder, rushorder string
	move                    *munfall.WPos
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *MoveTickTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
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
		// Get the path from the actor position to the actor position offset by our move vector.
		path := t.world.WorldMap().GetPath(t.owner, t.owner.Pos(), t.owner.Pos().Add(t.move))
		for _, canmove := range t.world.GetTraitsImplementing(t.owner, (*Mover)(nil)) {
			// Check if all traits on this actor agree that we can move to the WPos in the path.
			// path.WPos(percent) lerps the position between two points in the path according
			// to the percentage given, if it is the end of the path it just returns the position.
			if !canmove.(Mover).CanMove(path.WPos(1)) {
				t.blocked = true
				break
			}
		}

		if t.blocked || munfall.IsNil(path.Next()) || path.IsEnd() {
			// We can't move as the path is blocked (path.Next() is nil) so we destroy this actor.
			t.blocked = true
			t.world.IssueGlobalOrder(&munfall.Order{Order: t.blockedOrder})
			t.world.AddFrameEndTask(func() {
				t.owner.Kill()
				theGame.ActorRegistry().DisposeActor(t.owner, t.world)
			})
		} else {
			// Move according to the path we generated, this sends a NotifyMove call to every
			// trait on the actor that implements it.
			t.world.WorldMap().Move(t.owner, path.Last(), 1)
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

// Mover is used by traits to provide a simple question interface for movement.
type Mover interface {
	CanMove(*munfall.WPos) bool
}

// SpawnActorFollowerTrait spawns an actor at a specific offset to this actor
// when it is added to the world.
type SpawnActorFollowerTrait struct {
	world munfall.World
	owner munfall.Actor

	offset *munfall.WPos
	actor  munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *SpawnActorFollowerTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.offset = parameters["Offset"].(*munfall.WPos)
	t.actor = theGame.ActorRegistry().CreateActor(parameters["Actor"].(string), nil, world, false)
}

// Owner returns the actor this trait is attached on.
func (t *SpawnActorFollowerTrait) Owner() munfall.Actor {
	return t.owner
}

// NotifyAddedToWorld adds the actor to the world that is supposed to be spawned.
func (t *SpawnActorFollowerTrait) NotifyAddedToWorld() {
	t.actor.SetPos(t.owner.Pos().Add(t.offset))
	t.world.AddToWorld(t.actor)
}

// NotifyMove called when the actor moves.
func (t *SpawnActorFollowerTrait) NotifyMove(old, new *munfall.WPos) {
	path := t.world.WorldMap().GetPath(t.actor, t.actor.Pos(), t.actor.Pos().Add(new.Subtract(old)))
	t.world.WorldMap().Move(t.actor, path.Last(), 1)
}

// CanMove tests if this position is legal.
func (t *SpawnActorFollowerTrait) CanMove(pos *munfall.WPos) bool {
	path := t.world.WorldMap().GetPath(t.actor, t.actor.Pos(), t.actor.Pos().Add(t.owner.Pos().Add(t.offset).Subtract(pos)))
	return !munfall.IsNil(path.Next())
}

// SpawnActorOrderTrait spawns an actor at a specific spawnpoint when it recieves a specific order.
type SpawnActorOrderTrait struct {
	world munfall.World
	owner munfall.Actor

	spawnpoint   *munfall.WPos
	actors       []string
	order        string
	currentActor munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *SpawnActorOrderTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.spawnpoint = parameters["SpawnPoint"].(*munfall.WPos)
	t.actors = parameters["Actors"].([]string)
	t.order = parameters["Order"].(string)
}

// Owner returns the actor this trait is attached on.
func (t *SpawnActorOrderTrait) Owner() munfall.Actor {
	return t.owner
}

// ResolveOrder spawns a specified actor when a specified order is resolved.
func (t *SpawnActorOrderTrait) ResolveOrder(order *munfall.Order) {
	if order.Order == t.order {
		a := t.actors[rand.Intn(len(t.actors))]
		w := t.world
		w.AddFrameEndTask(func() {
			if t.currentActor != nil {
				t.currentActor.Kill()
				theGame.ActorRegistry().DisposeActor(t.currentActor, t.world)
			}

			t.currentActor = theGame.ActorRegistry().CreateActor(a, nil, w, false)
			t.currentActor.SetPos(t.spawnpoint)
			w.AddToWorld(t.currentActor)
		})
	}
}

// ClearRowTrait clears full rows and moves cells down a block when it recives a specific order.
type ClearRowTrait struct {
	world munfall.World
	owner munfall.Actor

	order, moveorder string
	stepsize         float32
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *ClearRowTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner

	t.order = parameters["Order"].(string)
	t.moveorder = parameters["MoveOrder"].(string)
	t.stepsize = parameters["StepSize"].(float32)
}

// Owner returns the actor this trait is attached on.
func (t *ClearRowTrait) Owner() munfall.Actor {
	return t.owner
}

// ResolveOrder checks for full rows in the map, clears them and moves all the
// actors that contain the MoveControlTrait and are disabled down 1 cell if a full
// row was found.
func (t *ClearRowTrait) ResolveOrder(order *munfall.Order) {
	if order.Order == t.order {
		w := t.world
		w.AddFrameEndTask(func() {
			row := uint(0)
			for {
				occupied, size := t.checkRow(row)
				if occupied <= 0 {
					// We stop when we are outside the map or we hit an empty row.
					break
				}

				if occupied == size {
					t.clearRow(row)
					t.collapse(row)
				}

				row++
			}
		})
	}
}

// checkRow calculates the occupied cells in the row and the row size.
// returns -1 -1 if the row is outside the map.
func (t *ClearRowTrait) checkRow(row uint) (int, int) {
	pos := munfall.MPos{X: 0, Y: row}
	if !t.world.WorldMap().InsideMapMPos(&pos) {
		return -1, -1
	}

	nonemptycells := 0
	rowsize := 0
	for {
		if !t.world.WorldMap().InsideMapMPos(&pos) {
			break
		}

		cell := t.world.WorldMap().CellAt(&pos)
		if len(cell.Space()) != 0 {
			nonemptycells++
		}

		pos.X++
		rowsize++
	}

	return nonemptycells, rowsize
}

// clearRow iterates over a row and deletes all actors that take space in it.
func (t *ClearRowTrait) clearRow(row uint) {
	pos := munfall.MPos{X: 0, Y: row}
	for {
		if !t.world.WorldMap().InsideMapMPos(&pos) {
			break
		}

		cell := t.world.WorldMap().CellAt(&pos)
		if len(cell.Space()) != 0 {
			for _, space := range cell.Space() {
				actor := space.Trait().Owner()
				munfall.Logger.Info(actor)
				// We just kill every actor that we find in the space, we never do it twice
				// as by design we don't have an actor that takes more space then 1 cell.
				actor.Kill()
				theGame.ActorRegistry().DisposeActor(t.owner, t.world)
			}
		}

		pos.X++
	}
}

// collapse moves all the actors in the specified row down one cell,
// returns true if the operation succeeded.
func (t *ClearRowTrait) collapse(row uint) bool {
	pos := munfall.MPos{X: 0, Y: row}
	down := t.world.WorldMap().ConvertToWPos(&munfall.MPos{X: 1, Y: 0})
	if !t.world.WorldMap().InsideMapMPos(&pos) {
		return false
	}

	for {
		if !t.world.WorldMap().InsideMapMPos(&pos) {
			break
		}

		cell := t.world.WorldMap().CellAt(&pos)
		if len(cell.Space()) == 0 {
			for _, space := range cell.Space() {
				actor := space.Trait().Owner()
				t.world.WorldMap().GetPath(actor, actor.Pos(), actor.Pos().Subtract(down))
			}
		}

		pos.X++
	}

	return true
}
