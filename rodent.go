package main

func (b *Board) findRodent() *Position {
	for x, row := range b.Objects {
		for y, obj := range row {
			if obj == Rodent {
				return &Position{Row: x, Column: y}
			}
		}
	}

	return nil
}

func moveRodent(b *Board, move *Move) {
	rodent := b.findRodent()
	newRow, newColumn := rodent.Row+move.Row, rodent.Column+move.Column

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

	if b.move(newRow, newColumn, move.Row, move.Column) {
		b.Objects[rodent.Row][rodent.Column] = Empty
		b.Objects[newRow][newColumn] = Rodent
	}
}
