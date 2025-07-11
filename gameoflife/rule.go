package gameoflife

// Rule interface defines the structure for rules that can be applied to cells in the Game of Life.
// The Apply method takes a cell, its alive status, the count of its live neighbors,
// and a pointer to the GameOfLife instance. It returns a boolean indicating whether the cell
// should be alive in the next generation based on the rules defined by the implementing type.
// This allows for flexible rule definitions, enabling different behaviors in the Game of Life simulation.
type Rule interface {
	// Apply returns true if the cell should be alive in the next generation.
	Apply(cell Cell, alive bool, neighborCount int, g *GameOfLife) bool
}

type ConwayRule struct{}

func (r ConwayRule) Apply(cell Cell, alive bool, neighborCount int, g *GameOfLife) bool {
	if alive {
		// Underpopulation or Overcrowding
		return neighborCount == 2 || neighborCount == 3
	}
	// Reproduction
	return neighborCount == 3
}

type NoTopLeftNeighborRule struct{}

func (r NoTopLeftNeighborRule) Apply(cell Cell, alive bool, neighborCount int, g *GameOfLife) bool {
	topLeftCell := g._wrapCellWithinUniverse(Cell{cell.R - 1, cell.C - 1})
	_, hasTopLeft := g.universe[topLeftCell]
	return !hasTopLeft && alive
}
