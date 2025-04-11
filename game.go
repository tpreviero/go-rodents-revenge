package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"time"
)

type Object int

const (
	Empty Object = iota
	Rodent
	RodentSinkHole
	Obstacle
	Wall
	Cat
	CatResting
	Cheese
	SinkHole
	Trap
)

type Position struct {
	Row    int
	Column int
}

type Move struct {
	Row    int
	Column int
}

func (p *Position) after(move *Move) *Position {
	return &Position{
		Row:    p.Row + move.Row,
		Column: p.Column + move.Column,
	}
}

type Board struct {
	Objects         [][]Object
	LastCatUpdate   time.Time
	InSinkHoleSince time.Time
	RemainingWaves  int
	rodentEaten     []*Position
}

type BoardCustomization func(position *Position) Object

func (b *Board) at(p *Position) Object {
	return b.Objects[p.Row][p.Column]
}

func (b *Board) set(p *Position, object Object) {
	b.Objects[p.Row][p.Column] = object
}

func (b *Board) distance(first, second *Position) float64 {
	return math.Abs(float64(first.Row-second.Row)) + math.Abs(float64(first.Column-second.Column))
}

type GameState int

const (
	Playing GameState = iota
	Pause
	GameOver
	Win
)

type Game struct {
	Board          *Board
	GameState      GameState
	Points         int
	RemainingLives int
	CurrentLevel   int
}

func (g *Game) PreviousLevel() {
	if g.CurrentLevel > 0 {
		g.CurrentLevel--
		customization := levelsToCustomization[g.CurrentLevel]
		g.GameState = Playing
		g.Board = NewBoard(customization)
	}
}

func (g *Game) NextLevel() {
	g.CurrentLevel++
	customization := levelsToCustomization[g.CurrentLevel]
	if customization == nil {
		g.GameState = Win
		return
	}
	g.Board = NewBoard(customization)
}

func NewGame() *Game {
	return &Game{
		Board:          NewBoard(levelsToCustomization[0]),
		GameState:      Playing,
		Points:         0,
		RemainingLives: 2,
		CurrentLevel:   0,
	}
}

var levelsToCustomization = map[int]BoardCustomization{
	0: func(position *Position) Object {
		if position.Row >= 4 && position.Row <= 18 && position.Column >= 4 && position.Column <= 18 {
			return Obstacle
		}
		return Empty
	},
	1: func(position *Position) Object {
		if rl.GetRandomValue(0, 100) < 5 {
			return Wall
		}
		if position.Row >= 4 && position.Row <= 18 && position.Column >= 4 && position.Column <= 18 {
			return Obstacle
		}
		return Empty
	},
	2: func(position *Position) Object {
		if rl.GetRandomValue(0, 100) < 5 {
			return Wall
		}
		if rl.GetRandomValue(0, 100) < 45 {
			return Obstacle
		}
		return Empty
	},
	3: func(position *Position) Object {
		if position.Row > 1 && position.Row < 21 && position.Column > 1 && position.Column < 21 {
			if rl.GetRandomValue(0, 100) < 2 {
				return SinkHole
			}
			if (position.Row%2 == 1 && position.Column%2 == 0) || (position.Row%2 == 0 && position.Column%2 == 1) {
				return Obstacle
			}
		}
		return Empty
	},
	4: func(position *Position) Object {
		if rl.GetRandomValue(0, 100) < 5 {
			return Wall
		}
		if rl.GetRandomValue(0, 100) < 5 {
			return SinkHole
		}
		if rl.GetRandomValue(0, 100) < 45 {
			return Obstacle
		}
		return Empty
	},
	5: func(position *Position) Object {
		if position.Row > 3 && position.Row < 19 && position.Column > 3 && position.Column < 19 {
			if (position.Row%2 == 1 && position.Column%2 == 0) || (position.Row%2 == 0 && position.Column%2 == 1) {
				return Wall
			}
		}
		if rl.GetRandomValue(0, 100) < 5 {
			return SinkHole
		}
		if rl.GetRandomValue(0, 100) < 25 {
			return Obstacle
		}
		return Empty
	},
	6: func(position *Position) Object {
		if rl.GetRandomValue(0, 100) < 5 {
			return SinkHole
		}
		if position.Row >= 4 && position.Row <= 18 && position.Column >= 4 && position.Column <= 18 {
			return Obstacle
		}
		if rl.GetRandomValue(0, 100) == 0 {
			return Trap
		}
		return Empty
	},
}

func NewBoard(customization BoardCustomization) *Board {
	board := &Board{
		LastCatUpdate:  time.Now(),
		RemainingWaves: 4,
		rodentEaten:    make([]*Position, 0),
	}

	board.Objects = make([][]Object, 23)

	for i := range board.Objects {
		board.Objects[i] = make([]Object, 23)
		for j := range board.Objects[i] {
			if i == 0 || i == 22 || j == 0 || j == 22 {
				board.Objects[i][j] = Wall
			} else if i == 11 && j == 11 {
				board.Objects[i][j] = Rodent
			} else {
				board.Objects[i][j] = customization(&Position{Row: i, Column: j})
			}
		}
	}

	return board
}
