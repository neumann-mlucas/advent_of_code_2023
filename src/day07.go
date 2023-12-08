package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const testInputP1 = `
32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day07.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solve(testInputP1, false)
	if testResultP1 != 6440 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 6440)
	}

	testResultP2 := solve(testInputP1, true)
	if testResultP2 != 5905 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 5905)
	}

	resultP1 := solve(string(content), false)
	fmt.Println(resultP1)

	resultP2 := solve(string(content), true)
	fmt.Println(resultP2)
}

func solve(input string, p2 bool) int {
	input = strings.TrimSpace(input)
	var hands []*Hand
	for _, line := range strings.Split(input, "\n") {
		hands = append(hands, parseHand(line))
	}

	var total int
	sort.Slice(hands, func(i, j int) bool { return hands[i].Score(p2) < hands[j].Score(p2) })
	for i, h := range hands {
		total += h.Bid * (i + 1)
	}

	return total
}

func parseHand(line string) *Hand {
	play := strings.Fields(line)
	bid, _ := strconv.Atoi(play[1])
	return &Hand{Cards: play[0], Bid: bid}
}

type Hand struct {
	Cards string
	Bid   int
}

func (h *Hand) Score(joker bool) int {

	var hist map[int]string
	if joker {
		jcards := handleJoker(h.Cards)
		hist = countCarts(jcards)
	} else {
		hist = countCarts(h.Cards)
	}

	var has [6]bool
	for i := 0; i <= 5; i++ {
		_, has[i] = hist[i]
	}

	scoreRank := scoreCardRank(h.Cards, "J23456789TQKA")

	if has[5] {
		// Five of a Kind
		return 70000000000 + scoreRank
	} else if has[4] {
		// Four of a Kind
		return 60000000000 + scoreRank
	} else if has[3] && has[2] {
		// Full House
		return 50000000000 + scoreRank
	} else if has[3] {
		// Three of a kind
		return 40000000000 + scoreRank
	} else if has[2] && len(hist[2]) == 2 {
		// Two Pair
		return 30000000000 + scoreRank
	} else if has[2] {
		// One Pair
		return 20000000000 + scoreRank
	} else {
		return 10000000000 + scoreRank
	}
}

func scoreCardRank(cards, strenghts string) int {
	mul, score := 8, 0
	for _, c := range cards {
		score += int(math.Pow10(mul)) * (strings.IndexRune(strenghts, c) + 1)
		mul -= 2
	}
	return score
}

func countCarts(cards string) map[int]string {
	hist := map[rune]int{}
	for _, c := range cards {
		hist[c]++
	}

	revHist := map[int]string{}
	for k, v := range hist {
		revHist[v] += string(k)

	}
	return revHist
}

func handleJoker(hand string) string {
	if hand == "JJJJJ" {
		return "22222"
	}

	common_one := mostCommon(hand)
	common_two := mostCommon(strings.Replace(hand, common_one, "", -1))

	if common_one != "J" {
		return strings.Replace(hand, "J", common_one, -1)
	} else {
		return strings.Replace(hand, "J", common_two, -1)
	}
}

func mostCommon(str string) string {
	max, maxchr := 0, ""
	for _, c := range str {
		chr := string(c)
		count := strings.Count(str, chr)
		if count > max {
			max = count
			maxchr = chr
		}
	}
	return maxchr
}
