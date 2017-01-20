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
type Grid struct {
	activePiece   *piece
	data          []bool
	rows, columns int
	gravCounter   float32
}

// CreateGrid creates a grid object correctly.
func CreateGrid(rows, columns int) *Grid {
	g := new(Grid)
	g.data = make([]bool, columns+rows*columns, columns+rows*columns)
	g.rows = rows
	g.columns = columns

	return g
}

// SpawnPiece spawns a new piece at the top of the grid.
func (g *Grid) spawnPiece() bool {
	g.activePiece = createPiece(g)
	g.activePiece.SetPosition(float32(g.columns/2), float32(g.rows-1))
	return g.activePiece.TryMove(0, 0)
}

// Move moves the active piece the given vector if possible.
func (g *Grid) Move(x, y float32) {
	g.activePiece.TryMove(x, y)
}

// IntegrateBlock Adds a given blodk to the grid.
func (g *Grid) IntegrateBlock(x, y int) {
	g.data[x+y*g.columns] = true
}

// NotifyAdded runs when the grid gets added to a world.
func (g *Grid) NotifyAdded(world *logic.World) {
	g.spawnPiece()
}

// Mesh Renderable interface
func (g *Grid) Mesh() *render.Mesh {
	mesh := new(render.Mesh)
	/*for i, exists := range g.data {
		if exists {
			x := float32(i%5) - 2
			y := float32(i/5) - 2
			mesh.Points = []float32{
				g.data[i],
			}
		}
	}*/
	return mesh
}

// Pos Renderable interface
func (g *Grid) Pos() (float32, float32) {
	return 0, 0
}

// Color Renderable interface
func (g *Grid) Color() uint32 {
	return render.ToColor(255, 0, 0, 255)
}

// Tick runs when the world ticks.
func (g *Grid) Tick(deltaUnit float32) {
	if g.gravCounter >= 1 {
		g.gravCounter = 0
		/*if !g.activePiece.TryMove(0, -1) {
			g.activePiece.Integrate()
			g.spawnPiece()
		}*/
	}

	g.gravCounter += deltaUnit
}

// Render2D renders the grid.
func (g *Grid) Render2D() []render.Renderable {
	return []render.Renderable{g.activePiece}
}
