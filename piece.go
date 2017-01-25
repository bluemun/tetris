// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris block.go Defines blocks and pieces used by the grid.
package main

import (
	"github.com/bluemun/engine"
	"github.com/bluemun/engine/graphics/render"
)

type orderpack struct {
	left, enabled bool
}

type piece struct {
	x, y           float32
	gravReset      float32
	gravCounter    float32
	inputCounter   float32
	mr, ml, rr, rl bool
	g              Grid
	blocks         [25]bool
	world          engine.World
	owner          engine.Actor
}

// CreatePiece creates a piece to be used with the grid.
func createPiece(g Grid) *piece {
	p := &piece{g: g}
	p.gravReset = 0.2
	p.blocks[2+3*5] = true
	p.blocks[2+2*5] = true
	p.blocks[2+1*5] = true
	p.blocks[3+1*5] = true
	return p
}

// NotifyAdded runs when the grid gets added to a world.
func (p *piece) NotifyAdded(owner engine.Actor) {
	c, r := p.g.Size()
	p.owner = owner
	p.world = owner.World()
	x, y := p.g.Pos()
	p.SetPosition(x+float32(c/2), y+float32(r-1))
	if !p.TryMove(0, 0) {
		p.world.RemoveActor(owner)
	}
}

// Tick runs when the world ticks.
func (p *piece) Tick(deltaUnit float32) {
	if p.gravCounter >= p.gravReset {
		p.gravCounter = 0
		if !p.TryMove(0, -1) {
			p.Integrate()
			p.g.spawnPiece()
			p.world.AddFrameEndTask(func() {
				p.world.RemoveActor(p.owner)
			})
		}
	}

	if p.inputCounter <= 0 && (p.mr || p.ml) {
		var x float32
		if p.mr {
			x++
		}
		if p.ml {
			x--
		}
		p.TryControlMove(x, 0)
	}

	p.gravCounter += deltaUnit
	p.inputCounter -= deltaUnit
}

// Render2D renders the grid.
func (p *piece) Render2D() []engine.Renderable {
	return []engine.Renderable{p}
}

func (p *piece) ResolveOrder(order *engine.Order) {
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
	}
}

func (p *piece) TryControlMove(x, y float32) bool {
	if p.inputCounter <= 0 && p.TryMove(x, y) {
		p.inputCounter = 0.01
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
}

func (p *piece) collides(array [25]bool, x, y float32) bool {
	c, r := p.g.Size()
	for i, exists := range array {
		if exists {
			bx, by := p.blockPos(i)
			bx += x
			by += y
			if bx < 0 || bx >= float32(c) || by < 0 {
				return true
			}

			if by >= float32(r) {
				continue
			}

			if p.g.GetBlock(int(bx), int(by)) {
				return true
			}
		}
	}

	return false
}

func (p *piece) Integrate() {
	for i, exists := range p.blocks {
		if exists {
			bx, by := p.blockPos(i)
			p.g.SetBlock(int(bx), int(by))
		}
	}
}

func (p *piece) Pos() (float32, float32) {
	return p.x, p.y
}

func (p *piece) Color() uint32 {
	return render.ToColor(255, 0xff, 0, 255)
}

func (p *piece) Mesh() *engine.Mesh {
	mesh := &engine.Mesh{}
	var offset uint32
	for i, exists := range p.blocks {
		if exists {
			x := float32(i%5 - 2)
			y := float32(i/5 - 2)
			mesh.Points = append(mesh.Points, x, y, 0, x+0.9, y, 0, x, y+0.9, 0, x+0.9, y+0.9, 0)
			mesh.Triangles = append(mesh.Triangles, offset+0, offset+1, offset+2, offset+1, offset+2, offset+3)
			offset += 4
		}
	}

	return mesh
}

func (p *piece) SetPosition(x, y float32) {
	p.x = x
	p.y = y
}

func (p *piece) blockPos(i int) (float32, float32) {
	x, y := p.g.Pos()
	return p.x - x + float32(i%5) - 2, p.y - y + float32(i/5) - 2
}
