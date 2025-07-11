package gameoflife

import (
	"testing"
)

func Test_wrapCellWithinUniverse(t *testing.T) {
	u := &GameOfLife{
		numRows: 5,
		numCols: 5,
	}

	tests := []struct {
		name string
		cell Cell
		want Cell
	}{
		{"top-left-corner", Cell{R: 0, C: 0}, Cell{R: 0, C: 0}},
		{"bottom-right-corner", Cell{R: 4, C: 4}, Cell{R: 4, C: 4}},
		{"center", Cell{R: 2, C: 2}, Cell{R: 2, C: 2}},
		{"negative-row", Cell{R: -1, C: 0}, Cell{R: 4, C: 0}},
		{"negative-col", Cell{R: 0, C: -1}, Cell{R: 0, C: 4}},
		{"row-out-of-bounds", Cell{R: 5, C: 0}, Cell{R: 0, C: 0}},
		{"col-out-of-bounds", Cell{R: 0, C: 5}, Cell{R: 0, C: 0}},
		{"both-out-of-bounds", Cell{R: 5, C: 5}, Cell{R: 0, C: 0}},
		{"both-negative", Cell{R: -1, C: -1}, Cell{R: 4, C: 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u._wrapCellWithinUniverse(tt.cell)
			if got != tt.want {
				t.Errorf("_wrapCellWithinUniverse(%v) = %v; want %v", tt.cell, got, tt.want)
			}
		})
	}
}

