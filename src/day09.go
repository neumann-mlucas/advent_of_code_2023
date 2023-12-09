package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const testInputP1 = `
0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`

func main() {

	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day09.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 114 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 114)
	}

	testResultP2 := solveP2(testInputP1)
	if testResultP2 != 2 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 2)
	}

	resultP1 := solveP1(string(content))
	fmt.Println(resultP1)

	resultP2 := solveP2(string(content))
	fmt.Println(resultP2)
}

func solveP1(input string) int {
	input = strings.TrimSpace(input)

	var total int
	for _, line := range strings.Split(input, "\n") {
		curve := parseNumbers(line)
		total += predictNext(curve)
	}

	return total
}

func solveP2(input string) int {
	input = strings.TrimSpace(input)

	var total int
	for _, line := range strings.Split(input, "\n") {
		curve := parseNumbers(line)
		total += predictPrev(curve)
		fmt.Println()
	}

	return total
}

func predictNext(curve []int) int {
	if isDxContant(curve) {
		return 0
	}
	dx := DerivateDiscrete(curve)
	return curve[len(curve)-1] + predictNext(dx)

}

func predictPrev(curve []int) int {
	fmt.Println(curve)
	if isDxContant(curve) {
		return 0
	}
	dx := DerivateDiscrete(curve)
	return curve[0] - predictPrev(dx)
}

func DerivateDiscrete(curve []int) []int {
	var dx []int
	for i, n := range curve[1:] {
		dx = append(dx, n-curve[i])

	}
	return dx
}

func isDxContant(curve []int) bool {
	for _, n := range curve {
		if n != 0 {
			return false
		}
	}
	return true
}

func parseNumbers(str string) []int {
	var numbers []int
	for _, word := range strings.Fields(str) {
		n, _ := strconv.Atoi(word)
		numbers = append(numbers, n)
	}
	return numbers
}
