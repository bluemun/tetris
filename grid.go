// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris grid.go Defines a grid struct used to store tetris blocks.
package main

import (
	"github.com/bluemun/graphics"
)

// Grid holds blocks that are in the game.
type Grid struct {
	activePiece   *piece
	data          [][]*block
	rows, columns int
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
	g.activePiece.SetPosition(g.columns/2, 1)
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

// Render Renders the whole grid.
func (g *Grid) Render(r *graphics.Renderer) {
	r.DrawRectangle(0, 0, 1, 1)
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
