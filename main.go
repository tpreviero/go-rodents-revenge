package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	game := newGame()

	ui := &UI{}
	ui.Init()
	defer ui.Close()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		game.Update()
		ui.Draw(game.Board)

		rl.EndDrawing()
	}
}

var keyPressedToMove = map[int32]*Move{
	rl.KeyUp:    {-1, 0},
	rl.KeyDown:  {1, 0},
	rl.KeyLeft:  {0, -1},
	rl.KeyRight: {0, 1},
}

func (g *Game) Update() {
	if g.GameState == Playing {
		g.Board.updateCats()
	}

	rodent := g.Board.findRodent()
	if rodent == nil {
		// the rodent has been eaten by a cat
		g.GameState = GameOver
	}

	if g.GameState == Playing {
		move, ok := keyPressedToMove[rl.GetKeyPressed()]
		if ok {
			moveRodent(g.Board, move)
		}
	}
}

func (b *Board) move(row, colum, dRow, dColumn int) bool {
	newRow, newColumn := row+dRow, colum+dColumn

	if b.Objects[newRow][newColumn] == Wall {
		return false
	}

	if b.Objects[newRow][newColumn] == Empty {
		b.Objects[newRow][newColumn] = b.Objects[row][colum]
		b.Objects[row][colum] = Empty
		return true
	}

	if b.move(newRow, newColumn, dRow, dColumn) {
		b.Objects[newRow][newColumn] = b.Objects[row][colum]
		b.Objects[row][colum] = Empty
		return true
	}

	return false
}
