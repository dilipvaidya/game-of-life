package gameoflife

type SEED_PATTERN int

const (
	Default SEED_PATTERN = iota
	Glider
)

func (s SEED_PATTERN) String() string {
	switch s {
	case Glider:
		return "glider"
	default:
		return "default"
	}
}

// GetSeedGrid returns a map representing the initial seed grid for Conway's Game of Life,
// based on the specified seed pattern and grid dimensions (row, col).
// The map keys are Cell arrays representing cell coordinates, and the values are booleans
// indicating whether the cell is alive (true) or dead (false).
// Supported seed patterns are defined by the SEED_PATTERN type.
// If an unsupported pattern is provided, a default seed pattern is returned.
func GetSeedGrid(seedPattern SEED_PATTERN, row, col int) map[Cell]struct{} {
	switch seedPattern {
	case Glider:
		return getSeedGliderPattern(row, col)
	default:
		return getDefaultSeedPattern(row, col)
	}
}

// getSeedGliderPattern returns a map representing the initial seed positions of a glider pattern
// for Conway's Game of Life, centered within a grid of the specified number of rows and columns.
// The map keys are Cell arrays representing the (row, column) coordinates of live cells.
// The glider pattern is a well-known configuration that moves diagonally across the grid over time.
//
// Parameters:
//
//	rows - the number of rows in the grid
//	cols - the number of columns in the grid
//
// Returns:
//
//	A map with keys as Cell coordinates and values as bool (true for live cells), representing
//	the initial positions of the glider pattern centered in the grid.
func getSeedGliderPattern(rows, cols int) map[Cell]struct{} {
	seedGrid := make(map[Cell]struct{})

	// Offset to center the glider
	r, c := rows/2, cols/2

	seedGrid[Cell{r, c + 1}] = struct{}{}
	seedGrid[Cell{r + 1, c + 2}] = struct{}{}
	seedGrid[Cell{r + 2, c}] = struct{}{}
	seedGrid[Cell{r + 2, c + 1}] = struct{}{}
	seedGrid[Cell{r + 2, c + 2}] = struct{}{}

	return seedGrid
}

// getDefaultSeedPattern returns a map representing the default seed pattern
// for Conway's Game of Life, centered within a grid of the specified number
// of rows and columns. The pattern consists of three vertically aligned live
// cells (a "blinker") centered in the grid.
//
// Parameters:
//
//	rows - the number of rows in the grid
//	cols - the number of columns in the grid
//
// Returns:
//
//	A map with keys as Cell coordinates and values as booleans indicating
//	live cells in the initial seed pattern.
func getDefaultSeedPattern(rows, cols int) map[Cell]struct{} {
	seedGrid := make(map[Cell]struct{})

	// Offset to center the glider
	r, c := rows/2, cols/2

	seedGrid[Cell{r - 1, c}] = struct{}{}
	seedGrid[Cell{r, c}] = struct{}{}
	seedGrid[Cell{r + 1, c}] = struct{}{}

	return seedGrid
}
