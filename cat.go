package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
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
	cats := b.findAllCats()
	if len(cats) == 0 && b.RemainingWaves > 0 {
		b.RemainingWaves--
		b.respawnCats()
	}

	b.transformTrappedCatsToCheese()
	b.moveCats()
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
	anotherRodent := b.findAnotherRodent()

	bestPositionRodent := b.aStarPathfinding(cat, rodent)
	bestPositionAnotherRodent := b.aStarPathfinding(cat, anotherRodent)

	if bestPositionRodent != nil && bestPositionAnotherRodent == nil {
		// only a path to the rodent
		if b.at(bestPositionRodent) == Rodent || b.at(bestPositionRodent) == RodentSinkHole {
			b.rodentDeath = append(b.rodentDeath, bestPositionRodent)
		}
		b.set(bestPositionRodent, Cat)
		b.set(cat, Empty)
		return
	}

	if bestPositionRodent == nil && bestPositionAnotherRodent != nil {
		// only a path to another rodent
		if b.at(bestPositionAnotherRodent) == AnotherRodent || b.at(bestPositionAnotherRodent) == AnotherRodentSinkHole {
			b.rodentDeath = append(b.rodentDeath, bestPositionAnotherRodent)
		}
		b.set(bestPositionAnotherRodent, Cat)
		b.set(cat, Empty)
		return
	}

	possibleMoves := b.getPossibleMoves(cat)
	var bestLegalPosition *Position
	if b.distance(cat, rodent) <= b.distance(cat, anotherRodent) {
		bestLegalPosition = b.minimizeDistance(cat, rodent, possibleMoves)
	} else {
		bestLegalPosition = b.minimizeDistance(cat, anotherRodent, possibleMoves)
	}
	if bestLegalPosition != nil {
		b.set(bestLegalPosition, Cat)
		b.set(cat, Empty)
		return
	}

	if len(possibleMoves) > 0 {
		firstMove := possibleMoves[0]
		b.set(cat.after(firstMove), Cat)
		b.set(cat, Empty)
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

	return b.at(position) == Empty || (b.at(position) == Rodent || b.at(position) == AnotherRodent) || (b.at(position) == RodentSinkHole || b.at(position) == AnotherRodentSinkHole)
}

func (b *Board) minimizeDistance(cat, rodent *Position, possibleMoves []*Move) *Position {
	var nextPosition *Position
	minDistance := math.MaxFloat64

	for _, move := range possibleMoves {
		distance := b.distance(cat.after(move), rodent)
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

func (b *Board) respawnCats() {
	for {
		x := rl.GetRandomValue(1, 22)
		y := rl.GetRandomValue(1, 22)
		position := &Position{int(x), int(y)}
		rodent := b.findRodent()
		anotherRodent := b.findAnotherRodent()

		if b.at(position) == Empty && b.distance(position, rodent) > 5 && b.distance(position, anotherRodent) > 5 {
			b.set(position, Cat)
			break
		}
	}
}
