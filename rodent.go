package main

import rl "github.com/gen2brain/raylib-go/raylib"

func (b *Board) findRodent() *Position {
	for x, row := range b.Objects {
		for y, obj := range row {
			if obj == Rodent {
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
