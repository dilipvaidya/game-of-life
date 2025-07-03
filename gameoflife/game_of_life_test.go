package gameoflife

import (
	"testing"
)

func Test_isValidNeighbour(t *testing.T) {
	u := &GameOfLife{
		numRows: 5,
		numCols: 5,
	}

	tests := []struct {
		name     string
		row, col int
		want     bool
	}{
		{"top-left-corner", 0, 0, true},
		{"bottom-right-corner", 4, 4, true},
		{"center", 2, 2, true},
		{"negative-row", -1, 0, false},
		{"negative-col", 0, -1, false},
		{"row-out-of-bounds", 5, 0, false},
		{"col-out-of-bounds", 0, 5, false},
		{"both-out-of-bounds", 5, 5, false},
		{"both-negative", -1, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u._isValidNeighbour(tt.row, tt.col)
			if got != tt.want {
				t.Errorf("_isValidNeighbour(%d, %d) = %v; want %v", tt.row, tt.col, got, tt.want)
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
		wantUniverse map[[2]int]bool
	}{
		{
			name: "Default pattern 5x5",
			args: args{5, 5, Default},
			wantUniverse: map[[2]int]bool{
				{1, 2}: true, {2, 2}: true, {3, 2}: true,
			},
		},
		{
			name: "Glider pattern 5x5",
			args: args{5, 5, Glider},
			wantUniverse: map[[2]int]bool{
				{2, 3}: true, {3, 4}: true, {4, 2}: true, {4, 3}: true, {4, 4}: true,
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
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}

	tests := []struct {
		name           string
		currentPos     [2]int
		wantAliveCells map[[2]int]int
	}{
		{
			name:       "center cell",
			currentPos: [2]int{1, 1},
			wantAliveCells: map[[2]int]int{
				{0, 0}: 1, {0, 1}: 1, {0, 2}: 1,
				{1, 0}: 1, {1, 2}: 1,
				{2, 0}: 1, {2, 1}: 1, {2, 2}: 1,
			},
		},
		{
			name:       "top-left corner",
			currentPos: [2]int{0, 0},
			wantAliveCells: map[[2]int]int{
				{0, 1}: 1, {1, 0}: 1, {1, 1}: 1,
			},
		},
		{
			name:       "bottom-right corner",
			currentPos: [2]int{2, 2},
			wantAliveCells: map[[2]int]int{
				{1, 1}: 1, {1, 2}: 1, {2, 1}: 1,
			},
		},
		{
			name:       "edge cell",
			currentPos: [2]int{0, 1},
			wantAliveCells: map[[2]int]int{
				{0, 0}: 1, {0, 2}: 1,
				{1, 0}: 1, {1, 1}: 1, {1, 2}: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliveMap := make(map[[2]int]int)
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

func TestCreateNextGeneration_Blinker(t *testing.T) {
	// Blinker pattern (period 2 oscillator)
	// Generation 0 (vertical):
	// . . .
	// X X X
	// . . .
	initialUniverse := map[[2]int]struct{}{
		{1, 0}: struct{}{},
		{1, 1}: struct{}{},
		{1, 2}: struct{}{},
	}
	u := &GameOfLife{
		universe: initialUniverse,
		numRows:  3,
		numCols:  3,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
	u.CreateNextGeneration()
	wantUniverse := map[[2]int]bool{
		{0, 1}: true,
		{1, 1}: true,
		{2, 1}: true,
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

func TestCreateNextGeneration_BlockStillLife(t *testing.T) {
	// Block pattern (still life)
	// . X X
	// . X X
	initialUniverse := map[[2]int]struct{}{
		{0, 1}: struct{}{}, {0, 2}: struct{}{},
		{1, 1}: struct{}{}, {1, 2}: struct{}{},
	}
	u := &GameOfLife{
		universe: initialUniverse,
		numRows:  4,
		numCols:  4,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
	u.CreateNextGeneration()
	wantUniverse := map[[2]int]bool{
		{0, 1}: true, {0, 2}: true,
		{1, 1}: true, {1, 2}: true,
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
	initialUniverse := map[[2]int]struct{}{
		{1, 1}: struct{}{},
	}
	u := &GameOfLife{
		universe: initialUniverse,
		numRows:  3,
		numCols:  3,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
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
	initialUniverse := map[[2]int]struct{}{
		{0, 1}: struct{}{},
		{1, 0}: struct{}{},
		{1, 2}: struct{}{},
	}
	u := &GameOfLife{
		universe: initialUniverse,
		numRows:  3,
		numCols:  3,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
	u.CreateNextGeneration()
	wantUniverse := map[[2]int]bool{
		{0, 1}: true,
		{1, 1}: true, // new cell born
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
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
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
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}

func BenchmarkCreateNextGeneration_100x100(b *testing.B) {
	// Create a dense 100x100 universe with a checkerboard pattern
	numRows, numCols := 100, 100
	universe := make(map[[2]int]struct{})
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			if (r+c)%2 == 0 {
				universe[[2]int{r, c}] = struct{}{}
			}
		}
	}
	u := &GameOfLife{
		universe: universe,
		numRows:  numRows,
		numCols:  numCols,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}

func BenchmarkCreateNextGeneration_1000x1000(b *testing.B) {
	// Create a dense 1000x1000 universe with a checkerboard pattern
	numRows, numCols := 1000, 1000
	universe := make(map[[2]int]struct{})
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			if (r+c)%2 == 0 {
				universe[[2]int{r, c}] = struct{}{}
			}
		}
	}
	u := &GameOfLife{
		universe: universe,
		numRows:  numRows,
		numCols:  numCols,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}

func BenchmarkCreateNextGeneration_1000x1000_Sparse(b *testing.B) {
	// Create a sparse 1000x1000 universe with a few live cells
	numRows, numCols := 1000, 1000
	universe := make(map[[2]int]struct{})
	for i := 0; i < 10000; i++ {
		universe[[2]int{i % numRows, (i * 31) % numCols}] = struct{}{}
	}
	u := &GameOfLife{
		universe: universe,
		numRows:  numRows,
		numCols:  numCols,
		neighbouringCells: [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.CreateNextGeneration()
	}
}
