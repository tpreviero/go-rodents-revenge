package main

import (
	"fmt"
	"math"
	"time"
)

var catsPossivleMoves = []*Move{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

func (b *Board) updateCats() {
	currentTime := time.Now()
	if currentTime.Sub(b.LastCatUpdate) >= config.CatUpdateInterval {

		cats := b.findAllCats()

		for _, cat := range cats {
			b.moveCatAt(cat)
		}

		b.LastCatUpdate = currentTime
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

func (b *Board) getPossibleMoves(cat *Position) []*Move {
	var moves []*Move

	for _, dir := range catsPossivleMoves {
		nextPosition := cat.after(dir)
		if b.isWalkable(nextPosition) {
			moves = append(moves, dir)
		}
	}
	return moves
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

func (b *Board) isWalkable(position *Position) bool {
	if position.Row < 0 || position.Row >= len(b.Objects) || position.Column < 0 || position.Column >= len(b.Objects[0]) {
		return false
	}

	return b.Objects[position.Row][position.Column] == Empty || b.Objects[position.Row][position.Column] == Rodent
}

func (b *Board) minimizeDistance(cat, rodent *Position, possibleMoves []*Move) *Position {
	var nextPosition *Position
	minDistance := math.MaxFloat64

	for _, move := range possibleMoves {
		distance := math.Abs(float64(cat.Row+move.Row-rodent.Row)) + math.Abs(float64(cat.Column+move.Column-rodent.Column))
		if distance < minDistance {
			minDistance = distance
			candidate := cat.after(move)
			if b.isWalkable(candidate) {
				nextPosition = candidate
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
