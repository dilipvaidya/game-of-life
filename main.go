package main

import (
	"time"
)

func main() {
	universeGrid := [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0},
	}
	gameOfLife := CreateUniverse(universeGrid)
	gameOfLife.Run(10, 500*time.Millisecond)
}
