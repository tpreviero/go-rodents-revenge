package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

func main() {
	game := NewGame()

	ui := NewUI()
	defer ui.Close()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		game.Update()
		ui.Draw(game)

		rl.EndDrawing()
	}
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

	if len(g.Board.findAllCats()) == 0 && g.Board.RemainingNumberOfWaves == 0 {
		g.NextLevel()
	}

	if g.GameState == Playing {
		g.moveRodent()
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

	if b.at(position) == Rodent && b.at(next) == SinkHole {
		b.set(position, Empty)
		b.set(next, SinkHoleRodent)
		b.InSinkHoleSince = time.Now()
		return false
	}

	if b.at(position) == SinkHoleRodent {
		if time.Now().Sub(b.InSinkHoleSince) >= config.SinkHoleDuration {
			b.set(position, Rodent)
		}
		return false
	}

	if b.at(position) == Obstacle && b.at(next) == SinkHole {
		return true
	}

	if b.at(next) == Wall {
		return false
	}

	if b.at(next) == Empty || b.at(next) == Cheese {
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
