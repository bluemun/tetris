// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Package tetris render.go Defines a renderable that is used by traits to
// render actors in the world.
package main

import (
	"github.com/bluemun/munfall"
)

type squereRenderable struct {
	mesh  *munfall.Mesh
	pos   *munfall.WPos
	color uint32
}

// CreateSquereRenderable creates a squereRenderable that has the information to
// render a squere at the given position.
func CreateSquereRenderable(pos *munfall.WPos, color uint32, halfsize float32) munfall.Renderable {
	return &squereRenderable{
		mesh: &munfall.Mesh{
			Points: []float32{
				-halfsize, -halfsize, 0,
				-halfsize, halfsize, 0,
				halfsize, -halfsize, 0,
				halfsize, halfsize, 0,
			},
			Triangles: []uint32{
				0, 1, 2,
				1, 2, 3,
			},
		},
		pos:   pos,
		color: color,
	}
}

func (s *squereRenderable) Mesh() *munfall.Mesh {
	return s.mesh
}

func (s *squereRenderable) Pos() *munfall.WPos {
	return s.pos
}

func (s *squereRenderable) Color() uint32 {
	return s.color
}
