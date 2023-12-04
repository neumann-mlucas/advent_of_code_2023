package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const testInputP1 = `
467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day03.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := SumParts(testInputP1)
	if testResultP1 != 4361 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d", testResultP1, 4361)
	}

	// testResultP2 := CalcPowerSet(testInputP1)
	// if testResultP2 != 2286 {
	// 	fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d", testResultP2, 2286)
	// }

	inputP1 := strings.TrimSpace(string(content))
	resultP1 := SumParts(inputP1)

	fmt.Println(resultP1)

	// inputP2 := strings.TrimSpace(string(content))
	// resultP2 := CalcPowerSet(inputP2)
	// fmt.Println(resultP2)
}

func SumParts(engine string) int {
	var numbers []Number
	engine = strings.TrimSpace(engine)
	engineList := strings.Split(engine, "\n")

	allowed := toAllowedPossitions(engineList)

	for lineNo, line := range engineList {
		n := FindNumbers(lineNo, line)
		numbers = append(numbers, n...)
	}

	var total int
	for _, number := range numbers {
		if number.isAllowed(allowed) {
			total += number.toInt()
		}

	}

	return total
}

func toAllowedPossitions(grid []string) [][]bool {
	N, M := len(grid), len(grid[0])

	arr := make([][]bool, N)
	for i := range arr {
		arr[i] = make([]bool, M)
	}

	for i, line := range grid {
		for j, char := range line {
			if !strings.ContainsRune("1234567890.", char) {
				setNeighbors(arr, i, j)
			}
		}
	}

	return arr
}

func setNeighbors(arr [][]bool, i, j int) {
	for _, m := range []int{i - 1, i, i + 1} {
		if m < 0 || m >= len(arr) {
			continue
		}
		for _, n := range []int{j - 1, j, j + 1} {
			if n < 0 || n >= len(arr[m]) {
				continue
			}
			arr[m][n] = true
		}
	}
}

func FindNumbers(lineNo int, line string) []Number {
	var numbers []Number
	var currNumber Number

	for i, c := range line {
		if strings.ContainsRune("1234567890", c) {
			currNumber.idxs = append(currNumber.idxs, i)
			currNumber.chrs = append(currNumber.chrs, c)
		} else {
			currNumber.line = lineNo
			if len(currNumber.idxs) >= 1 {
				numbers = append(numbers, currNumber)
			}
			currNumber = Number{}
		}
	}

	if len(currNumber.idxs) >= 1 {
		currNumber.line = lineNo
		numbers = append(numbers, currNumber)
	}
	return numbers
}

type Number struct {
	line int
	chrs []rune
	idxs []int
}

func (n *Number) toInt() int {
	value, err := strconv.Atoi(string(n.chrs))
	if err != nil {
		return 0
	} else {
		return value
	}
}

func (n *Number) isAllowed(grid [][]bool) bool {
	for _, idx := range n.idxs {
		if grid[n.line][idx] {
			return true
		}
	}
	return false
}
