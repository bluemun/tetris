// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris block.go Defines blocks and pieces used by the grid.
package main

import (
	"github.com/bluemun/engine/graphics/render"
	"github.com/bluemun/engine/logic"
)

type piece struct {
	x, y        float32
	gravCounter float32
	g           Grid
	blocks      [25]bool
	world       logic.World
	owner       logic.Actor
}

// CreatePiece creates a piece to be used with the grid.
func createPiece(g Grid) *piece {
	p := &piece{g: g}
	p.blocks[2+3*5] = true
	p.blocks[2+2*5] = true
	p.blocks[2+1*5] = true
	p.blocks[3+1*5] = true
	return p
}

// NotifyAdded runs when the grid gets added to a world.
func (p *piece) NotifyAdded(owner logic.Actor) {
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
	if p.gravCounter >= 0.1 {
		p.gravCounter = 0
		if !p.TryMove(0, -1) {
			p.Integrate()
			p.g.spawnPiece()
			p.world.AddFrameEndTask(func() {
				p.world.RemoveActor(p.owner)
			})
		}
	}

	p.gravCounter += deltaUnit
}

// Render2D renders the grid.
func (p *piece) Render2D() []render.Renderable {
	return []render.Renderable{p}
}

func (p *piece) Pos() (float32, float32) {
	return p.x, p.y
}

func (p *piece) Color() uint32 {
	return render.ToColor(255, 0xff, 0, 255)
}

func (p *piece) Mesh() *render.Mesh {
	mesh := new(render.Mesh)
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

func (p *piece) TryMove(x, y float32) bool {
	c, r := p.g.Size()
	for i, exists := range p.blocks {
		if exists {
			bx, by := p.blockPos(i)
			bx += x
			by += y
			if bx < 0 || bx >= float32(c) || by < 0 {
				return false
			}

			if by >= float32(r) {
				continue
			}

			if p.g.GetBlock(int(bx), int(by)) {
				return false
			}
		}
	}

	p.x += x
	p.y += y
	return true
}

func (p *piece) TryRotate() {

}

func (p *piece) Integrate() {
	for i, exists := range p.blocks {
		if exists {
			bx, by := p.blockPos(i)
			p.g.SetBlock(int(bx), int(by))
		}
	}
}
