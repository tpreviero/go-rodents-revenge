package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"strconv"
	"strings"
)

type Animation struct {
	texture      rl.Texture2D
	frameCounter int
}

type UI struct {
	gameTextures map[Object]rl.Texture2D
	animations   map[Position]*Animation
	rodentDeath  rl.Texture2D
	rodentLives  rl.Texture2D
	showHelp     bool
}

func NewUI() *UI {
	ui := UI{
		gameTextures: map[Object]rl.Texture2D{},
		animations:   map[Position]*Animation{},
	}
	ui.init()
	return &ui
}

func (ui *UI) init() {
	rl.InitWindow(int32(config.SquareSize*23), int32((config.SquareSize*23)+config.StatusBarHeight), "Go, Rodent's Revenge!")
	rl.SetTargetFPS(60)

	ui.LoadTextures()
}

func (ui *UI) Close() {
	rl.CloseWindow()
}

func (ui *UI) LoadTextures() {
	ui.gameTextures[Rodent] = rl.LoadTexture("assets/rodent.png")
	ui.gameTextures[RodentSinkHole] = rl.LoadTexture("assets/sinkhole-rodent.png")
	ui.gameTextures[Cat] = rl.LoadTexture("assets/cat.png")
	ui.gameTextures[CatResting] = rl.LoadTexture("assets/cat-rest.png")
	ui.gameTextures[Cheese] = rl.LoadTexture("assets/cheese.png")
	ui.gameTextures[Obstacle] = rl.LoadTexture("assets/obstacle.png")
	ui.gameTextures[Wall] = rl.LoadTexture("assets/wall.png")
	ui.gameTextures[SinkHole] = rl.LoadTexture("assets/sinkhole.png")
	ui.gameTextures[Trap] = rl.LoadTexture("assets/trap.png")

	ui.rodentLives = rl.LoadTexture("assets/rodent-lives.png")
	ui.rodentDeath = rl.LoadTexture("assets/rodent-death.png")
}

func (ui *UI) Draw(g *Game) {

	for i := range g.Board.rodentDeath {
		ui.animations[*g.Board.rodentDeath[i]] = &Animation{
			texture: ui.rodentDeath,
		}
	}
	g.Board.rodentDeath = make([]*Position, 0)

	var offset = int32(config.StatusBarHeight)
	rl.DrawRectangle(0, 0, int32(config.SquareSize*23), offset, rl.LightGray)

	for i := range g.RemainingLives {
		rl.DrawTextureEx(ui.rodentLives, rl.NewVector2(float32(config.SquareSize+(i*config.SquareSize)), float32(config.SquareSize)), 0, float32(config.SquareSize)/float32(ui.rodentLives.Width), rl.White)
	}

	level := "Level: " + strconv.Itoa(g.CurrentLevel+1)
	textWidth := rl.MeasureText(level, int32(config.SquareSize))
	rl.DrawText(level, int32(config.SquareSize*22)-textWidth, 0, int32(config.SquareSize), rl.Black)

	score := "Score: " + strconv.Itoa(g.Points)
	textWidth = rl.MeasureText(score, int32(config.SquareSize))
	rl.DrawText(score, int32(config.SquareSize*22)-textWidth, int32(config.SquareSize), int32(config.SquareSize), rl.Black)

	help := "? for help"
	textWidth = rl.MeasureText(help, int32(config.SquareSize))
	rl.DrawText(help, int32(config.SquareSize*22)-textWidth, int32(config.SquareSize*2), int32(config.SquareSize), rl.Black)

	for i := range g.Board.Objects {
		for j := range g.Board.Objects[i] {
			rl.DrawRectangle(int32(j*config.SquareSize), offset+int32(i*config.SquareSize), int32(config.SquareSize), int32(config.SquareSize), rl.NewColor(195, 195, 0, 255))

			animation := ui.animations[Position{Row: i, Column: j}]
			if animation != nil {
				if animation.Finished() {
					ui.animations[Position{Row: i, Column: j}] = nil
				} else {
					animation.Draw(Position{Row: i, Column: j})
				}
			} else {
				rl.DrawTextureEx(ui.gameTextures[g.Board.Objects[i][j]], rl.NewVector2(float32(j*config.SquareSize), float32(offset+int32(i*config.SquareSize))), 0, float32(config.SquareSize)/float32(ui.gameTextures[g.Board.Objects[i][j]].Width), rl.White)
			}
		}
	}
	if config.DrawGrid {
		drawGrid()
	}

	if g.GameState == GameOver {
		ui.displayText("Game Over", config.SquareSize)
	}

	if g.GameState == Win {
		ui.displayText("You win!", config.SquareSize)
	}

	if g.GameState == Pause {
		ui.displayText("Paused. Press P to continue.", config.SquareSize)
	}

	if ui.showHelp {
		ui.displayText("Arrow keys: Move the rodent (8 directions)\n"+
			"P: Pause the game\n"+
			"Right Shift + UP: Increase difficulty (speeds up the cats)\n"+
			"Right Shift + DOWN: Decrease difficulty (speeds up the cats)\n"+
			"Right Shift + RIGHT: Skip to the next level\n"+
			"Right Shift + LEFT: Go back to the previous level\n"+
			"?: Toggle this help screen\n"+
			"ESC: Quit the game",
			config.SquareSize/2)
	}
}

func (ui *UI) displayText(text string, fontSize int) {
	textWidth := rl.MeasureText(text, int32(fontSize))
	boxWidth := textWidth + 20
	lines := strings.Count(text, "\n")
	boxHeight := int32(fontSize+(fontSize*lines)) + 10
	if lines > 1 {
		boxHeight += 10
	}
	x := int32((config.SquareSize*23)/2) - boxWidth/2
	y := int32((23*config.SquareSize)/2) - boxHeight/2

	rl.DrawRectangle(x, y, boxWidth, boxHeight, rl.White)
	rl.DrawText(text, x+10, y+5, int32(fontSize), rl.Black)
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

func (a *Animation) Finished() bool {
	frameCount := a.texture.Width / int32(config.TextureSquareSize)
	currentFrame := (a.frameCounter / 10) % int(frameCount)
	return a.frameCounter > int((frameCount-1)*10) && currentFrame == 0
}

func (a *Animation) Draw(p Position) {
	frameCount := a.texture.Width / int32(config.TextureSquareSize)
	currentFrame := (a.frameCounter / 10) % int(frameCount)

	var offset = int32(config.StatusBarHeight)
	frameRect := rl.NewRectangle(float32(currentFrame*config.TextureSquareSize), 0, float32(config.TextureSquareSize), float32(config.TextureSquareSize))
	destRect := rl.NewRectangle(float32(p.Column*config.SquareSize), float32(offset+int32(p.Row*config.SquareSize)), float32(config.SquareSize), float32(config.SquareSize))
	rl.DrawTexturePro(a.texture, frameRect, destRect, rl.Vector2{}, 0, rl.White)

	a.frameCounter++
}
