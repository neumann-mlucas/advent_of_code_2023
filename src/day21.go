package main

import (
	"fmt"
	"os"
	"strings"
)

const testInputP1 = `
...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........
`

type pos struct {
    i, j int
}

func (p pos) String() string {
    return fmt.Sprintf("<Position x:%d y:%d>", p.i, p.j)
}

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day21.txt"
	}
	//
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1, 6)
	if testResultP1 != 16 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 16)
	}
    fmt.Println(testResultP1)


	// testResultP2 := solveP2(testInputP1)
	// if testResultP2 != 64 {
	// 	fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 64)
	// }

	resultP1 := solveP1(string(content), 64)
	fmt.Println(resultP1)

	resultP2 := solveP1(string(content), 26501365)
	fmt.Println(resultP2)
}

func solveP1(input string, n int) int {
	input = strings.TrimSpace(input)
	gridStr := strings.Split(input, "\n")

	grid := make([][]rune, len(gridStr))
	for i := range gridStr {
		grid[i] = make([]rune, len(gridStr[0]))
		grid[i] = []rune(gridStr[i])
	}

    fmt.Println(findStart(grid))

    positions := map[pos]bool{findStart(grid):true}
    steps := map[pos]bool{}

    for i := 0; i < n; i ++ {
        fmt.Println(i)
        for position := range positions {
            for _, step := range findSteps(grid, position) {
                steps[step] = true
            }
        }
        positions, steps = steps, positions
        clear(steps)
    }
    return len(positions)
}

func findStart(grid [][]rune) pos {
	for i, line := range grid {
		for j, c := range line {
			if c == 'S' {
                return pos{i, j}
			}
		}
	}
	return pos{0, 0}
}

func findSteps(grid [][]rune, p pos) []pos {
    out := []pos{}

    if (p.i + 1) < len(grid) && grid[p.i +1][p.j] != '#' {
        out = append(out, pos{p.i+1, p.j})
    }

    if (p.i - 1) >= 0 && grid[p.i -1][p.j] != '#' {
        out = append(out, pos{p.i-1, p.j})
    }

    if (p.j + 1) < len(grid[0]) && grid[p.i][p.j + 1] != '#' {
        out = append(out, pos{p.i, p.j+1})
    }

    if (p.j - 1) >= 0 && grid[p.i][p.j - 1] != '#' {
        out = append(out, pos{p.i, p.j-1})
    }

    return out
}
