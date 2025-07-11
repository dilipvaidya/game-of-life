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

// Cell represents a cell in the Game of Life universe with its row and column indices.
// It is used as a key in the universe map to track live cells.
type Cell struct {
	R, C int
}

type GameOfLife struct {
	universe          map[Cell]struct{}
	numRows           int
	numCols           int
	neighbouringCells []Cell
	rules             []Rule
}

// CreateSeedUniverse create seed universe based on the given row, col and seed pattern
// It initializes the universe with the specified seed pattern and returns a pointer to GameOfLife.
// If the row or col is less than or equal to zero, it returns nil.
func CreateSeedUniverse(row, col int, seedPattern SEED_PATTERN, rules ...Rule) *GameOfLife {
	// should omit `input == nil` check; len() for nil slices is defined as zero
	if row <= 0 || col <= 0 {
		return nil
	}

	// return data strcture GameOfLife with universe as a new copy of the input grid.
	return &GameOfLife{
		universe: GetSeedGrid(seedPattern, row, col),
		numRows:  row,
		numCols:  col,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: rules,
	}
}

// Display displays the current state of the Game of Life universe to the standard output.
// Alive cells are represented by whiteChar, and dead cells by blackChar.
// The universe is printed row by row, with each cell separated by a space.
func (g GameOfLife) Display() {
	fmt.Println("==============")

	for rowIndex := range g.numRows {
		for colIndex := range g.numCols {
			if _, ok := g.universe[Cell{rowIndex, colIndex}]; ok {
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
func (g *GameOfLife) Run(generations int, delay time.Duration) {
	fmt.Printf("Original Generation:\n")
	g.Display()
	for i := 1; i <= generations; i++ {
		fmt.Print("\033[H\033[2J") // Clear screen before printing next frame
		fmt.Printf("Generation: %d\n", i)
		g.CreateNextGeneration()
		g.Display()
		time.Sleep(delay)
	}
}

// _wrapCellWithinUniverse return border bound cooridinates for row and col if the given cell coordinates (neighbourCellRow, neighbourCellCol) beyond the boundry
// are within the valid bounds of the Game of Life grid, otherwise false.
func (g *GameOfLife) _wrapCellWithinUniverse(cell Cell) Cell {
	// Wrap coordinates using modulo for toroidal (wrap-around) universe
	return Cell{(cell.R + g.numRows) % g.numRows, (cell.C + g.numCols) % g.numCols}
}

// _markNeighbourAlive increments the count of alive neighbours for each valid neighbouring cell
// around the given currentPosition in the neighborCounts. It iterates over all possible neighbouring
// cell offsets defined in u.neighbouringCells, calculates their absolute positions, checks if
// they are valid using _isValidNeighbour, and updates the neighborCounts accordingly.
//
// Parameters:
//   - currentPosition: a [2]int array representing the row and column of the current cell.
//   - neighborCounts: a pointer to a map that tracks the count of alive neighbours for each cell.
func (g *GameOfLife) _markNeighbourAlive(currentPosition Cell, neighborCounts *map[Cell]int) {

	for _, neighbourCell := range g.neighbouringCells {
		neighbourCell = g.
			_wrapCellWithinUniverse(
				Cell{(currentPosition.R + neighbourCell.R), (currentPosition.C + neighbourCell.C)},
			)
		(*neighborCounts)[neighbourCell]++
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
// updated:
// #3 will be: more than 4 neight would be overcrowding
// #5 cell has live neighbour on top left, it will always die
//
// The new universe is created which replaces existing one.
func (g *GameOfLife) CreateNextGeneration() {
	neighborCounts := make(map[Cell]int)
	newUniverse := make(map[Cell]struct{})

	// for every live cell, count the number of live neighbours
	for cell := range g.universe {
		g._markNeighbourAlive(cell, &neighborCounts)
	}

	// Union all the cells that have live neighbours
	// This is done to ensure that we consider all cells that could potentially become alive
	candidates := make(map[Cell]struct{})
	for cell := range neighborCounts {
		candidates[cell] = struct{}{}
	}
	// Also include all currently alive cells
	for cell := range g.universe {
		candidates[cell] = struct{}{}
	}

	// Iterate through all candidate cells to apply the rules
	// and determine their next state based on the number of live neighbours.
	// This is where the rules are applied to determine if a cell should be alive or dead
	// based on the neighborCounts.
	for cell := range candidates {
		// Check if the cell is currently alive
		_, isCellAlive := g.universe[cell]
		neighborCount := neighborCounts[cell]

		// Apply each rule to determine if the cell should be alive in the next generation
		for _, rule := range g.rules {
			if rule.Apply(cell, isCellAlive, neighborCount, g) {
				newUniverse[cell] = struct{}{}
				break // If any rule applies, we can stop checking further rules for this cell
			}
		}
	}

	g.universe = newUniverse
}
