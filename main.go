package main

import (
	"time"

	"github.com/dilipvaidya/game-of-life/gameoflife"
)

func main() {
	rows := 5
	cols := 5

	game := gameoflife.CreateSeedUniverse(rows, cols, gameoflife.Default)
	game.Run(25, 500*time.Millisecond)
}
