package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const testInputP1 = `
Time:      7  15   30
Distance:  9  40  200
`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day06.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := mulWaystoWin(testInputP1)
	if testResultP1 != 288 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 288)
	}

	testResultP2 := numWaysToWinP2(testInputP1)
	if testResultP2 != 71503 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 71503)
	}

	resultP1 := mulWaystoWin(string(content))
	fmt.Println(resultP1)

	resultP2 := numWaysToWinP2(string(content))
	fmt.Println(resultP2)
}

func mulWaystoWin(input string) int {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	times, distances := parseNumbers(lines[0]), parseNumbers(lines[1])

	total := 1
	for i := 0; i < len(times); i++ {
		total = total * numWaysToWin(times[i], distances[i])

	}
	return total
}

func numWaysToWin(t, s int) int {
	tf, sf := float64(t), float64(s)+0.001
	roots := solveBaskara(-1, tf, -sf)
	r1, r2 := roots[0], roots[1]
	return int(math.Floor(max(r1, r2))-math.Ceil(min(r1, r2))) + 1

}

func numWaysToWinP2(input string) int {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	time, distance := parseNumber(lines[0]), parseNumber(lines[1])
	return numWaysToWin(time, distance)

}

func solveBaskara(a, b, c float64) []float64 {
	delta := b*b - 4*a*c
	r1 := (-b + math.Sqrt(delta)) / 2 * a
	r2 := (-b - math.Sqrt(delta)) / 2 * a
	return []float64{r1, r2}
}

func parseNumbers(str string) []int {
	var numbers []int
	for _, word := range strings.Fields(str) {
		if isNumeric(word) {
			n, _ := strconv.Atoi(word)
			numbers = append(numbers, n)
		}
	}
	return numbers
}

func parseNumber(str string) int {
	var number []rune
	for _, c := range str {
		if unicode.IsDigit(c) {
			number = append(number, c)
		}
	}
	val, _ := strconv.Atoi(string(number))
	return val
}

func isNumeric(str string) bool {
	for _, c := range str {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
