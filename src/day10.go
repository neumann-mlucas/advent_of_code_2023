package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const testInputP1 = `
7-F7-
.FJ|7
SJLL7
|F--J
LJ.LJ
`

const testInputP2 = `
FF7FSF7F7F7F7F7F---7
L|LJ||||||||||||F--J
FL-7LJLJ||||||LJL-77
F--JF--7||LJLJ7F7FJ-
L---JF-JLJ.||-FJLJJ7
|F|F-JF---7F7-L7L|7|
|FFJF7L7F-JF7|JL---7
7-L-JL7||F7|L7F-7F7|
L.L7LFJ|||||FJL7||LJ
L7JLJL-JLJLJL--JLJ.L`

func main() {

	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day10.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 8 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 8)
	}

	testResultP2 := solveP2(testInputP2)
	if testResultP2 != 10 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 10)
	}

	resultP1 := solveP1(string(content))
	fmt.Println(resultP1)

	// resultP2 := solveP2(string(content))
	// fmt.Println(resultP2)
}

func solveP1(input string) int {
	input = strings.TrimSpace(input)

	// initialize maze
	grid := strings.Split(input, "\n")
	dist := make([][]int, len(grid))
	for i := range dist {
		dist[i] = make([]int, len(grid[0]))
	}
	maze := &Maze{grid, dist}

	// fmt.Println(maze)
	for _, p := range findStartingPipes(maze) {
		maze.Visit(p, maze.Start(), 0)

	}
	// fmt.Println(maze.StringDist())
	return maze.MaxDist()
}

func solveP2(input string) int {
	input = strings.TrimSpace(input)
	// initialize maze
	grid := strings.Split(input, "\n")
	dist := make([][]int, len(grid))
	for i := range dist {
		dist[i] = make([]int, len(grid[0]))
	}
	maze := &Maze{grid, dist}
	fmt.Println(maze)
	return 0
}

type Coord struct {
	i int
	j int
}

type Maze struct {
	grid []string
	dist [][]int
}

func (m *Maze) String() string {
	mapper := func(c rune) rune {
		v, ok := map[rune]rune{
			'-': '━', '|': '┃', 'L': '┗', 'J': '┛', '7': '┓', 'F': '┏', '.': ' ', 'S': '■',
		}[c]
		if ok {
			return v
		} else {
			return c
		}
	}
	repr := ""
	for _, line := range m.grid {
		repr += strings.Map(mapper, line)
		repr += "\n"
	}
	return repr
}

func (m *Maze) StringDist() string {
	repr := ""
	for _, line := range m.dist {
		for _, val := range line {
			repr += fmt.Sprintf("% 5d ", val)
		}
		repr += "\n"
	}
	return repr
}

func (m *Maze) MaxDist() int {
	var d int
	for _, line := range m.dist {
		for _, val := range line {
			d = max(d, val)
		}
	}
	return d
}

func (m *Maze) Start() Coord {
	for i := range m.grid {
		for j := range m.grid[i] {
			if m.grid[i][j] == 'S' {
				return Coord{i, j}
			}
		}
	}
	return Coord{-1, -1}
}

func (m *Maze) Get(i, j int) (rune, error) {
	if i >= 0 && i <= len(m.grid) && j >= 0 && j <= len(m.grid[0]) {
		return rune(m.grid[i][j]), nil
	}
	return '0', errors.New("Index i,j Out of Bounds")
}

func (m *Maze) Visit(curr, prev Coord, count int) (int, error) {
	r, err := m.Get(curr.i, curr.j)
	if err != nil {
		return 0, err
	}
	if r == 'S' {
		return count, nil
	}

	count += 1
	if m.dist[curr.i][curr.j] == 0 {
		m.dist[curr.i][curr.j] = count
	} else {
		m.dist[curr.i][curr.j] = min(m.dist[curr.i][curr.j], count)
	}

	switch r {
	case '-':
		if prev.j < curr.j { // prev is east of curr
			return m.Visit(Coord{curr.i, curr.j + 1}, curr, count)
		} else {
			return m.Visit(Coord{curr.i, curr.j - 1}, curr, count)
		}
	case '|':
		if prev.i < curr.i { // prev is north of curr
			return m.Visit(Coord{curr.i + 1, curr.j}, curr, count)
		} else {
			return m.Visit(Coord{curr.i - 1, curr.j}, curr, count)
		}
	case 'L': // '┗'
		if prev.i < curr.i { // prev is to nort of current
			return m.Visit(Coord{curr.i, curr.j + 1}, curr, count)
		} else {
			return m.Visit(Coord{curr.i - 1, curr.j}, curr, count)
		}
	case 'J': // '┛'
		if prev.i < curr.i { // prev is to nort of current
			return m.Visit(Coord{curr.i, curr.j - 1}, curr, count)
		} else {
			return m.Visit(Coord{curr.i - 1, curr.j}, curr, count)
		}
	case 'F': // '┏'
		if prev.i > curr.i { // prev is to south of current
			return m.Visit(Coord{curr.i, curr.j + 1}, curr, count)
		} else {
			return m.Visit(Coord{curr.i + 1, curr.j}, curr, count)
		}
	case '7': // '┓'
		if prev.i > curr.i { // prev is to south of current
			return m.Visit(Coord{curr.i, curr.j - 1}, curr, count)
		} else {
			return m.Visit(Coord{curr.i + 1, curr.j}, curr, count)
		}
	default:
		return 0, errors.New("Unexpect Char in Grid")
	}
}

func findStartingPipes(maze *Maze) []Coord {
	s := maze.Start()
	validPos := []Coord{}

	r, err := maze.Get(s.i, s.j+1)
	if err == nil && strings.ContainsRune("-J7", r) {
		validPos = append(validPos, Coord{s.i, s.j + 1})
	}

	r, err = maze.Get(s.i, s.j-1)
	if err == nil && strings.ContainsRune("-LF", r) {
		validPos = append(validPos, Coord{s.i, s.j - 1})
	}

	r, err = maze.Get(s.i+1, s.j)
	if err == nil && strings.ContainsRune("|JL", r) {
		validPos = append(validPos, Coord{s.i + 1, s.j})
	}

	r, err = maze.Get(s.i-1, s.j)
	if err == nil && strings.ContainsRune("|7F", r) {
		validPos = append(validPos, Coord{s.i - 1, s.j})
	}

	return validPos
}
