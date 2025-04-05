package main

import (
	"fmt"
	"math"
	"time"
)

var catsPossibleMoves = []*Move{
	{0, -1},
	{0, 1},
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

func (b *Board) updateCats() {
	currentTime := time.Now()
	if currentTime.Sub(b.LastCatUpdate) >= config.CatUpdateInterval {

		b.transformTrappedCatsToCheese()
		b.moveCats()

		b.LastCatUpdate = currentTime
	}
}

func (b *Board) moveCats() {
	cats := b.findAllCats()
	for _, cat := range cats {
		b.moveCat(cat)
	}
}

func (b *Board) transformTrappedCatsToCheese() {
	cats := b.findAllCats()
	allCatsResting := true
	for _, cat := range cats {
		if b.at(cat) != CatResting {
			allCatsResting = false
			break
		}
	}
	if allCatsResting {
		for _, cat := range cats {
			b.set(cat, Cheese)
		}
	}
}

func (b *Board) moveCat(cat *Position) {
	rodent := b.findRodent()

	bestPosition := b.aStarPathfinding(cat, rodent)
	if bestPosition != nil {
		b.set(bestPosition, Cat)
		b.set(cat, Empty)
		fmt.Println("BEST")
		return
	}

	possibleMoves := b.getPossibleMoves(cat)
	bestLegalPosition := b.minimizeDistance(cat, rodent, possibleMoves)
	if bestLegalPosition != nil {
		b.set(bestLegalPosition, Cat)
		b.set(cat, Empty)
		fmt.Println("MINIMIZE")
		return
	}

	if len(possibleMoves) > 0 {
		firstMove := possibleMoves[0]
		b.set(cat.after(firstMove), Cat)
		b.set(cat, Empty)
		fmt.Println("JUMPING AROUND")
		return
	} else {
		b.set(cat, CatResting)
		return
	}
}

func (b *Board) getPossibleMoves(cat *Position) []*Move {
	var moves []*Move

	for _, dir := range catsPossibleMoves {
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

	heuristic := func(a, b *Position) int {
		return int(math.Max(math.Abs(float64(a.Row-b.Row)), math.Abs(float64(a.Column-b.Column))))
	}

	startNode := &node{pos: *cat, g: 0, h: heuristic(cat, rodent), f: heuristic(cat, rodent)}
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
		for _, d := range catsPossibleMoves {
			neighborPos := current.pos.after(d)
			if closedSet[*neighborPos] || !b.isWalkable(neighborPos) {
				continue
			}

			gScore := current.g + 1

			neighborNode, exists := openSet[*neighborPos]
			if !exists || gScore < neighborNode.g {
				neighborNode = &node{
					pos:    *neighborPos,
					g:      gScore,
					h:      heuristic(neighborPos, rodent),
					f:      gScore + heuristic(neighborPos, rodent),
					parent: current,
				}
				openSet[*neighborPos] = neighborNode
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
			if obj == Cat || obj == CatResting {
				catPositions = append(catPositions, &Position{Row: x, Column: y})
			}
		}
	}
	return catPositions
}
