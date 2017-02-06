// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris block.go Defines blocks and pieces used by the grid.
package main

import (
	"github.com/bluemun/munfall"
)

// MoveOrderTrait moves the actor when the given order is issued.
type MoveOrderTrait struct {
	Order string

	world munfall.World
	owner munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *MoveOrderTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner
}

// Owner returns the actor this trait is attached on.
func (t *MoveOrderTrait) Owner() munfall.Actor {
	return t.owner
}

// ResolveOrder moves the actor when it recieves a specific order.
func (t *MoveOrderTrait) ResolveOrder(order *munfall.Order) {
	if order.Order == "move" {
		pos := order.Value.(*munfall.WPos)
		t.owner.SetPos(pos)
	}
}

// MoveTickTrait bla
type MoveTickTrait struct {
	world munfall.World
	owner munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *MoveTickTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner
}

// Owner returns the actor this trait is attached on.
func (t *MoveTickTrait) Owner() munfall.Actor {
	return t.owner
}

// RenderCellBodyTrait bla
type RenderCellBodyTrait struct {
	world munfall.World
	owner munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *RenderCellBodyTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner
}

// Owner returns the actor this trait is attached on.
func (t *RenderCellBodyTrait) Owner() munfall.Actor {
	return t.owner
}

// CellBodyTrait bla
type CellBodyTrait struct {
	world munfall.World
	owner munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *CellBodyTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner
}

// Owner returns the actor this trait is attached on.
func (t *CellBodyTrait) Owner() munfall.Actor {
	return t.owner
}

// ActorSpawner bla
type ActorSpawner struct {
	world munfall.World
	owner munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *ActorSpawner) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner
}

// Owner returns the actor this trait is attached on.
func (t *ActorSpawner) Owner() munfall.Actor {
	return t.owner
}

// RowClearer bla
type RowClearer struct {
	world munfall.World
	owner munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (t *RowClearer) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	t.world = world
	t.owner = owner
}

// Owner returns the actor this trait is attached on.
func (t *RowClearer) Owner() munfall.Actor {
	return t.owner
}

/*
type orderpack struct {
	left, enabled bool
}

type PieceControlTrait struct {
	inputCounter   float32
	mr, ml, rr, rl bool
	world          munfall.World
	owner          munfall.Actor
}

// Initialize is used by the ActorRegistry to initialize this trait.
func (p *PieceControlTrait) Initialize(world munfall.World, owner munfall.Actor, parameters map[string]interface{}) {
	p.world = world
	p.owner = owner
}

// Owner returns the actor this trait is attached on.
func (p *PieceControlTrait) Owner() munfall.Actor {
	return p.owner
}

// ResolveOrder moves the actor when it recieves a specific order.
func (p *PieceControlTrait) ResolveOrder(order *munfall.Order) {
	switch order.Order {
	case "rush":
		p.gravReset = 0
	case "move":
		pack := order.Value.(*orderpack)
		switch pack.left {
		case true:
			p.ml = pack.enabled
			if pack.enabled {
				p.TryControlMove(-1, 0)
			}
		case false:
			p.mr = pack.enabled
			if pack.enabled {
				p.TryControlMove(1, 0)
			}
		}
	case "rotate":
		pack := order.Value.(*orderpack)
		switch pack.left {
		case true:
			p.rl = pack.enabled
			if pack.enabled {
				p.TryRotate(true)
			}
		case false:
			p.rr = pack.enabled
			if pack.enabled {
				p.TryRotate(false)
			}
		}

func (p *piece) TryControlMove(x, y float32) bool {
	if p.inputCounter <= 0 && p.TryMove(x, y) {
		p.inputCounter = 0.05
		return true
	}

	return false
}

func (p *piece) TryMove(x, y float32) bool {
	if p.collides(p.blocks, x, y) {
		return false
	}

	p.x += x
	p.y += y
	return true
}

func (p *piece) TryRotate(left bool) {
	var newBlocks [25]bool
	if left {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				newBlocks[x+y*5] = p.blocks[5-y-1+x*5]
			}
		}
	} else {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				newBlocks[x+y*5] = p.blocks[y+(5-x-1)*5]
			}
		}
	}

	if !p.collides(newBlocks, 0, 0) {
		p.blocks = newBlocks
	}
}*/
