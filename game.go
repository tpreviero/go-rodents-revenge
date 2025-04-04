package main

import "time"

type Object int

const (
	Empty Object = iota
	Rodent
	Obstacle
	Wall
	Cat
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

type Board struct {
	Objects       [][]Object
	LastCatUpdate time.Time
}

type GameState int

const (
	Playing GameState = iota
	Pause
	GameOver
	Win
)

type Game struct {
	Board     *Board
	GameState GameState
	Points    int
}

func newGame() *Game {
	return &Game{
		Board:     NewBoard(),
		GameState: Playing,
		Points:    0,
	}
}

func NewBoard() *Board {
	board := &Board{
		LastCatUpdate: time.Now(),
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
