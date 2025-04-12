package main

import "time"

type Config struct {
	SquareSize         int
	TextureSquareSize  int
	StatusBarHeight    int
	CatUpdateIntervals map[GameSpeed]time.Duration
	DrawGrid           bool
	CheesePoints       int
	SinkHoleDuration   time.Duration
}

var config = &Config{
	SquareSize:        34,
	TextureSquareSize: 12,
	StatusBarHeight:   3 * 34,
	CatUpdateIntervals: map[GameSpeed]time.Duration{
		Snail:   2 * time.Second,
		Slow:    1 * time.Second,
		Medium:  750 * time.Millisecond,
		Fast:    500 * time.Millisecond,
		Blazing: 250 * time.Millisecond,
	},
	DrawGrid:         false,
	CheesePoints:     100,
	SinkHoleDuration: 10 * time.Second,
}
