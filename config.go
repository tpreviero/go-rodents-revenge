package main

import "time"

type Config struct {
	SquareSize        int
	CatUpdateInterval time.Duration
	DrawGrid          bool
}

var config = &Config{
	SquareSize:        34,
	CatUpdateInterval: 1 * time.Second,
	DrawGrid:          false,
}
