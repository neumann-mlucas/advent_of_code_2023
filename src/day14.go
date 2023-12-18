package main

import (
	"fmt"
	"os"
	"strings"
)

const testInputP1 = `
O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....
`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day14.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 136 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 136)
	}

	testResultP2 := solveP2(testInputP1)
	if testResultP2 != 64 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 64)
	}

	resultP1 := solveP1(string(content))
	fmt.Println(resultP1)

	resultP2 := solveP2(string(content))
	fmt.Println(resultP2)
}

func solveP1(input string) int {
	input = strings.TrimSpace(input)
	gridStr := strings.Split(input, "\n")

	grid := make([][]rune, len(gridStr))
	for i := range gridStr {
		grid[i] = make([]rune, len(gridStr[0]))
		grid[i] = []rune(gridStr[i])
	}

	for i, line := range grid {
		for j, c := range line {
			if c == 'O' {
				moveRockNorth(grid, i, j)
			}
		}
	}
	return calcWeight(grid)
}

func solveP2(input string) int {
	input = strings.TrimSpace(input)
	gridStr := strings.Split(input, "\n")

	grid := make([][]rune, len(gridStr))
	oldGrid := make([][]rune, len(gridStr))
	for i := range gridStr {
		grid[i] = make([]rune, len(gridStr[0]))
		oldGrid[i] = make([]rune, len(gridStr[0]))
		grid[i] = []rune(gridStr[i])
	}

	for n := 0; n <= 100; n++ {
		grid = doCycle(grid)
	}

	var series []int
	for n := 0; n <= 1_000; n++ {
		grid = doCycle(grid)
		series = append(series, calcWeight(grid))
	}

	match := findPeriodicPattern(series)
	return match[(1_000_000_000-101)%(len(match))-1]
}

func calcWeight(grid [][]rune) int {
	var total int
	for i, line := range grid {
		for _, c := range line {
			if c == 'O' {
				total += (len(grid) - i)
			}
		}
	}
	return total
}

func doCycle(grid [][]rune) [][]rune {
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 'O' {
				moveRockNorth(grid, i, j)
			}
		}
	}

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 'O' {
				moveRockWest(grid, i, j)
			}
		}
	}

	for i := len(grid) - 1; i >= 0; i-- {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 'O' {
				moveRockSouth(grid, i, j)
			}
		}
	}

	for i := 0; i < len(grid); i++ {
		for j := len(grid[i]) - 1; j >= 0; j-- {
			if grid[i][j] == 'O' {
				moveRockEast(grid, i, j)
			}
		}
	}

	return grid
}

func moveRockNorth(grid [][]rune, i, j int) {
	for canMove(grid, i-1, j) {
		grid[i][j], grid[i-1][j] = grid[i-1][j], grid[i][j]
		i--
	}
}

func moveRockSouth(grid [][]rune, i, j int) {
	for canMove(grid, i+1, j) {
		grid[i][j], grid[i+1][j] = grid[i+1][j], grid[i][j]
		i++
	}
}

func moveRockWest(grid [][]rune, i, j int) {
	for canMove(grid, i, j-1) {
		grid[i][j], grid[i][j-1] = grid[i][j-1], grid[i][j]
		j--
	}
}

func moveRockEast(grid [][]rune, i, j int) {
	for canMove(grid, i, j+1) {
		grid[i][j], grid[i][j+1] = grid[i][j+1], grid[i][j]
		j++
	}
}

func canMove(grid [][]rune, i, j int) bool {
	inbounds := i >= 0 && i < len(grid) && j >= 0 && j < len(grid[i])
	return inbounds && grid[i][j] == '.'
}

func findPeriodicPattern(slice []int) []int {
	n := len(slice)

	for size := 1; size <= n/2; size++ {
		pattern := make([]int, size)
		copy(pattern, slice[:size])
		pattern = append(pattern, pattern...)
		if compareSlices(pattern, slice[size*2:size*4]) {
			return pattern[:size]
		}
	}
	return nil
}

func compareSlices(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
