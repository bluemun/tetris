package main

type grid struct {
	activePiece   *piece
	data          [][]*block
	rows, columns int
}

func CreateGrid(rows, columns int) *grid {
	g := new(grid)

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

func (g *grid) SpawnPiece() bool {
	g.activePiece = CreatePiece(g)
	g.activePiece.SetPosition(g.columns/2, 1)
	return g.activePiece.TryMove(0, 0)
}

func (g *grid) Move(x, y int) {
	g.activePiece.TryMove(x, y)
}

func (g *grid) IntegrateBlock(b *block) {
	g.data[b.Y][b.X] = b
}

func (g *grid) Render() {
}

func (g *grid) String() string {
	var str string
	for y, row := range g.data {
		for x, cell := range row {
			var cb *block = cell
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
