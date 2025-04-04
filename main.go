package main

import (
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

func main() {
	board := NewBoard()

	ui := &UI{}
	ui.Init()
	defer ui.Close()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		board.Update()
		ui.Draw(board)

		rl.EndDrawing()
	}
}

func (b *Board) Update() {
	currentTime := time.Now()
	if currentTime.Sub(b.LastCatUpdate) >= config.CatUpdateInterval {
		cats := b.findAllCats()

		for _, cat := range cats {
			b.moveCatAt(cat)
		}

		b.LastCatUpdate = currentTime
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		moveRodent(b, -1, 0)
	} else if rl.IsKeyPressed(rl.KeyDown) {
		moveRodent(b, 1, 0)
	} else if rl.IsKeyPressed(rl.KeyLeft) {
		moveRodent(b, 0, -1)
	} else if rl.IsKeyPressed(rl.KeyRight) {
		moveRodent(b, 0, 1)
	}
}

func (b *Board) moveCatAt(cat *Position) {
	rodent := b.findRodent()

	bestMove := b.aStarPathfinding(cat, rodent)
	if bestMove != nil {
		b.Objects[bestMove.Row][bestMove.Column] = Cat
		b.Objects[cat.Row][cat.Column] = Empty
		fmt.Println("BEST")
		return
	}

	possibleMoves := b.getPossibleMoves(cat)
	bestLegalMove := b.minimizeDistance(cat, rodent, possibleMoves)
	if bestLegalMove != nil {
		b.Objects[bestLegalMove.Row][bestLegalMove.Column] = Cat
		b.Objects[cat.Row][cat.Column] = Empty
		fmt.Println("MINIMIZE")
		return
	}

	if len(possibleMoves) > 0 {
		firstMove := possibleMoves[0]
		b.Objects[firstMove.Row][firstMove.Column] = Cat
		b.Objects[cat.Row][cat.Column] = Empty
		fmt.Println("JUMPING AROUND")
		return
	}
}

func (b *Board) getPossibleMoves(cat *Position) []Move {
	var moves []Move
	directions := []Move{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	for _, dir := range directions {
		nextPosition := Position{Row: cat.Row + dir.Row, Column: cat.Column + dir.Column}
		if b.isWalkable(nextPosition) {
			moves = append(moves, dir)
		}
	}
	return moves
}

type Node struct {
	pos       Position
	cost      int
	heuristic int
	parent    *Node
	index     int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (b *Board) aStarPathfinding(cat, rodent *Position) *Position {
	type node struct {
		pos     Position
		g, h, f int
		parent  *node
	}

	allPossibleMoves := []Move{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	heuristic := func(a, b Position) int {
		return abs(a.Row-b.Row) + abs(a.Column-b.Column) // Manhattan Distance
	}

	startNode := &node{pos: *cat, g: 0, h: heuristic(*cat, *rodent), f: heuristic(*cat, *rodent)}
	openSet := map[Position]*node{startNode.pos: startNode}
	closedSet := make(map[Position]bool)

	for len(openSet) > 0 {
		// Get node with lowest f-score
		var current *node
		for _, n := range openSet {
			if current == nil || n.f < current.f {
				current = n
			}
		}
		delete(openSet, current.pos)
		closedSet[current.pos] = true

		if current.pos == *rodent {
			// Reconstruct path to get the next position after start
			for current.parent != nil && current.parent.pos != *cat {
				current = current.parent
			}
			return &current.pos
		}

		// Expand neighbors
		for _, d := range allPossibleMoves {
			neighborPos := Position{current.pos.Row + d.Row, current.pos.Column + d.Column}
			if closedSet[neighborPos] || !b.isWalkable(neighborPos) {
				continue
			}

			gScore := current.g + 1

			neighborNode, exists := openSet[neighborPos]
			if !exists || gScore < neighborNode.g {
				neighborNode = &node{
					pos:    neighborPos,
					g:      gScore,
					h:      heuristic(neighborPos, *rodent),
					f:      gScore + heuristic(neighborPos, *rodent),
					parent: current,
				}
				openSet[neighborPos] = neighborNode
			}
		}
	}

	return nil // No path found
}

func (b *Board) isWalkable(position Position) bool {
	if position.Row < 0 || position.Row >= len(b.Objects) || position.Column < 0 || position.Column >= len(b.Objects[0]) {
		return false
	}

	return b.Objects[position.Row][position.Column] == Empty || b.Objects[position.Row][position.Column] == Rodent
}

func (b *Board) minimizeDistance(cat, rodent *Position, possibleMoves []Move) *Position {
	var nextPosition *Position
	minDistance := math.MaxFloat64

	for _, move := range possibleMoves {
		distance := math.Abs(float64(cat.Row+move.Row-rodent.Row)) + math.Abs(float64(cat.Column+move.Column-rodent.Column))
		if distance < minDistance {
			minDistance = distance
			candidate := Position{cat.Row + move.Row, cat.Column + move.Column}
			if b.isWalkable(candidate) {
				nextPosition = &candidate
			}
		}
	}
	return nextPosition
}

func (b *Board) findAllCats() []*Position {
	var catPositions []*Position
	for x, row := range b.Objects {
		for y, obj := range row {
			if obj == Cat {
				catPositions = append(catPositions, &Position{Row: x, Column: y})
			}
		}
	}
	return catPositions
}

func (b *Board) findRodent() *Position {
	for x, row := range b.Objects {
		for y, obj := range row {
			if obj == Rodent {
				return &Position{Row: x, Column: y}
			}
		}
	}

	panic("no rodent found")
}

func moveRodent(b *Board, dRow, dColumn int) {
	rodent := b.findRodent()
	newRow, newColumn := rodent.Row+dRow, rodent.Column+dColumn

	if newRow < 0 || newRow >= len(b.Objects[0]) || newColumn < 0 || newColumn >= len(b.Objects) {
		return
	}

	if b.Objects[newRow][newColumn] == Wall {
		return
	}

	if b.Objects[newRow][newColumn] == Empty {
		b.Objects[rodent.Row][rodent.Column] = Empty
		b.Objects[newRow][newColumn] = Rodent
		return
	}

	//TODO: if moving into a cat should makes the game over

	if b.move(newRow, newColumn, dRow, dColumn) {
		b.Objects[rodent.Row][rodent.Column] = Empty
		b.Objects[newRow][newColumn] = Rodent
	}
}

func (b *Board) move(row, colum, dRow, dColumn int) bool {
	newRow, newColumn := row+dRow, colum+dColumn

	if b.Objects[newRow][newColumn] == Wall {
		return false
	}

	if b.Objects[newRow][newColumn] == Empty {
		b.Objects[newRow][newColumn] = b.Objects[row][colum]
		b.Objects[row][colum] = Empty
		return true
	}

	if b.move(newRow, newColumn, dRow, dColumn) {
		b.Objects[newRow][newColumn] = b.Objects[row][colum]
		b.Objects[row][colum] = Empty
		return true
	}

	return false
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
