package main

import (
	"fmt"
	"os"
	"strings"
)

const testInputP1 = `
#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day12.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solve(testInputP1, 0)
	if testResultP1 != 405 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 405)
	}

	testResultP2 := solve(testInputP1, 1)
	if testResultP2 != 400 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 400)
	}

	resultP1 := solve(string(content), 0)
	fmt.Println(resultP1)

	resultP2 := solve(string(content), 1)
	fmt.Println(resultP2)
}

func solve(input string, tol int) int {
	input = strings.TrimSpace(input)
	grids := strings.Split(input, "\n\n")

	var total int
	for _, g := range grids {
		grid := strings.Split(g, "\n")
		n := findReflectionV(grid, tol)
		if n == 0 {
			n = findReflectionH(grid, tol)
		}
		total += n
	}

	return total
}

func findReflectionV(grid []string, tol int) int {

	for col := range grid[0] {
		if hasReflectionLine(grid, col, tol) {
			return col
		}
	}
	return 0
}

func findReflectionH(str []string, tol int) int {
	str = transposeStrings(str)
	val := findReflectionV(str, tol)
	return val * 100
}

func hasReflectionLine(grid []string, col, tol int) bool {
	var totalDiff int
	for _, line := range grid {
		left, right := reverse(line[:col]), line[col:]
		rlen := min(len(left), len(right))
		if rlen == 0 {
			return false
		}
		totalDiff += countDiff(left[:rlen], right[:rlen])
	}

	return totalDiff == tol
}

func countDiff(this, other string) int {
	var count int
	for i := 0; i < min(len(this), len(other)); i++ {
		if this[i] != other[i] {
			count++
		}
	}
	return count
}

func transposeStrings(str []string) []string {
	if len(str) == 0 {
		return []string{}
	}

	width := len(str[0])
	transposed := make([]string, width)

	for _, row := range str {
		for j, ch := range row {
			transposed[j] += string(ch)
		}
	}
	return transposed
}

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}
