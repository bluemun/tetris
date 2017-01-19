// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris block.go Defines blocks and pieces used by the grid.
package main

import (
	"github.com/bluemun/engine/graphics/render"
)

type block struct {
	X, Y int
}

type piece struct {
	X, Y int
	g    *Grid
	b    [4]*block
}

// CreatePiece creates a piece to be used with the grid.
func createPiece(g *Grid) *piece {
	p := new(piece)
	p.g = g
	p.b[0] = &block{X: 0, Y: 1}
	p.b[1] = &block{X: 0, Y: 0}
	p.b[2] = &block{X: 0, Y: -1}
	p.b[3] = &block{X: 1, Y: -1}
	return p
}

func (p *piece) SetPosition(x, y int) {
	p.X = x
	p.Y = y
}

func (p *piece) TryMove(x, y int) bool {
	for _, pb := range p.b {
		if p.X+pb.X+x < 0 || p.X+pb.X+x >= p.g.columns || p.Y+pb.Y+y < 0 || p.Y+pb.Y+y >= p.g.rows {
			return false
		}

		var cell = p.g.data[p.Y+pb.Y+y][p.X+pb.X+x]
		if cell == nil {
			continue
		}

		for _, nb := range p.b {
			if cell == nb {
				continue
			}
		}

		return false
	}

	p.X += x
	p.Y += y
	return true
}

func (p *piece) TryRotate() {

}

func (p *piece) Integrate() {
	for i, pb := range p.b {
		pb.X += p.X
		pb.Y += p.Y
		p.g.IntegrateBlock(pb)
		p.b[i] = nil
	}

	p.g = nil
}

// Render Renders the whole grid.
func (p *piece) Render(r render.Renderer) {
	for _, block := range p.b {
		//r.DrawRectangle(float32(p.X+block.X), float32(p.Y+block.Y), 1, 1)
		r.DrawRectangle(float32(p.X+block.X)+0.1, float32(p.Y+block.Y)+0.1, 1-0.1, 1-0.1)
	}
}
