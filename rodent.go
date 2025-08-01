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

var keyToMoveAnother = map[int32]*Move{
	rl.KeyW: {-1, 0},
	rl.KeyS: {1, 0},
	rl.KeyA: {0, -1},
	rl.KeyD: {0, 1},
	rl.KeyQ: {-1, -1},
	rl.KeyE: {-1, 1},
	rl.KeyZ: {1, -1},
	rl.KeyC: {1, 1},
}

func (g *Game) moveRodent() {
	resultingMove := &Move{0, 0}
	for key, move := range keyToMove {
		if rl.IsKeyPressed(key) {
			resultingMove = resultingMove.compose(move)
		}
	}

	if resultingMove.Row != 0 || resultingMove.Column != 0 {
		g.move(g.Board.findRodent(), resultingMove)
	}

	resultingMove = &Move{0, 0}
	for key, move := range keyToMoveAnother {
		if rl.IsKeyPressed(key) {
			resultingMove = resultingMove.compose(move)
		}
	}
	if resultingMove.Row != 0 || resultingMove.Column != 0 {
		g.move(g.Board.findAnotherRodent(), resultingMove)
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

func (b *Board) findAnotherRodent() *Position {
	for x, row := range b.Objects {
		for y, obj := range row {
			if obj == AnotherRodent || obj == AnotherRodentSinkHole {
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

func (g *Game) respawnAnotherRodent() {
	for {
		x := rl.GetRandomValue(1, 22)
		y := rl.GetRandomValue(1, 22)
		position := &Position{int(x), int(y)}

		if g.Board.at(position) == Empty {
			g.Board.set(position, AnotherRodent)
			break
		}
	}
}