func TestCreateSeedUniverse_DefaultAndGliderPatterns(t *testing.T) {
	type args struct {
		row         int
		col         int
		seedPattern SEED_PATTERN
	}
	// Use the real GetSeedGrid implementation for these tests
	tests := []struct {
		name         string
		args         args
		wantUniverse map[Cell]struct{}
	}{
		{
			name: "Default pattern 5x5",
			args: args{5, 5, Default},
			wantUniverse: map[Cell]struct{}{
				{1, 2}: {}, {2, 2}: {}, {3, 2}: {},
			},
		},
		{
			name: "Glider pattern 5x5",
			args: args{5, 5, Glider},
			wantUniverse: map[Cell]struct{}{
				{2, 3}: {}, {3, 4}: {}, {4, 2}: {}, {4, 3}: {}, {4, 4}: {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateSeedUniverse(tt.args.row, tt.args.col, tt.args.seedPattern)
			if got == nil {
				t.Fatalf("CreateSeedUniverse returned nil for valid input")
			}
			for cell := range tt.wantUniverse {
				if _, isAlive := got.universe[cell]; !isAlive {
					t.Errorf("expected cell %v to be alive", cell)
				}
			}
			if len(got.neighbouringCells) != 8 {
				t.Errorf("neighbouringCells = %v; want 8 directions", got.neighbouringCells)
			}
		})
	}
}

func Test_markNeighbourAlive(t *testing.T) {
	u := &GameOfLife{
		numRows: 3,
		numCols: 3,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
	}

	tests := []struct {
		name           string
		currentPos     Cell
		wantAliveCells map[Cell]int
	}{
		{
			name:       "center cell",
			currentPos: Cell{R: 1, C: 1},
			wantAliveCells: map[Cell]int{
				{0, 0}: 1, {0, 1}: 1, {0, 2}: 1,
				{1, 0}: 1, {1, 2}: 1,
				{2, 0}: 1, {2, 1}: 1, {2, 2}: 1,
			},
		},
		{
			name:       "top-left corner",
			currentPos: Cell{R: 0, C: 0},
			wantAliveCells: map[Cell]int{
				{2, 2}: 1, {2, 0}: 1, {2, 1}: 1,
				{0, 2}: 1, {0, 1}: 1,
				{1, 0}: 1, {1, 1}: 1, {1, 2}: 1,
			},
		},
		{
			name:       "bottom-right corner",
			currentPos: Cell{R: 2, C: 2},
			wantAliveCells: map[Cell]int{
				{1, 1}: 1, {1, 2}: 1, {1, 0}: 1,
				{2, 1}: 1, {2, 0}: 1,
				{0, 1}: 1, {0, 2}: 1, {0, 0}: 1,
			},
		},
		{
			name:       "edge cell",
			currentPos: Cell{R: 0, C: 1},
			wantAliveCells: map[Cell]int{
				{2, 0}: 1, {2, 1}: 1, {2, 2}: 1,
				{0, 0}: 1, {0, 2}: 1,
				{1, 0}: 1, {1, 1}: 1, {1, 2}: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliveMap := make(map[Cell]int)
			u._markNeighbourAlive(tt.currentPos, &aliveMap)
			if len(aliveMap) != len(tt.wantAliveCells) {
				t.Errorf("got %d alive cells, want %d", len(aliveMap), len(tt.wantAliveCells))
			}
			for cell, wantCount := range tt.wantAliveCells {
				if gotCount := aliveMap[cell]; gotCount != wantCount {
					t.Errorf("cell %v: got count %d, want %d", cell, gotCount, wantCount)
				}
			}
		})
	}
}

/*
func TestCreateNextGeneration_Blinker(t *testing.T) {
	// Blinker pattern (period 2 oscillator)
	// Generation 0 (vertical):
	// . . .
	// X X X
	// . . .
	initialUniverse := map[Cell]struct{}{
		{1, 0}: {},
		{1, 1}: {},
		{1, 2}: {},
	}
	u := &GameOfLife{
		universe: initialUniverse,
		numRows:  3,
		numCols:  3,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	u.CreateNextGeneration()
	wantUniverse := map[Cell]struct{}{
		{0, 1}: {},
		{1, 1}: {},
		{2, 1}: {},
	}
	if len(u.universe) != len(wantUniverse) {
		t.Errorf("got %d alive cells, want %d", len(u.universe), len(wantUniverse))
	}
	for cell := range wantUniverse {
		if _, isAlive := u.universe[cell]; !isAlive {
			t.Errorf("expected cell %v to be alive in next generation", cell)
		}
	}
}
*/

func TestCreateNextGeneration_BlockStillLife(t *testing.T) {
	// Block pattern (still life)
	// . X X
	// . X X
	initialUniverse := map[Cell]struct{}{
		{0, 1}: {}, {0, 2}: {},
		{1, 1}: {}, {1, 2}: {},
	}
	u := &GameOfLife{
		universe: initialUniverse,
		numRows:  4,
		numCols:  4,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	u.CreateNextGeneration()
	wantUniverse := map[Cell]struct{}{
		{0, 1}: {}, {0, 2}: {},
		{1, 1}: {}, {1, 2}: {},
	}
	if len(u.universe) != len(wantUniverse) {
		t.Errorf("got %d alive cells, want %d", len(u.universe), len(wantUniverse))
	}
	for cell := range wantUniverse {
		if _, isAlive := u.universe[cell]; !isAlive {
			t.Errorf("expected cell %v to be alive in next generation", cell)
		}
	}
}

func TestCreateNextGeneration_Underpopulation(t *testing.T) {
	// Single live cell should die
	initialUniverse := map[Cell]struct{}{
		{1, 1}: {},
	}
	u := &GameOfLife{
		universe: initialUniverse,
		numRows:  3,
		numCols:  3,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	u.CreateNextGeneration()
	if len(u.universe) != 0 {
		t.Errorf("expected all cells to be dead due to underpopulation, got %v", u.universe)
	}
}

func TestCreateNextGeneration_Reproduction(t *testing.T) {
	// Dead cell with exactly three live neighbours becomes alive
	// . X .
	// X . X
	// . . .
	initialUniverse := map[Cell]struct{}{
		{0, 1}: {},
		{1, 0}: {},
		{1, 2}: {},
	}
	u := &GameOfLife{
		universe: initialUniverse,
		numRows:  3,
		numCols:  3,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	u.CreateNextGeneration()
	wantUniverse := map[Cell]struct{}{
		{0, 1}: {},
		{1, 1}: {}, // new cell born
	}
	for cell := range wantUniverse {
		if _, isAlive := u.universe[cell]; !isAlive {
			t.Errorf("expected cell %v to be alive in next generation", cell)
		}
	}
}

func BenchmarkCreateNextGeneration_100x100_Glider(b *testing.B) {
	// Create a dense 100x100 universe with a checkerboard pattern
	numRows, numCols := 100, 100
	// Using the Glider seed pattern for the benchmark
	u := &GameOfLife{
		universe: GetSeedGrid(Glider, numRows, numCols),
		numRows:  numRows,
		numCols:  numCols,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}

func BenchmarkCreateNextGeneration_1000x1000_Glider(b *testing.B) {
	// Create a dense 1000x1000 universe with a checkerboard pattern
	numRows, numCols := 1000, 1000
	// Using the Glider seed pattern for the benchmark
	u := &GameOfLife{
		universe: GetSeedGrid(Glider, numRows, numCols),
		numRows:  numRows,
		numCols:  numCols,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}

func BenchmarkCreateNextGeneration_100x100(b *testing.B) {
	// Create a dense 100x100 universe with a checkerboard pattern
	numRows, numCols := 100, 100
	universe := make(map[Cell]struct{})
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			if (r+c)%2 == 0 {
				universe[Cell{r, c}] = struct{}{}
			}
		}
	}
	u := &GameOfLife{
		universe: universe,
		numRows:  numRows,
		numCols:  numCols,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}

func BenchmarkCreateNextGeneration_1000x1000(b *testing.B) {
	// Create a dense 1000x1000 universe with a checkerboard pattern
	numRows, numCols := 1000, 1000
	universe := make(map[Cell]struct{})
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			if (r+c)%2 == 0 {
				universe[Cell{r, c}] = struct{}{}
			}
		}
	}
	u := &GameOfLife{
		universe: universe,
		numRows:  numRows,
		numCols:  numCols,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}

func BenchmarkCreateNextGeneration_1000x1000_Sparse(b *testing.B) {
	// Create a sparse 1000x1000 universe with a few live cells
	numRows, numCols := 1000, 1000
	universe := make(map[Cell]struct{})
	for i := 0; i < 10000; i++ {
		universe[Cell{i % numRows, (i * 31) % numCols}] = struct{}{}
	}
	u := &GameOfLife{
		universe: universe,
		numRows:  numRows,
		numCols:  numCols,
		neighbouringCells: []Cell{
			{R: -1, C: -1}, {R: -1, C: 0}, {R: -1, C: 1},
			{R: 0, C: -1}, {R: 0, C: 1},
			{R: 1, C: -1}, {R: 1, C: 0}, {R: 1, C: 1},
		},
		rules: []Rule{RuleFactory(ConwayRuleType)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}
