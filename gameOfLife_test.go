package main

import (
	"reflect"
	"testing"
)

func TestDeepCopyUniverse(t *testing.T) {
	tests := []struct {
		name     string
		original [][]int
		modify   func(copy, original [][]int) bool
	}{
		{
			name:     "Empty",
			original: [][]int{},
			modify: func(copy, original [][]int) bool {
				return &original == &copy
			},
		},
		{
			name:     "SingleRow",
			original: [][]int{{1, 2, 3}},
			modify: func(copy, original [][]int) bool {
				if &original[0] == &copy[0] {
					return true
				}
				copy[0][0] = 99
				return original[0][0] == 99
			},
		},
		{
			name: "MultipleRows",
			original: [][]int{
				{1, 0, 1},
				{0, 1, 0},
			},
			modify: func(copy, original [][]int) bool {
				for i := range original {
					if &original[i] == &copy[i] {
						return true
					}
				}
				copy[1][2] = 42
				return original[1][2] == 42
			},
		},
		// {
		// 	name:     "NilRows",
		// 	original: [][]int{nil, {1, 2}},
		// 	modify: func(copy, original [][]int) bool {
		// 		return copy[0] != nil
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copy := _deepCopyUniverse(tt.original)
			if !reflect.DeepEqual(tt.original, copy) {
				t.Errorf("Expected %v, got %v", tt.original, copy)
			}
			if tt.modify != nil && tt.modify(copy, tt.original) {
				t.Errorf("Deep copy failed for test case %s", tt.name)
			}
		})
	}
}

func TestIsValidNeighbour(t *testing.T) {
	universe := [][]int{
		{1, 0, 1},
		{0, 1, 0},
		{1, 1, 1},
	}
	gol := CreateUniverse(universe)

	tests := []struct {
		name  string
		row   int
		col   int
		valid bool
	}{
		{"TopLeft", 0, 0, true},
		{"TopRight", 0, 2, true},
		{"BottomLeft", 2, 0, true},
		{"BottomRight", 2, 2, true},
		{"NegativeRow", -1, 1, false},
		{"NegativeCol", 1, -1, false},
		{"RowTooLarge", 3, 1, false},
		{"ColTooLarge", 1, 3, false},
		{"BothNegative", -1, -1, false},
		{"BothTooLarge", 3, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gol._isValidNeighbour(tt.row, tt.col)
			if got != tt.valid {
				t.Errorf("Expected %v for (%d,%d), got %v", tt.valid, tt.row, tt.col, got)
			}
		})
	}
}

// TestGameOfLife__getNeighbourStats tests the _getNeighbourStats method of GameOfLife
// It checks the number of live and dead neighbours for various cell positions in the universe.
func TestGameOfLife__getNeighbourStats(t *testing.T) {
	neighbours := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	tests := []struct {
		name     string
		universe [][]int
		row      int
		col      int
		wantLive int
		wantDead int
	}{
		{
			name: "center cell, all live neighbours",
			universe: [][]int{
				{1, 1, 1},
				{1, 0, 1},
				{1, 1, 1},
			},
			row:      1,
			col:      1,
			wantLive: 8,
			wantDead: 0,
		},
		{
			name: "center cell, all dead neighbours",
			universe: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			row:      1,
			col:      1,
			wantLive: 0,
			wantDead: 8,
		},
		{
			name: "corner cell, some out of bounds",
			universe: [][]int{
				{1, 0, 1},
				{0, 1, 0},
				{1, 0, 1},
			},
			row:      0,
			col:      0,
			wantLive: 1, // (0,1) is 0, (1,0) is 0, (1,1) is 1
			wantDead: 2,
		},
		{
			name: "edge cell, mixed neighbours",
			universe: [][]int{
				{1, 1, 0},
				{0, 1, 1},
				{1, 0, 0},
			},
			row:      0,
			col:      1,
			wantLive: 3, // (0,0)=2, (0, 1)=3, (0,2)=3, (1,0)=4, (1,1)=4, (1,2)=2, (2,0)=1, (2,1)=3, (2,2)=2
			wantDead: 2,
		},
		{
			name: "single cell universe",
			universe: [][]int{
				{1},
			},
			row:      0,
			col:      0,
			wantLive: 0,
			wantDead: 0,
		},
		{
			name: "non-square universe, edge cell",
			universe: [][]int{
				{1, 0, 1, 0},
				{0, 1, 0, 1},
			},
			row:      0,
			col:      2,
			wantLive: 2, // (0,1)=0, (0,3)=0, (1,1)=1, (1,2)=0, (1,3)=1
			wantDead: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &GameOfLife{
				universe:          tt.universe,
				numRows:           len(tt.universe),
				numCols:           len(tt.universe[0]),
				neighbouringCells: neighbours,
			}
			gotLive, gotDead := u._getNeighbourStats(tt.row, tt.col)
			if gotLive != tt.wantLive || gotDead != tt.wantDead {
				t.Errorf("got (%d, %d), want (%d, %d)", gotLive, gotDead, tt.wantLive, tt.wantDead)
			}
		})
	}
}

// TestPlayNextGeneration tests the PlayNextGeneration method of GameOfLife
// It checks the next generation of the universe based on various initial configurations.
// Each test case includes an initial universe and the expected next generation.
func TestPlayNextGeneration(t *testing.T) {
	tests := []struct {
		name         string
		initial      [][]int
		expectedNext [][]int
	}{
		{
			name: "Block (Still Life)",
			initial: [][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 0},
				{0, 1, 1, 0},
				{0, 0, 0, 0},
			},
			expectedNext: [][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 0},
				{0, 1, 1, 0},
				{0, 0, 0, 0},
			},
		},
		{
			name: "Blinker (Oscillator)",
			initial: [][]int{
				{0, 0, 0, 0, 0},
				{0, 0, 1, 0, 0},
				{0, 0, 1, 0, 0},
				{0, 0, 1, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expectedNext: [][]int{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 1, 1, 1, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		},
		{
			name: "Toad (Oscillator)",
			initial: [][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 1},
				{1, 1, 1, 0},
				{0, 0, 0, 0},
			},
			expectedNext: [][]int{
				{0, 0, 1, 0},
				{1, 0, 0, 1},
				{1, 0, 0, 1},
				{0, 1, 0, 0},
			},
		},
		{
			name: "All Dead",
			initial: [][]int{
				{0, 0},
				{0, 0},
			},
			expectedNext: [][]int{
				{0, 0},
				{0, 0},
			},
		},
		{
			name: "Single Live Cell Dies",
			initial: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			expectedNext: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
		},
		{
			name: "Reproduction",
			initial: [][]int{
				{0, 1, 0},
				{1, 0, 1},
				{0, 0, 0},
			},
			expectedNext: [][]int{
				{0, 1, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gol := CreateUniverse(tt.initial)
			gol.PlayNextGeneration()
			if !reflect.DeepEqual(gol.universe, tt.expectedNext) {
				t.Errorf("After PlayNextGeneration, got:\n%v\nwant:\n%v", gol.universe, tt.expectedNext)
			}
		})
	}
}
