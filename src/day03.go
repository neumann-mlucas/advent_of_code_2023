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

	testResultP2 := GetSumGearRatios(testInputP1)
	if testResultP2 != 467835 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d", testResultP2, 467835)
	}

	input := strings.TrimSpace(string(content))

	resultP1 := SumParts(input)
	fmt.Println(resultP1)

	resultP2 := GetSumGearRatios(input)
	fmt.Println(resultP2)
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

func (n *Number) isNeighour(i, j int) bool {
	if n.line < i-1 || n.line > i+1 {
		return false
	}

	for _, idx := range n.idxs {
		if idx >= j-1 && idx <= j+1 {
			return true
		}
	}
	return false
}

type Coord struct {
	i int
	j int
}

func FindGears(grid []string, out chan Coord) {
	for i, line := range grid {
		for j, char := range line {
			if char == '*' {
				out <- Coord{i, j}
			}
		}
	}
	close(out)
}

func FindAdjacentNumber(grid []string, inp chan Coord, out chan []Number) {
	for coord := range inp {
		var numbers []Number
		// find numbers in line above Coord line and in lines above and  bellow
		for inc := -1; inc <= 1; inc++ {
			// check grid bounds
			if coord.i+inc >= 0 && coord.i+inc < len(grid) {
				// find numbers in line
				n := FindNumbers(coord.i+inc, grid[coord.i+inc])
				numbers = append(numbers, n...)
			}
		}

		var pair []Number
		for _, n := range numbers {
			// check if number is neighour of coord
			if n.isNeighour(coord.i, coord.j) {
				pair = append(pair, n)
			}

		}
		// send to next function
		out <- pair
	}
	close(out)
}

func SumGearRatios(inp chan []Number, out chan int) {
	var total int
	for pair := range inp {
		if len(pair) == 2 {
			total += pair[0].toInt() * pair[1].toInt()
		}
	}
	out <- total
	close(out)
}

func GetSumGearRatios(engine string) int {

	engine = strings.TrimSpace(engine)
	grid := strings.Split(engine, "\n")

	coordChan := make(chan Coord)
	numChan := make(chan []Number)
	intChan := make(chan int)

	go FindGears(grid, coordChan)

	go FindAdjacentNumber(grid, coordChan, numChan)

	go SumGearRatios(numChan, intChan)

	total := <-intChan
	return total

}
