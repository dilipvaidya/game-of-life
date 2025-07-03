package main

import (
	"fmt"
	"time"
)

const (
	VALUE_DEAD_CELL = 0
	VALUE_LIVE_CELL = 1

	// ANSI escape codes for black and white background
	blackChar = "\033[40m \033[0m"
	whiteChar = "\033[47m \033[0m"
)

// Game of Life implementation in Go
type GameOfLife struct {
	// todo: optimize the data structure to use a single slice instead of a 2D slice.
	// this will help to reduce the memory usage and improve performance.
	universe          [][]int
	numRows           int
	numCols           int
	neighbouringCells [][]int
}

func _deepCopyUniverse(originalUniverse [][]int) [][]int {
	// create a deep of the input grid
	newUniverse := make([][]int, len(originalUniverse))

	// iterate over the orginal grid and copy to new universer
	for currentRowIndex, currentRow := range originalUniverse {
		newUniverse[currentRowIndex] = make([]int, len(currentRow), cap(currentRow))
		copy(newUniverse[currentRowIndex], currentRow)
	}

	return newUniverse
}

// create GameOfLife Universe from the input
func CreateUniverse(inputGrid [][]int) *GameOfLife {
	// should omit `input == nil` check; len() for nil slices is defined as zero
	if len(inputGrid) <= 0 || len(inputGrid[0]) <= 0 {
		return nil
	}

	// create a deep of the input grid
	numberOfRows := len(inputGrid)
	numberOfCols := len(inputGrid[0])
	universe := _deepCopyUniverse(inputGrid)

	// return data strcture GameOfLife with universe as a new copy of the input grid.
	return &GameOfLife{
		universe: universe,
		numRows:  numberOfRows,
		numCols:  numberOfCols,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
}

func (u GameOfLife) Print() {
	fmt.Println("==============")

	for _, row := range u.universe {
		for _, value := range row {
			if value == VALUE_LIVE_CELL {
				fmt.Print(" ", whiteChar)
			} else {
				fmt.Print(" ", blackChar)
			}
		}
		fmt.Print("\n\n")
	}
}

// checkCellValidity
func (u *GameOfLife) _isValidNeighbour(neighbourCellRow int, neighbourCellCol int) bool {
	return neighbourCellRow >= 0 && neighbourCellRow < u.numRows &&
		neighbourCellCol >= 0 && neighbourCellCol < u.numCols
}

// get the status of the neghbours
func (u *GameOfLife) _getNeighbourStats(currentRow, currentCol int) (int, int) {
	numberOfLiveNeighbours := 0
	numberOfDeadNeughbours := 0

	for _, neighbourCell := range u.neighbouringCells {
		neighbourCellRow := currentRow + neighbourCell[0]
		neighbourCellCol := currentCol + neighbourCell[1]

		if u._isValidNeighbour(neighbourCellRow, neighbourCellCol) {
			if u.universe[neighbourCellRow][neighbourCellCol] == VALUE_LIVE_CELL {
				numberOfLiveNeighbours++
			} else {
				numberOfDeadNeughbours++
			}
		}
	}

	return numberOfLiveNeighbours, numberOfDeadNeughbours
}

func (u *GameOfLife) PlayNextGeneration() {
	// there are threee rules for live cell while one for dead.
	// while iterating over the grid, based on cell is alive or dead, apply rule accordingly
	// may be we need new copy of the grid, as we have to keep the previous one until we finish iterating over the previous one.

	newUniverse := make([][]int, u.numRows)
	for rowIndex, row := range u.universe {

		newUniverse[rowIndex] = make([]int, len(row), cap(row))

		for colIndex, value := range row {

			// find neighbour cell stats
			numberOfLiveNeighbours, _ := u._getNeighbourStats(rowIndex, colIndex)

			// find which rule to apply based on current cell status: live / dead
			if value == VALUE_DEAD_CELL {
				// dead cell, only one rule can be applied
				if numberOfLiveNeighbours == 3 {
					// 4. Any dead cell with exactly three live neighbours becomes a live cell,​ as if by reproduction.
					newUniverse[rowIndex][colIndex] = VALUE_LIVE_CELL
				}
			} else {
				// live cell, there are thee possibilities of rules
				// 1. Any live cell with fewer than two live neighbours dies​, as if caused by under­population.
				// 2. Any live cell with two or three live neighbours lives​ on to the next generation.
				// 3. Any live cell with more than three live neighbours dies​, as if by overcrowding.
				switch {
				case numberOfLiveNeighbours < 2:
					newUniverse[rowIndex][colIndex] = VALUE_DEAD_CELL
				case numberOfLiveNeighbours == 2 || numberOfLiveNeighbours == 3:
					newUniverse[rowIndex][colIndex] = VALUE_LIVE_CELL
				case numberOfLiveNeighbours > 3:
					newUniverse[rowIndex][colIndex] = VALUE_DEAD_CELL
				default:
					fmt.Println("unknown status")
				}
			}
		}
	}

	//u.universe = _deepCopyUniverse(newUniverse) // deepcopy is uneccessary as newUniverse is just constructed.
	u.universe = newUniverse
}

func (u *GameOfLife) Run(generations int, delay time.Duration) {
	for i := 0; i < generations; i++ {
		fmt.Printf("Generation: %d\n", i)
		u.Print()
		u.PlayNextGeneration()
		time.Sleep(delay)
	}
}
