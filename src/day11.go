package main

import (
	"fmt"
	"os"
	"strings"
)

const testInputP1 = `
...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....
`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day11.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solve(testInputP1, 2)
	if testResultP1 != 374 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 374)
	}

	testResultP2 := solve(testInputP1, 100)
	if testResultP2 != 8410 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 8410)
	}

	resultP1 := solve(string(content), 2)
	fmt.Println(resultP1)

	resultP2 := solve(string(content), 1_000_000)
	fmt.Println(resultP2)
}

func solve(input string, mul int) int {
	input = strings.TrimSpace(input)
	input = markGapV(input)
	input = markGapH(input)

	grid := strings.Split(input, "\n")
	starPos := findStarts(grid, mul)
	starPairs := combinations(starPos, 2)

	var total int
	for _, c := range starPairs {
		dist := Abs(c[0].i-c[1].i) + Abs(c[0].j-c[1].j)
		total += dist
	}

	return total

}

func solveP2(input string) int {
	input = strings.TrimSpace(input)
	input = markGapV(input)
	input = markGapH(input)

	grid := strings.Split(input, "\n")
	starPos := findStarts(grid, 100)
	starPairs := combinations(starPos, 2)

	var total int
	for _, c := range starPairs {
		dist := Abs(c[0].i-c[1].i) + Abs(c[0].j-c[1].j)
		fmt.Println(c, dist)
		total += dist
	}
	return total

}

func markGapV(str string) string {
	allDots := func(l string) bool {
		for _, c := range l {
			if c != '.' && c != 'G' {
				return false
			}
		}
		return true
	}

	expanded := ""
	for _, line := range strings.Split(str, "\n") {
		if allDots(line) {
			expanded += strings.Replace(line, ".", "G", -1) + "\n"
		} else {
			expanded += line + "\n"
		}
	}
	return expanded
}

func markGapH(str string) string {
	str = transposeStrings(str)
	str = markGapV(str)
	str = transposeStrings(str)
	return str
}

func findStarts(grid []string, mul int) []Coord {
	countSpaces := func(str string) int {
		return mul*strings.Count(str, "G") + strings.Count(str, ".") + strings.Count(str, "#")
	}
	var coords []Coord
	for i, line := range grid {
		for j, c := range line {
			if c == '#' {
				ii, jj := line[:j], getCol(grid, j)[:i]
				coords = append(coords, Coord{countSpaces(ii), countSpaces(jj)})
			}
		}
	}
	return coords
}

func getCol(str []string, ncol int) string {
	var col []rune
	for _, line := range str {
		for j, c := range line {
			if j == ncol {
				col = append(col, c)
			}
		}
	}
	return string(col)
}

func combinations(input []Coord, n int) [][]Coord {
	var result [][]Coord
	var recurse func(int, []Coord)
	recurse = func(start int, acc []Coord) {
		if len(acc) == n {
			result = append(result, append([]Coord{}, acc...))
			return
		}
		for i := start; i <= len(input)-n+len(acc); i++ {
			recurse(i+1, append(acc, input[i]))
		}
	}
	recurse(0, []Coord{})
	return result
}

func Abs(i int) int {
	mask := i >> (64 - 1)
	return (i + mask) ^ mask
}

func transposeStrings(str string) string {
	input := strings.Split(str, "\n")
	if len(input) == 0 {
		return ""
	}

	width := len(input[0])
	transposed := make([]string, width)

	for _, row := range input {
		for j, ch := range row {
			transposed[j] += string(ch)
		}
	}
	return strings.Join(transposed, "\n")
}

type Coord struct {
	i int
	j int
}
