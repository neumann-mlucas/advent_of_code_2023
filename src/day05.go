package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const testInputP1 = `
seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day05.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := findLowestLocation(testInputP1)
	if testResultP1 != 35 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 35)
	}

	// testResultP2 := findLowestLocationP2(testInputP1)
	// if testResultP2 != 46 {
	// 	fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 46)
	// }

	resultP1 := findLowestLocation(string(content))
	fmt.Println(resultP1)

	// resultP2 := findLowestLocationP2(string(content))
	// fmt.Println(resultP2)
}

func findLowestLocation(maps string) int {
	maps = strings.TrimSpace(maps)

	seeds := parseNumbers(maps[:strings.Index(maps, "\n\n")])

	var maplist [][]string
	for _, m := range strings.Split(maps, "\n\n")[1:] {
		maplist = append(maplist, strings.Split(m, "\n")[1:])
	}

	lowest := 999_999_999_999
	for _, seed := range seeds {
		soil := getFromChainMap(seed, maplist)
		lowest = min(soil, lowest)

	}

	return lowest
}

func getFromChainMap(value int, maps [][]string) int {
	for _, m := range maps {
		value = getFromMap(value, m)
	}
	return value
}

func getFromMap(value int, ranges []string) int {
	for _, strRange := range ranges {
		nrange := parseNumbers(strRange)
		dst, src, rlen := nrange[0], nrange[1], nrange[2]
		if value >= src && value <= src+rlen {
			value = value + (dst - src)
			break

		}
	}
	return value
}

func toRange(n ...int) map[int]int {
	mapping := map[int]int{}

	dst, src, rlen := n[0], n[1], n[2]
	for i := 0; i < rlen; i++ {
		mapping[src+i] = dst + i

	}
	return mapping
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

func isNumeric(str string) bool {
	for _, c := range str {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
