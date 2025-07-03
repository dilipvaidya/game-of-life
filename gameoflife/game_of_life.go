package gameoflife

import (
	"fmt"
	"time"
)

// Constants for cell values and display characters
const (
	VALUE_DEAD_CELL = 0
	VALUE_LIVE_CELL = 1
	blackChar       = "\033[40m \033[0m"
	whiteChar       = "\033[47m \033[0m"
)

type GameOfLife struct {
	universe          map[[2]int]struct{} // todo: [2]int can also be defined as `type Cell struct { Row int, Col int }`
	numRows           int
	numCols           int
	neighbouringCells [][]int
}

// CreateSeedUniverse create seed universe based on the given row, col and seed pattern
// It initializes the universe with the specified seed pattern and returns a pointer to GameOfLife.
// If the row or col is less than or equal to zero, it returns nil.
func CreateSeedUniverse(row, col int, seedPattern SEED_PATTERN) *GameOfLife {
	// should omit `input == nil` check; len() for nil slices is defined as zero
	if row <= 0 || col <= 0 {
		return nil
	}

	// return data strcture GameOfLife with universe as a new copy of the input grid.
	return &GameOfLife{
		universe: GetSeedGrid(seedPattern, row, col),
		numRows:  row,
		numCols:  col,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
}

// Display displays the current state of the Game of Life universe to the standard output.
// Alive cells are represented by whiteChar, and dead cells by blackChar.
// The universe is printed row by row, with each cell separated by a space.
func (u GameOfLife) Display() {
	fmt.Println("==============")

	for rowIndex := range u.numRows {
		for colIndex := range u.numCols {
			if _, ok := u.universe[[2]int{rowIndex, colIndex}]; ok {
				fmt.Print(" ", whiteChar)
			} else {
				fmt.Print(" ", blackChar)
			}
		}
		fmt.Print("\n\n")
	}
}

// Run simulates the Game of Life for a specified number of generations.
// It prints the initial state, then iteratively generates and prints each subsequent generation,
// pausing for the specified delay between generations.
//
// Parameters:
//
//	generations - the number of generations to simulate.
//	delay - the duration to wait between each generation.
func (u *GameOfLife) Run(generations int, delay time.Duration) {
	fmt.Printf("Original Generation:\n")
	u.Display()
	for i := 1; i <= generations; i++ {
		fmt.Printf("Generation: %d\n", i)
		u.CreateNextGeneration()
		u.Display()
		time.Sleep(delay)
	}
}

// _isValidNeighbour return true if the given cell coordinates (neighbourCellRow, neighbourCellCol)
// are within the valid bounds of the Game of Life grid, otherwise false.
func (u *GameOfLife) _isValidNeighbour(neighbourCellRow int, neighbourCellCol int) bool {
	return neighbourCellRow >= 0 && neighbourCellRow < u.numRows &&
		neighbourCellCol >= 0 && neighbourCellCol < u.numCols
}

// _markNeighbourAlive increments the count of alive neighbours for each valid neighbouring cell
// around the given currentPosition in the neighborCounts. It iterates over all possible neighbouring
// cell offsets defined in u.neighbouringCells, calculates their absolute positions, checks if
// they are valid using _isValidNeighbour, and updates the neighborCounts accordingly.
//
// Parameters:
//   - currentPosition: a [2]int array representing the row and column of the current cell.
//   - neighborCounts: a pointer to a map that tracks the count of alive neighbours for each cell.
func (u *GameOfLife) _markNeighbourAlive(currentPosition [2]int, neighborCounts *map[[2]int]int) {

	for _, neighbourCell := range u.neighbouringCells {
		neighbourCellRow := currentPosition[0] + neighbourCell[0]
		neighbourCellCol := currentPosition[1] + neighbourCell[1]

		if u._isValidNeighbour(neighbourCellRow, neighbourCellCol) {
			(*neighborCounts)[[2]int{neighbourCellRow, neighbourCellCol}]++
		}
	}
}

// CreateNextGeneration advances the Game of Life universe by one generation
// according to Conway's Game of Life rules. It calculates the next state of
// each cell based on the number of live neighbours:
//  1. Any live cell with fewer than two live neighbours dies (underpopulation).
//  2. Any live cell with two or three live neighbours lives on to the next generation.
//  3. Any live cell with more than three live neighbours dies (overcrowding).
//  4. Any dead cell with exactly three live neighbours becomes a live cell (reproduction).
//
// The new universe is created which replaces existing one.
func (u *GameOfLife) CreateNextGeneration() {
	neighborCounts := make(map[[2]int]int)
	newUniverse := make(map[[2]int]struct{})

	// for every live cell, count the number of live neighbours
	for key := range u.universe {
		u._markNeighbourAlive(key, &neighborCounts)
	}

	for key, value := range neighborCounts {
		_, isCellAlive := u.universe[[2]int{key[0], key[1]}]
		if value == 3 && !isCellAlive {
			newUniverse[[2]int{key[0], key[1]}] = struct{}{}
		}
	}

	for key := range u.universe {
		neighborCount := neighborCounts[[2]int{key[0], key[1]}]
		if neighborCount == 2 || neighborCount == 3 {
			newUniverse[[2]int{key[0], key[1]}] = struct{}{}
		}
	}

	u.universe = newUniverse
}
