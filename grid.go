// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris grid.go Defines a grid struct used to store tetris blocks.
package main

import (
	"github.com/bluemun/engine/graphics/render"
)

// Grid holds blocks that are in the game.
type Grid struct {
	activePiece                *piece
	data                       [][]*block
	rows, columns, gravCounter int
}

// CreateGrid creates a grid object correctly.
func CreateGrid(rows, columns int) *Grid {
	g := new(Grid)

	for y := 0; y < rows; y++ {
		var row []*block
		for x := 0; x < columns; x++ {
			row = append(row, nil)
		}
		g.data = append(g.data, row)
	}

	g.rows = rows
	g.columns = columns

	return g
}

// SpawnPiece spawns a new piece at the top of the grid.
func (g *Grid) SpawnPiece() bool {
	g.activePiece = createPiece(g)
	g.activePiece.SetPosition(g.columns/2, g.rows-1)
	return g.activePiece.TryMove(0, 0)
}

// Move moves the active piece the given vector if possible.
func (g *Grid) Move(x, y int) {
	g.activePiece.TryMove(x, y)
}

// IntegrateBlock Adds a given blodk to the grid.
func (g *Grid) IntegrateBlock(b *block) {
	g.data[b.Y][b.X] = b
}

// Update Updates the whole grid.
func (g *Grid) Update() {
	if g.gravCounter == 30 {
		g.gravCounter = 0
		if !g.activePiece.TryMove(0, -1) {
			g.activePiece.Integrate()
			g.SpawnPiece()
		}
	}

	g.gravCounter++
}

// Render Renders the whole grid.
func (g *Grid) Render(r render.Renderer) {
	for _, row := range g.data {
		for _, cell := range row {
			if cell != nil {
				r.DrawRectangle(float32(cell.X)+0.1, float32(cell.Y)+0.1, 1-0.1, 1-0.1, render.ToColor(255, 0, 0, 255))
			}
		}
	}

	if g.activePiece != nil {
		g.activePiece.Render(r)
	}
}

func (g *Grid) String() string {
	var str string
	for y, row := range g.data {
		for x, cell := range row {
			var cb = cell
			if g.activePiece != nil {
				for _, ab := range g.activePiece.b {
					if ab != nil && ab.X == x && ab.Y == y {
						cb = ab
						break
					}
				}
			}

			if cb == nil {
				str += "o"
			} else {
				str += "x"
			}
		}
		str += "\n"
	}

	return str
}
