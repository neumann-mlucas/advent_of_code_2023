package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const testInputP1 = `rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day15.txt"
	}
	//
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 1320 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 1320)
	}

	testResultP2 := solveP2(testInputP1)
	if testResultP2 != 145 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 145)
	}

	resultP1 := solveP1(string(content))
	fmt.Println(resultP1)

	resultP2 := solveP2(string(content))
	fmt.Println(resultP2)
}

func solveP1(input string) int {
	input = strings.TrimSpace(input)
	codes := strings.Split(input, ",")

	var total int
	for _, c := range codes {
		total += hash(c)

	}

	return total
}

func solveP2(input string) int {
	input = strings.TrimSpace(input)
	codes := strings.Split(input, ",")

	boxes := make([][]string, 256)
	lenses := make([][]int, 256)

	for i := 0; i < 256; i++ {
		boxes[i] = make([]string, 0, 256)
		lenses[i] = make([]int, 0, 256)
	}

	var total int
	for _, c := range codes {
		label := c[:strings.IndexAny(c, "-=")]
		id := hash(label)
		labelIdx := index(boxes[id], label)

		if strings.ContainsRune(c, '-') {
			if labelIdx != -1 {
				boxes[id] = append(boxes[id][:labelIdx], boxes[id][labelIdx+1:]...)
				lenses[id] = append(lenses[id][:labelIdx], lenses[id][labelIdx+1:]...)
			}
		} else if strings.ContainsRune(c, '=') {
			lens, _ := strconv.Atoi(c[strings.IndexRune(c, '=')+1:])
			if labelIdx == -1 {
				boxes[id] = append(boxes[id], label)
				lenses[id] = append(lenses[id], lens)
			} else {
				boxes[id][labelIdx] = label
				lenses[id][labelIdx] = lens

			}
		}

	}
	for i, box := range boxes {
		for j := range box {
			total += (i + 1) * (j + 1) * lenses[i][j]
		}

	}

	return total
}

func hash(input string) int {
	var total int
	for _, ch := range input {
		total += int(ch)
		total *= 17
		total = total % 256
	}
	return total
}

func index(input []string, elem string) int {
	for i, s := range input {
		if s == elem {
			return i
		}
	}
	return -1
}
