package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const testInputP1 = `
Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day04.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}
	input := strings.TrimSpace(string(content))

	testInputP1 := strings.TrimSpace(testInputP1)
	testResultP1 := SumGames(testInputP1)
	if testResultP1 != 13 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d", testResultP1, 13)
	}

	testResultP2 := SumScratchCards(testInputP1)
	if testResultP2 != 30 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d", testResultP2, 30)
	}

	resultP1 := SumGames(input)
	fmt.Println(resultP1)

	resultP2 := SumScratchCards(input)
	fmt.Println(resultP2)

}

func SumGames(games string) int {
	games = strings.TrimSpace(games)
	var total int
	for _, g := range strings.Split(games, "\n") {
		game := parseGame(g)
		total += game.Score()
	}
	return total
}

func SumScratchCards(games string) int {
	totalCards := map[int]int{}

	games = strings.TrimSpace(games)
	for _, g := range strings.Split(games, "\n") {
		game := parseGame(g)
		for i := game.id; i <= game.id+game.Matching(); i++ {
			if i == game.id {
				totalCards[i]++
			} else {
				totalCards[i] += totalCards[game.id]
			}
		}
	}

	var total int
	for _, v := range totalCards {
		total += v
	}
	return total
}

func parseGame(inp string) *Game {
	idStr := inp[strings.Index(inp, "Card")+4 : strings.Index(inp, ":")]
	id, _ := strconv.Atoi(strings.TrimSpace(idStr))

	handSlice := inp[strings.Index(inp, ":")+1 : strings.Index(inp, "|")]
	hand := strings.Split(strings.TrimSpace(handSlice), " ")

	tableSlice := inp[strings.Index(inp, "|")+1:]
	table := strings.Split(strings.TrimSpace(tableSlice), " ")

	tableMap := map[string]struct{}{}
	for _, n := range table {
		tableMap[n] = struct{}{}
	}
	delete(tableMap, "")

	return &Game{id: id, hand: hand, table: tableMap}
}

type Game struct {
	id    int
	hand  []string
	table map[string]struct{}
}

func (g *Game) String() string {
	return fmt.Sprintf("Card %03d:\n\thand: %v\n\ttable: %v", g.id, g.hand, g.table)
}

func (g *Game) Matching() int {
	var m int
	for _, c := range g.hand {
		if _, ok := g.table[c]; ok {
			m++
		}
	}
	return m
}

func (g *Game) Score() int {
	m := float64(g.Matching())
	if m > 0 {
		return int(math.Pow(2, m-1))
	} else {
		return 0
	}

}
