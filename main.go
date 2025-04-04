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

var keyToMove = map[int32]*Move{
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
		move, ok := keyToMove[rl.GetKeyPressed()]
		if ok {
			g.Board.move(rodent, move)
		}
	}
}

func (b *Board) move(position *Position, move *Move) bool {
	next := position.after(move)

	if b.at(position) == Rodent && b.at(next) == Cat {
		b.set(position, Empty)
		return true
	}

	if b.at(next) == Wall {
		return false
	}

	if b.at(next) == Empty {
		b.set(next, b.at(position))
		b.set(position, Empty)
		return true
	}

	if b.move(next, move) {
		b.set(next, b.at(position))
		b.set(position, Empty)
		return true
	}

	return false
}
