package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UI struct {
	rodent      rl.Texture2D
	rodentLives rl.Texture2D
	cat         rl.Texture2D
	catResting  rl.Texture2D
	cheese      rl.Texture2D
	obstacle    rl.Texture2D
	wall        rl.Texture2D
}

func (ui *UI) Init() {
	rl.InitWindow(int32(config.SquareSize*23), int32((config.SquareSize*23)+config.StatusBarHeight), "Go Rodent's Revenge")
	rl.SetTargetFPS(60)

	ui.LoadTextures()
}

func (ui *UI) Close() {
	rl.CloseWindow()
}

func (ui *UI) LoadTextures() {
	ui.rodent = rl.LoadTexture("assets/rodent.png")
	ui.rodentLives = rl.LoadTexture("assets/rodent-lives.png")
	ui.cat = rl.LoadTexture("assets/cat.png")
	ui.catResting = rl.LoadTexture("assets/cat-rest.png")
	ui.cheese = rl.LoadTexture("assets/cheese.png")
	ui.obstacle = rl.LoadTexture("assets/obstacle.png")
	ui.wall = rl.LoadTexture("assets/wall.png")
}

func (ui *UI) Draw(g *Game) {
	var offset = int32(config.StatusBarHeight)
	rl.DrawRectangle(0, 0, int32(config.SquareSize*23), offset, rl.LightGray)

	for i := range g.RamainingLives {
		rl.DrawTexture(ui.rodentLives, int32(config.SquareSize+(i*config.SquareSize)), int32(config.SquareSize), rl.White)
	}

	text := strconv.Itoa(g.Points)
	textWidth := rl.MeasureText(text, int32(config.SquareSize))
	rl.DrawText(text, int32(config.SquareSize*22)-textWidth, int32(config.SquareSize), int32(config.SquareSize), rl.Black)

	for i := range g.Board.Objects {
		for j := range g.Board.Objects[i] {
			rl.DrawRectangle(int32(j*config.SquareSize), offset+int32(i*config.SquareSize), int32(config.SquareSize), int32(config.SquareSize), rl.NewColor(195, 195, 0, 255))
			if g.Board.Objects[i][j] == Wall {
				rl.DrawTexture(ui.wall, int32(j*config.SquareSize), offset+int32(i*config.SquareSize), rl.White)
			} else if g.Board.Objects[i][j] == Obstacle {
				rl.DrawTexture(ui.obstacle, int32(j*config.SquareSize), offset+int32(i*config.SquareSize), rl.White)
			} else if g.Board.Objects[i][j] == Rodent {
				rl.DrawTexture(ui.rodent, int32(j*config.SquareSize), offset+int32(i*config.SquareSize), rl.White)
			} else if g.Board.Objects[i][j] == Cat {
				rl.DrawTexture(ui.cat, int32(j*config.SquareSize), offset+int32(i*config.SquareSize), rl.White)
			} else if g.Board.Objects[i][j] == CatResting {
				rl.DrawTexture(ui.catResting, int32(j*config.SquareSize), offset+int32(i*config.SquareSize), rl.White)
			} else if g.Board.Objects[i][j] == Cheese {
				rl.DrawTexture(ui.cheese, int32(j*config.SquareSize), offset+int32(i*config.SquareSize), rl.White)
			}
		}
	}
	if config.DrawGrid {
		drawGrid()
	}
}

func drawGrid() {
	for i := range 23 * config.SquareSize {
		rl.DrawLineV(
			rl.NewVector2(float32(config.SquareSize*i), 0),
			rl.NewVector2(float32(config.SquareSize*i), float32(23*config.SquareSize)),
			rl.LightGray,
		)
	}

	for i := range 23 * config.SquareSize {
		rl.DrawLineV(
			rl.NewVector2(0, float32(config.SquareSize*i)),
			rl.NewVector2(float32(23*config.SquareSize), float32(config.SquareSize*i)),
			rl.LightGray,
		)
	}
}
