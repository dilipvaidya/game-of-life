# Game of Life exercise
> Implementation of zero player game - https://en.wikipedia.org/wiki/Conway's_Game_of_Life

## Definition​:
<details>

<p>
The universe of the Game of Life is an infinite two­dimensional orthogonal grid of square cells, each of
which is in one of two possible states, alive​ or dead.​ Every cell interacts with its eight neighbours, which
are the cells that are horizontally, vertically, or diagonally adjacent.
</p>

</details>

## Rules​:
<details>

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

</details>

## Objective​:
<details>

1. Implement game of life data structures and algorithm
2. Demonstrate that game of life algorithm works

> Note: the program has to run and work properly (working program is better than in­progress design).

> Note: use a ‘Glider’ pattern placed in the middle of 25x25 cell universe for this exercise.
</details>

## Guidelines:
<details>

1. Please limit yourself to no more than 2 hours for this exercise.
2. Use any language and/or frameworks you’d like.
3. No actual “UI” is required ­ console output is fine.
4. Be ready to explain your design decisions and how you might improve/expand upon your solution.
5. Please submit your solution using Github or Dropbox or Google Drive or etc.
6. Include any setup details as needed to make your solution run.
7. Please email us if you have any questions.
</details>

## Implementation Details
How to run?
```bash
git clone https://github.com/dilipvaidya/game-of-life.git
cd game-of-life
go run .
```


#### Brute Force algorithm
<details>

1. Let's call a game grid as `universe` which is a 2-dimentional array of integers with possibility of only binary values - either '1' (live) or '0' (dead). Every time tick will hold the current copy of the universe.
2. While producing next generation at next tick, a new universe will be constrcuted of same original size. Original/previois snapshot of universe will be traversed and the game rules will be applied to find the possible values for new universe (cells are dead or alive)
3. Before calculating next generation of universe, a original/previous universe will be traversed to find an alive and dead neighbours of each of the cell (8 neighbours) which then further will be traverse to finalize the cell values in new universe.

##### Time and Space Complexity

<p>

Assume there are `n` number of rows and `m` number of columns in the universe
1. Time complexity: 
> O(n x m x 8) -> O(n x m); 8 is the time
2. Space Complexity:
> O(n x m) 

##### limitations
1. This solution won't scale well for bigger values of `n` and `m` like 1000.

</details>

#### Optimized algorithm
<details>

###### ideas
1. use better data structure. Universe has only binary values, could that be used?

</details>

### Next steps
1. find better data structure to hold the universe/grid to improve time and space complexity.
2. Write performance test suite.
3. User of glider pattern to generate the seed universe.
4. Modulerize and clean code