package main

import "C"
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

		if rl.IsKeyDown(rl.KeyRightShift) && rl.IsKeyPressed(rl.KeySlash) {
			ui.showHelp = !ui.showHelp
		}

		game.Update()
		ui.Draw(game)

		rl.EndDrawing()
	}
}

func (g *Game) Update() {
	if rl.IsKeyPressed(rl.KeyP) {
		if g.GameState == Pause {
			g.GameState = Playing
		} else {
			g.GameState = Pause
		}
	}
	if rl.IsKeyDown(rl.KeyRightShift) && rl.IsKeyPressed(rl.KeyRight) {
		g.NextLevel()
		return
	}
	if rl.IsKeyDown(rl.KeyRightShift) && rl.IsKeyPressed(rl.KeyLeft) {
		g.PreviousLevel()
		return
	}
	if rl.IsKeyDown(rl.KeyRightShift) && rl.IsKeyPressed(rl.KeyUp) {
		if g.GameSpeed < Blazing {
			g.GameSpeed++
		}
		return
	}
	if rl.IsKeyDown(rl.KeyRightShift) && rl.IsKeyPressed(rl.KeyDown) {
		if g.GameSpeed > Snail {
			g.GameSpeed--
		}
		return
	}

	if g.GameState == Playing {
		currentTime := time.Now()
		if currentTime.Sub(g.Board.LastCatUpdate) >= g.catUpdateInterval() {
			g.Board.updateCats()
			g.Board.LastCatUpdate = currentTime
		}
	}

	rodent := g.Board.findRodent()
	if rodent == nil {
		// the rodent has been eaten by a cat
		if g.RemainingLives == 0 {
			g.GameState = GameOver
			return
		}

		g.RemainingLives--
		g.respawnRodent()
	}

	if len(g.Board.findAllCats()) == 0 && g.Board.RemainingWaves == 0 {
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

	if b.at(position) == Rodent && b.at(next) == Trap {
		b.set(position, Empty)
		b.set(next, Empty)
		g.Board.rodentDeath = append(g.Board.rodentDeath, next)
		return false
	}

	if b.at(position) == Rodent && b.at(next) == SinkHole {
		b.set(position, Empty)
		b.set(next, RodentSinkHole)
		b.InSinkHoleSince = time.Now()
		return false
	}

	if b.at(position) == RodentSinkHole {
		if time.Now().Sub(b.InSinkHoleSince) >= config.SinkHoleDuration {
			b.set(position, Rodent)
		}
		return false
	}

	if b.at(position) == Obstacle && b.at(next) == SinkHole {
		return true
	}

	if b.at(position) == Obstacle && b.at(next) == Trap {
		return false
	}

	if b.at(position) == Cat && b.at(next) == SinkHole {
		return false
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
