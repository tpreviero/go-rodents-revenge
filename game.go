package main

import (
	"math"
	"time"
)

type Object int

const (
	Empty Object = iota
	Rodent
	Obstacle
	Wall
	Cat
	CatResting
	Cheese
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
	Objects                [][]Object
	LastCatUpdate          time.Time
	RemainingNumberOfWaves int
}

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
	RamainingLives int
}

func newGame() *Game {
	return &Game{
		Board:          NewBoard(),
		GameState:      Playing,
		Points:         0,
		RamainingLives: 2,
	}
}

func NewBoard() *Board {
	board := &Board{
		LastCatUpdate:          time.Now(),
		RemainingNumberOfWaves: 3,
	}

	board.Objects = make([][]Object, 23)
	for i := range board.Objects {
		board.Objects[i] = make([]Object, 23)
		for j := range board.Objects[i] {
			if i == 0 || i == 22 || j == 0 || j == 22 {
				board.Objects[i][j] = Wall
			} else if i == 11 && j == 11 {
				board.Objects[i][j] = Rodent
			} else if i >= 4 && i <= 18 && j >= 4 && j <= 18 {
				board.Objects[i][j] = Obstacle
			} else {
				board.Objects[i][j] = Empty
			}
		}
	}

	board.Objects[1][1] = Cat

	return board
}
