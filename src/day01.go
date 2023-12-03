package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const testInputP1 = `1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet`

const testInputP2 = `two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day01.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := sumCalibrationValuesP1(testInputP1)
	if testResultP1 != 142 {
		fmt.Fprintf(os.Stderr, "Test input fail: %d, expected: %d", testResultP1, 142)
	}

	testResultP2 := sumCalibrationValuesP2(testInputP2)
	if testResultP2 != 281 {
		fmt.Fprintf(os.Stderr, "Test input fail: %d, expected: %d", testResultP2, 281)
		panic("asda")
	}

	inputP1 := strings.TrimSpace(string(content))
	resultP1 := sumCalibrationValuesP1(inputP1)
	fmt.Println(resultP1)

	inputP2 := strings.TrimSpace(string(content))
	resultP2 := sumCalibrationValuesP2(inputP2)
	fmt.Println(resultP2)

}

func sumCalibrationValuesP2(calibrations string) int {
	var total int

	for _, s := range strings.Split(calibrations, "\n") {
		cal, err := cleanString(TranslateString(s))
		if err == nil {
			total += cal
		}
	}
	return total
}

func sumCalibrationValuesP1(calibrations string) int {
	var total int

	for _, s := range strings.Split(calibrations, "\n") {
		cal, err := cleanString(s)
		if err == nil {
			total += cal
		}
	}
	return total
}

func TranslateString(s string) string {
	var trans = map[string]string{
		"one":   "o1e",
		"two":   "t2o",
		"three": "t3e",
		"four":  "f4r",
		"five":  "f5e",
		"six":   "s6x",
		"seven": "s7n",
		"eight": "e8t",
		"nine":  "n9e",
	}

	for k, v := range trans {
		s = strings.ReplaceAll(s, k, v)
	}
	return s

}

func cleanString(s string) (int, error) {
	var clean []rune
	for _, c := range s {
		if unicode.IsDigit(c) {
			clean = append(clean, c)
		}
	}
	if len(clean) >= 1 {
		clean = []rune{clean[0], clean[len(clean)-1]}
		return strconv.Atoi(string(clean))
	}
	return 0, errors.New("bad string input")
}
