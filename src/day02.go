package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const testInputP1 = `Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day02.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := CheckGames(testInputP1)
	if testResultP1 != 8 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d", testResultP1, 8)
	}

	testResultP2 := CalcPowerSet(testInputP1)
	if testResultP2 != 2286 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d", testResultP2, 2286)
	}

	inputP1 := strings.TrimSpace(string(content))
	resultP1 := CheckGames(inputP1)
	fmt.Println(resultP1)

	inputP2 := strings.TrimSpace(string(content))
	resultP2 := CalcPowerSet(inputP2)
	fmt.Println(resultP2)
}

func CheckGames(games string) int {
	limit := &Play{red: 12, green: 13, blue: 14}
	var total int

	for _, g := range strings.Split(games, "\n") {
		game := parseGame(g)
		if game.IsValid(limit) {
			total += game.id
		}

	}
	return total
}

func CalcPowerSet(games string) int {
	var total int

	for _, g := range strings.Split(games, "\n") {
		game := parseGame(g)
		minPlay := game.findMin()
		total += minPlay.red * minPlay.green * minPlay.blue
	}
	return total
}

type Game struct {
	id    int
	plays []*Play
}

func (g *Game) String() string {
	var strPlay string
	for _, p := range g.plays {
		strPlay += "\n\t" + p.String()

	}
	return fmt.Sprintf("Game %03d: %s", g.id, strPlay)
}

func (g *Game) IsValid(limit *Play) bool {
	for _, p := range g.plays {
		if !p.IsValid(limit) {
			return false
		}
	}
	return true
}

func (g *Game) findMin() *Play {
	var mRed, mGreen, mBlue int
	for _, p := range g.plays {
		mRed = max(mRed, p.red)
		mGreen = max(mGreen, p.green)
		mBlue = max(mBlue, p.blue)
	}
	return &Play{red: mRed, green: mGreen, blue: mBlue}
}

type Play struct {
	red   int
	blue  int
	green int
}

func (p *Play) String() string {
	return fmt.Sprintf("Play: red: %02d, blue: %02d, green: %02d", p.red, p.blue, p.green)
}

func (p *Play) IsValid(limit *Play) bool {
	return p.red <= limit.red && p.blue <= limit.blue && p.green <= limit.green
}

func parseGame(inp string) *Game {
	parts := strings.SplitN(inp, ":", 2)
	id, _ := strconv.Atoi(strings.TrimSpace(strings.Replace(parts[0], "Game", "", 1)))
	var plays []*Play
	for _, p := range strings.Split(parts[1], ";") {
		play := parsePlay(p)
		plays = append(plays, play)
	}
	return &Game{id: id, plays: plays}
}

func parsePlay(inp string) *Play {
	var nRed, nBlue, nGreen int

	reRed := regexp.MustCompile(`(\d+) red`)
	reGreen := regexp.MustCompile(`(\d+) green`)
	reBlue := regexp.MustCompile(`(\d+) blue`)

	mRed := reRed.FindStringSubmatch(inp)
	if mRed != nil && len(mRed) > 1 {
		nRed, _ = strconv.Atoi(mRed[1])
	}
	mBlue := reBlue.FindStringSubmatch(inp)
	if mBlue != nil && len(mBlue) > 1 {
		nBlue, _ = strconv.Atoi(mBlue[1])
	}
	mGreen := reGreen.FindStringSubmatch(inp)
	if mGreen != nil && len(mGreen) > 1 {
		nGreen, _ = strconv.Atoi(mGreen[1])
	}
	return &Play{red: nRed, green: nGreen, blue: nBlue}
}
