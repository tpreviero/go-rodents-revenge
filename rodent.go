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
