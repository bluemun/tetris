// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris grid.go Defines a grid struct used to store tetris blocks.
package main

import (
	"github.com/bluemun/engine/graphics/render"
	"github.com/bluemun/engine/logic"
)

// Grid holds blocks that are in the game.
type Grid interface {
	logic.TraitNotifyAdded
	render.TraitRender2D
	render.Renderable

	Size() (int, int)
	spawnPiece()

	GetBlock(x, y int) bool
	SetBlock(x, y int)
}

type grid struct {
	x, y          float32
	data          []bool
	rows, columns int
	world         logic.World
}

// CreateGrid creates a grid object correctly.
func CreateGrid(world logic.World, rows, columns int) Grid {
	return &grid{
		x:       -float32(columns) / 2,
		y:       -float32(rows) / 2,
		data:    make([]bool, columns+rows*columns, columns+rows*columns),
		rows:    rows,
		columns: columns,
		world:   world,
	}
}

// SpawnPiece spawns a new piece at the top of the grid.
func (g *grid) spawnPiece() {
	g.world.CreateActor(func() logic.Trait {
		return createPiece(g)
	})
}

func (g *grid) Size() (int, int) {
	return g.columns, g.rows
}

func (g *grid) GetBlock(x, y int) bool {
	return g.data[x+y*g.columns]
}

// IntegrateBlock Adds a given blodk to the grid.
func (g *grid) SetBlock(x, y int) {
	g.data[x+y*g.columns] = true
}

// NotifyAdded runs when the grid gets added to a world.
func (g *grid) NotifyAdded(owner logic.Actor) {
	g.world = owner.World()
	g.spawnPiece()
}

// Mesh Renderable interface
func (g *grid) Mesh() *render.Mesh {
	mesh := &render.Mesh{}
	c, _ := g.Size()
	var offset uint32
	for i, exists := range g.data {
		if exists {
			x := float32(i % c)
			y := float32(i / c)
			mesh.Points = append(mesh.Points,
				x, y, 0,
				x+0.9, y, 0,
				x, y+0.9, 0,
				x+0.9, y+0.9, 0,
			)
			mesh.Triangles = append(mesh.Triangles,
				offset+0, offset+1, offset+2,
				offset+1, offset+2, offset+3,
			)
			offset += 4
		}
	}

	return mesh
}

// Pos Renderable interface
func (g *grid) Pos() (float32, float32) {
	return g.x, g.y
}

// Color Renderable interface
func (g *grid) Color() uint32 {
	return render.ToColor(255, 0, 0, 255)
}

// Render2D renders the grid.
func (g *grid) Render2D() []render.Renderable {
	return []render.Renderable{g}
}
