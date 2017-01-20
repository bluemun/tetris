// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris block.go Defines blocks and pieces used by the grid.
package main

import (
	"github.com/bluemun/engine/graphics/render"
)

type piece struct {
	x, y   float32
	g      *Grid
	blocks [25]bool
}

// CreatePiece creates a piece to be used with the grid.
func createPiece(g *Grid) *piece {
	p := new(piece)
	p.g = g
	p.blocks[2+1*5] = true
	p.blocks[2+2*5] = true
	p.blocks[2+3*5] = true
	p.blocks[3+3*5] = true
	return p
}

func (p *piece) Pos() (float32, float32) {
	return float32(p.x), float32(p.y)
}

func (p *piece) Color() uint32 {
	return render.ToColor(255, 0, 0, 255)
}

func (p *piece) Mesh() *render.Mesh {
	mesh := new(render.Mesh)
	mesh.Points = []float32{}
	mesh.Triangles = []uint32{}
	var j uint32
	for i, exists := range p.blocks {
		if exists {
			x := float32(i%5) - 2
			y := float32(i/5) - 2
			mesh.Points = append(mesh.Points, x, y, 0, x+0.9, y, 0, x, y+0.9, 0, x+0.9, y+0.9, 0)
			mesh.Triangles = append(mesh.Triangles, j+0, j+1, j+2, j+1, j+2, j+3)
			j += 4
		}
	}

	return mesh
}

func (p *piece) SetPosition(x, y float32) {
	p.x = x
	p.y = y
}

func (p *piece) blockPos(i int) (float32, float32) {
	return float32(p.x + float32(i%5) - 2), float32(p.y + float32(i/5) - 2)
}

func (p *piece) TryMove(x, y float32) bool {
	for i := range p.blocks {
		bx, by := p.blockPos(i)
		if bx+x < 0 || bx+x >= float32(p.g.columns) || by+y < 0 {
			return false
		}

		if by+y >= float32(p.g.rows) {
			continue
		}

		if p.g.data[int(x)+int(y)*p.g.columns] {
			return false
		}
	}

	p.x += x
	p.y += y
	return true
}

func (p *piece) TryRotate() {

}

func (p *piece) Integrate() {
	for i := range p.blocks {
		bx, by := p.blockPos(i)
		p.g.IntegrateBlock(int(bx), int(by))
	}

	p.g = nil
}
