# Game of Life exercise
> Implementation of zero player game - https://en.wikipedia.org/wiki/Conway's_Game_of_Life

<details>

## Definition​:

<p>
The universe of the Game of Life is an infinite two­dimensional orthogonal grid of square cells, each of
which is in one of two possible states, alive​ or dead.​ Every cell interacts with its eight neighbours, which
are the cells that are horizontally, vertically, or diagonally adjacent.
</p>


## Rules​:

At each step in time, the following transitions occur:
1. Any live cell with fewer than two live neighbours dies​, as if caused by under­population.
2. Any live cell with two or three live neighbours lives​ on to the next generation.
3. Any live cell with more than three live neighbours dies​, as if by overcrowding.
4. Any dead cell with exactly three live neighbours becomes a live cell,​ as if by reproduction.

<p>
The initial pattern constitutes the seed of the system. The first generation is created by applying the above
rules simultaneously to every cell in the seed—births and deaths occur simultaneously, and the discrete
moment at which this happens is sometimes called a tick (in other words, each generation is a pure
function of the preceding one). The rules continue to be applied repeatedly to create further generations.
</p>


## Objective​:

1. Implement game of life data structures and algorithm
2. Demonstrate that game of life algorithm works

> Note: the program has to run and work properly (working program is better than in­progress design).

> Note: use a ‘Glider’ pattern placed in the middle of 25x25 cell universe for this exercise.

## Guidelines:

1. Please limit yourself to no more than 2 hours for this exercise.
2. Use any language and/or frameworks you’d like.
3. No actual “UI” is required ­ console output is fine.
4. Be ready to explain your design decisions and how you might improve/expand upon your solution.
5. Please submit your solution using Github or Dropbox or Google Drive or etc.
6. Include any setup details as needed to make your solution run.
7. Please email us if you have any questions.

</details>

## Implementation Details

#### How to run?

1. Checkout latest code from github
    ```bash
    git clone https://github.com/dilipvaidya/game-of-life.git
    cd game-of-life
    ```

2. Take help to learn the command line arguments being accepted by the program
    ```shell 
    game-of-life % go run . -help                            
    Usage of /var/folders/xt/l7kq1dlx3fj8s3p15prcpy9w0000gn/T/go-build2668948479/b001/exe/game-of-life:
    This is a custom usage message.
      -cols int
            Number of columns in the universe (default 5)
      -rows int
            Number of rows in the universe (default 5)
      -seed string
            Seed pattern for the universe (default, glider) (default "default")
    ```

 3. execute with appropriate command line values
    ```shell
    game-of-life % go run .
    ```


### Brute Force algorithm
<details>

1. Let's call a game grid as `universe` which is a 2-dimentional array of integers with possibility of only binary values - either '1' (live) or '0' (dead). Every time tick will hold the current copy of the universe.
2. While producing next generation at next tick, a new universe will be constrcuted of same original size. Original/previois snapshot of universe will be traversed and the game rules will be applied to find the possible values for new universe (cells are dead or alive)
3. Before calculating next generation of universe, a original/previous universe will be traversed to find an alive and dead neighbours of each of the cell (8 neighbours) which then further will be traverse to finalize the cell values in new universe.

#### Time and Space Complexity
Assume there are `n` number of rows and `m` number of columns in the universe
1. Time complexity: 
> O(n x m x 8) -> O(n x m); for every cell, 8 neighbours will be travelled in constant time.
2. Space Complexity:
> O(n x m x 8) -> O(n x m); for every cell, 8 neighbours will be travelled in constant time

#### limitations
1. This solution won't scale well for larger sparse grid with values of `n` and `m` in 1000.

</details>

### Optimized algorithm [Selected and currently implemented]
<details>

> Hint: What if the universe grid is sparse with alive cells?

1. Let's call a game grid as a `universe` which is optimized to store only alive cells - set. Only cell that is maintaining its aliveness from previous generation or the cell that is reviving from dead will be added into the set.
2. While producing next generation universe, algorithm will - 
- generate a map of the neighbours to the alive node as key and count of their alive neighbours as value, 
- iterate over this map to find out if it can be revived (currently dead with 3 live neighbours), and add into set if so. 
- iterate over the alive cell's set to find if any of them will remain alive or will die per their count of alive neighbours (map above).

#### Time and Space Complexity

Assume there are `l` alive cells in the universe at current generation and alive cells are controlled to max `l`
1. Time complexity: 
> O(l x 8) -> O(l); for every cell, 8 neighbours will be travelled in constant time
2. Space Complexity: 
> O(l x 8) -> O(l); for every cell, 8 neighbours will be travelled in constant time

#### performance benchmark:

Performance benchmarks below proves, with optimized data structure and algorithm, running time is totally depend on the number of alive cells and not on the size of the universe grid.

1. Benchmark: 100X100 grid with glider pattern (5 out of 10000 cells alive)
    ```shell
    Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkCreateNextGeneration_100x100_Glider$ github.com/dilipvaidya/game-of-life/gameoflife

    goos: darwin
    goarch: amd64
    pkg: github.com/dilipvaidya/game-of-life/gameoflife
    cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
    BenchmarkCreateNextGeneration_100x100_Glider-12    	  615019	      1882 ns/op	     705 B/op	       7 allocs/op
    PASS
    ok  	github.com/dilipvaidya/game-of-life/gameoflife	2.479s
    ```

2. 1000X1000 grid with glider pattern (5 out of 1000000 cells alive)
    ```shell
    Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkCreateNextGeneration_1000x1000_Glider$ github.com/dilipvaidya/game-of-life/gameoflife

    goos: darwin
    goarch: amd64
    pkg: github.com/dilipvaidya/game-of-life/gameoflife
    cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
    BenchmarkCreateNextGeneration_1000x1000_Glider-12    	  636194	      1896 ns/op	     707 B/op	       7 allocs/op
    PASS
    ok  	github.com/dilipvaidya/game-of-life/gameoflife	2.582s
    ```

3. 100X100 grid with 50% live cells
    ```shell
    Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkCreateNextGeneration_100x100$ github.com/dilipvaidya/game-of-life/gameoflife

    goos: darwin
    goarch: amd64
    pkg: github.com/dilipvaidya/game-of-life/gameoflife
    cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
    BenchmarkCreateNextGeneration_100x100-12    	    4142	    276222 ns/op	  168596 B/op	     114 allocs/op
    PASS
    ok  	github.com/dilipvaidya/game-of-life/gameoflife	3.065s
    ```

4. 1000X1000 grid with 50% live cells
    ```shell
    Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkCreateNextGeneration_1000x1000$ github.com/dilipvaidya/game-of-life/gameoflife

    goos: darwin
    goarch: amd64
    pkg: github.com/dilipvaidya/game-of-life/gameoflife
    cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
    BenchmarkCreateNextGeneration_1000x1000-12    	      26	  39910036 ns/op	 8812968 B/op	    3033 allocs/op
    PASS
    ok  	github.com/dilipvaidya/game-of-life/gameoflife	10.706s
    ```

</details>

### Next steps
1. Performance improvment: Is it possible to use in-place updates to avoid runtime memory allocation? 
- every time next generation of the universe is created, new copy is being written and old will be garbage collected. 
- in place updates will avoind run time memory allocation drastically. 
- `sync.pool` can be used to pre-allocate active and passive universe/memory blocks which will be alternatively used while generating next generation which will reduce the runtime memory allocation drastically. Cleaup of the passive universe can be handed over to another thread. 
2. add UI component, a canvas kind, which will be updated for every next universe generation rather than displaying universe on the console - not a goal of exercise.