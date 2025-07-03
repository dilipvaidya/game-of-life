package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dilipvaidya/game-of-life/gameoflife"
)

// main function initializes the Game of Life universe based on user input flags
func main() {

	// Custom usage function
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("This is a custom usage message.")
		flag.PrintDefaults() // Prints default flag usage
	}

	// defining flags to accept user input
	seedPatternStr := flag.String("seed", gameoflife.Default.String(), "Seed pattern for the universe (default, glider)")
	rows := flag.Int("rows", 5, "Number of rows in the universe")
	cols := flag.Int("cols", 5, "Number of columns in the universe")

	// Parse the command line flags
	flag.Parse()

	// Validate seed pattern
	var seedPattern gameoflife.SEED_PATTERN
	switch *seedPatternStr {
	case gameoflife.Default.String():
		seedPattern = gameoflife.Default
	case gameoflife.Glider.String():
		seedPattern = gameoflife.Glider
	default:
		seedPattern = gameoflife.Default
	}

	// Create the Game of Life universe with the specified seed pattern and dimensions
	game := gameoflife.CreateSeedUniverse(*rows, *cols, seedPattern)
	game.Run(25, 500*time.Millisecond)
}
