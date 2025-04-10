package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UI struct {
	gameTextures map[Object]rl.Texture2D
	rodentLives  rl.Texture2D
}

func NewUI() *UI {
	ui := UI{
		gameTextures: map[Object]rl.Texture2D{},
	}
	ui.init()
	return &ui
}

func (ui *UI) init() {
	rl.InitWindow(int32(config.SquareSize*23), int32((config.SquareSize*23)+config.StatusBarHeight), "Go Rodent's Revenge")
	rl.SetTargetFPS(60)

	ui.LoadTextures()
}

func (ui *UI) Close() {
	rl.CloseWindow()
}

func (ui *UI) LoadTextures() {
	ui.gameTextures[Rodent] = rl.LoadTexture("assets/rodent.png")
	ui.gameTextures[Cat] = rl.LoadTexture("assets/cat.png")
	ui.gameTextures[CatResting] = rl.LoadTexture("assets/cat-rest.png")
	ui.gameTextures[Cheese] = rl.LoadTexture("assets/cheese.png")
	ui.gameTextures[Obstacle] = rl.LoadTexture("assets/obstacle.png")
	ui.gameTextures[Wall] = rl.LoadTexture("assets/wall.png")
	ui.gameTextures[SinkHole] = rl.LoadTexture("assets/sinkhole.png")
	ui.gameTextures[SinkHoleRodent] = rl.LoadTexture("assets/sinkhole-rodent.png")
	ui.gameTextures[Trap] = rl.LoadTexture("assets/trap.png")

	ui.rodentLives = rl.LoadTexture("assets/rodent-lives.png")
}

func (ui *UI) Draw(g *Game) {
	var offset = int32(config.StatusBarHeight)
	rl.DrawRectangle(0, 0, int32(config.SquareSize*23), offset, rl.LightGray)

	for i := range g.RamainingLives {
		//rl.DrawTexture(ui.rodentLives, int32(config.SquareSize+(i*config.SquareSize)), int32(config.SquareSize), rl.White)
		rl.DrawTextureEx(ui.rodentLives, rl.NewVector2(float32(config.SquareSize+(i*config.SquareSize)), float32(config.SquareSize)), 0, float32(config.SquareSize)/float32(ui.rodentLives.Width), rl.White)
	}

	text := strconv.Itoa(g.Points)
	textWidth := rl.MeasureText(text, int32(config.SquareSize))
	rl.DrawText(text, int32(config.SquareSize*22)-textWidth, int32(config.SquareSize), int32(config.SquareSize), rl.Black)

	for i := range g.Board.Objects {
		for j := range g.Board.Objects[i] {
			rl.DrawRectangle(int32(j*config.SquareSize), offset+int32(i*config.SquareSize), int32(config.SquareSize), int32(config.SquareSize), rl.NewColor(195, 195, 0, 255))
			rl.DrawTextureEx(ui.gameTextures[g.Board.Objects[i][j]], rl.NewVector2(float32(j*config.SquareSize), float32(offset+int32(i*config.SquareSize))), 0, float32(config.SquareSize)/float32(ui.gameTextures[g.Board.Objects[i][j]].Width), rl.White)
		}
	}
	if config.DrawGrid {
		drawGrid()
	}

	if g.GameState == GameOver {
		text := "Game Over"
		textWidth := rl.MeasureText(text, int32(config.SquareSize))
		boxWidth := textWidth + 20
		boxHeight := int32(config.SquareSize) + 10
		x := int32((config.SquareSize*23)/2) - boxWidth/2
		y := int32((23*config.SquareSize)/2) - boxHeight/2

		rl.DrawRectangle(x, y, boxWidth, boxHeight, rl.White)
		rl.DrawText(text, x+10, y+5, int32(config.SquareSize), rl.Black)
	}

	if g.GameState == Win {
		text := "You win!"
		textWidth := rl.MeasureText(text, int32(config.SquareSize))
		boxWidth := textWidth + 20
		boxHeight := int32(config.SquareSize) + 10
		x := int32((config.SquareSize*23)/2) - boxWidth/2
		y := int32((23*config.SquareSize)/2) - boxHeight/2

		rl.DrawRectangle(x, y, boxWidth, boxHeight, rl.White)
		rl.DrawText(text, x+10, y+5, int32(config.SquareSize), rl.Black)
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
