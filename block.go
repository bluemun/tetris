// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris block.go Defines blocks and pieces used by the grid.
package main

type block struct {
	X, Y int
}

type piece struct {
	g *Grid
	b [4]*block
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
	for _, pb := range p.b {
		pb.X = x
		pb.Y = y
	}
}

func (p *piece) TryMove(x, y int) bool {
	for _, pb := range p.b {
		if pb.X+x < 0 || pb.X+x >= p.g.rows || pb.Y+y < 0 || pb.Y+y >= p.g.columns {
			return false
		}

		var cell = p.g.data[pb.Y+y][pb.X+x]
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

	for _, pb := range p.b {
		pb.X += x
		pb.Y += y
	}

	return true
}

func (p *piece) TryRotate() {

}

func (p *piece) Integrate() {
	for i, pb := range p.b {
		p.g.IntegrateBlock(pb)
		p.b[i] = nil
	}

	p.g = nil
}
