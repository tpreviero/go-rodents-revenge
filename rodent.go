package main

import rl "github.com/gen2brain/raylib-go/raylib"

var keyToMove = map[int32]*Move{
	rl.KeyUp:    {-1, 0},
	rl.KeyKp8:   {-1, 0},
	rl.KeyDown:  {1, 0},
	rl.KeyKp2:   {1, 0},
	rl.KeyLeft:  {0, -1},
	rl.KeyKp4:   {0, -1},
	rl.KeyRight: {0, 1},
	rl.KeyKp6:   {0, 1},
	rl.KeyKp7:   {-1, -1},
	rl.KeyKp9:   {-1, 1},
	rl.KeyKp1:   {1, -1},
	rl.KeyKp3:   {1, 1},
}

func (g *Game) moveRodent() {
	move, ok := keyToMove[rl.GetKeyPressed()]
	if ok {
		g.move(g.Board.findRodent(), move)
	}
}

func (b *Board) findRodent() *Position {
	for x, row := range b.Objects {
		for y, obj := range row {
			if obj == Rodent || obj == RodentSinkHole {
				return &Position{Row: x, Column: y}
			}
		}
	}

	return nil
}

func (g *Game) respawnRodent() {
	for {
		x := rl.GetRandomValue(1, 22)
		y := rl.GetRandomValue(1, 22)
		position := &Position{int(x), int(y)}

		if g.Board.at(position) == Empty {
			g.Board.set(position, Rodent)
			break
		}
	}
}
