package main

import rl "github.com/gen2brain/raylib-go/raylib"

type UI struct {
	rodentImage   rl.Texture2D
	catImage      rl.Texture2D
	obstacleImage rl.Texture2D
	wallImage     rl.Texture2D
}

func (ui *UI) Init() {
	rl.InitWindow(int32(config.SquareSize*23), int32(config.SquareSize*23), "Go Rodent's Revenge")
	rl.SetTargetFPS(60)

	ui.rodentImage = rl.LoadTexture("assets/rodent.png")
	ui.catImage = rl.LoadTexture("assets/cat.png")
	ui.obstacleImage = rl.LoadTexture("assets/obstacle.png")
	ui.wallImage = rl.LoadTexture("assets/wall.png")
}

func (ui *UI) Close() {
	rl.CloseWindow()
}

func (ui *UI) Draw(b *Board) {
	for i := range b.Objects {
		for j := range b.Objects[i] {
			rl.DrawRectangle(int32(j*config.SquareSize), int32(i*config.SquareSize), int32(config.SquareSize), int32(config.SquareSize), rl.NewColor(195, 195, 0, 255))
			if b.Objects[i][j] == Wall {
				rl.DrawTexture(ui.wallImage, int32(j*config.SquareSize), int32(i*config.SquareSize), rl.White)
			} else if b.Objects[i][j] == Obstacle {
				rl.DrawTexture(ui.obstacleImage, int32(j*config.SquareSize), int32(i*config.SquareSize), rl.White)
			} else if b.Objects[i][j] == Rodent {
				rl.DrawTexture(ui.rodentImage, int32(j*config.SquareSize), int32(i*config.SquareSize), rl.White)
			} else if b.Objects[i][j] == Cat {
				rl.DrawTexture(ui.catImage, int32(j*config.SquareSize), int32(i*config.SquareSize), rl.White)
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
