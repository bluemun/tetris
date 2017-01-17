package main

type block struct {
	X, Y int
}

type piece struct {
	g *grid
	b [4]*block
}

func CreatePiece(g *grid) *piece {
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

		var cell *block = p.g.data[pb.Y+y][pb.X+x]
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
