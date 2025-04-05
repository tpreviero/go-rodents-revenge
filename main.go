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
		ui.Draw(game)

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
		if g.RamainingLives == 0 {
			g.GameState = GameOver
			return
		}

		g.RamainingLives--
		g.respawnRodent()
	}

	if g.GameState == Playing {
		move, ok := keyToMove[rl.GetKeyPressed()]
		if ok {
			g.move(rodent, move)
		}
	}
}

func (g *Game) move(position *Position, move *Move) bool {
	next := position.after(move)

	b := g.Board

	if b.at(position) == Rodent && b.at(next) == Cat {
		b.set(position, Empty)
		return true
	}

	if b.at(position) == Rodent && b.at(next) == Cheese {
		b.set(position, Empty)
		b.set(next, Rodent)
		g.Points += config.CheesePoints
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

	if g.move(next, move) {
		b.set(next, b.at(position))
		b.set(position, Empty)
		return true
	}

	return false
}
