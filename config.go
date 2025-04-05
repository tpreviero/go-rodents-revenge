package main

import "time"

type Config struct {
	SquareSize        int
	StatusBarHeight   int
	CatUpdateInterval time.Duration
	DrawGrid          bool
	CheesePoints      int
}

var config = &Config{
	SquareSize:        34,
	StatusBarHeight:   3 * 34,
	CatUpdateInterval: 1 * time.Second,
	DrawGrid:          false,
	CheesePoints:      100,
}
