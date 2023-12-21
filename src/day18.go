package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

const testInputP1 = `
R 6 (#70c710)
D 5 (#0dc571)
L 2 (#5713f0)
D 2 (#d2c081)
R 2 (#59c680)
D 2 (#411b91)
L 5 (#8ceee2)
U 2 (#caa173)
L 1 (#1b58a2)
U 2 (#caa171)
R 2 (#7807d2)
U 3 (#a77fa3)
L 2 (#015232)
U 2 (#7a21e3)`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day18.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 62 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 62)
	}

	testResultP2 := solveP2(testInputP1)
	if testResultP2 != 51 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 51)
	}

	resultP1 := solveP1(string(content))
	fmt.Println(resultP1)

	resultP2 := solveP2(string(content))
	fmt.Println(resultP2)
}

func solveP1(input string) int {
	input = strings.TrimSpace(input)
	instructions := []*Instruction{}

	for _, i := range strings.Split(input, "\n") {
		instructions = append(instructions, parseInstructionP1(i))
	}
	start := &Coord{0, 0}
	coords := processInstructions(start, instructions)
	return calcArea(coords)
}

func solveP2(input string) int {
	input = strings.TrimSpace(input)
	instructions := []*Instruction{}

	for _, i := range strings.Split(input, "\n") {
		instructions = append(instructions, parseInstructionP2(i))
		fmt.Println(parseInstructionP2(i))
	}
	// start := &Coord{0, 0}
	// coords := processInstructions(start, instructions)
	return 0
}

type Instruction struct {
	dir    string
	meters int
	color  string
}

type Coord struct {
	i, j int
}

func parseInstructionP1(input string) *Instruction {
	input = strings.TrimSpace(input)
	fields := strings.Fields(input)
	m, _ := strconv.Atoi(fields[1])
	return &Instruction{fields[0], m, fields[2]}
}

func parseInstructionP2(input string) *Instruction {
	input = strings.TrimSpace(input)

	color := strings.Fields(input)[2]
	color = strings.Trim(color, "#()")

	ndir, _ := strconv.Atoi(string(color[5]))
	meters, _ := strconv.ParseInt(color[:5], 16, 64)

	return &Instruction{string("RDLU"[ndir]), int(meters), color}

}

func processInstructions(start *Coord, insts []*Instruction) []*Coord {
	coords := []*Coord{start}
	dirFn := func(c *Coord) *Coord { return &Coord{c.i, c.j} }

	for _, inst := range insts {
		switch inst.dir {
		case "D":
			dirFn = func(c *Coord) *Coord { return &Coord{c.i + 1, c.j} }
		case "U":
			dirFn = func(c *Coord) *Coord { return &Coord{c.i - 1, c.j} }
		case "R":
			dirFn = func(c *Coord) *Coord { return &Coord{c.i, c.j + 1} }
		case "L":
			dirFn = func(c *Coord) *Coord { return &Coord{c.i, c.j - 1} }
		}

		last := coords[len(coords)-1]
		for i := 0; i < inst.meters; i++ {
			last = dirFn(last)
			coords = append(coords, last)
		}
	}
	return coords

}

func calcArea(coords []*Coord) int {
	grid := toGrid(coords)
	ff := NewFloodFill(grid)

	ff.Visit(Coord{0, 0})
	total := len(grid) * len(grid[0])
	exterior := ff.Area()

	return total - exterior
}

func makeGrid(m, n int) [][]bool {
	grid := make([][]bool, m)
	for i := 0; i < m; i++ {
		grid[i] = make([]bool, n)
	}
	return grid
}

func toGrid(coords []*Coord) [][]bool {
	vcoords, hcoords := []int{}, []int{}
	for _, c := range coords {
		vcoords = append(vcoords, c.i)
		hcoords = append(hcoords, c.j)
	}

	m := (slices.Max(vcoords) + 2) - (slices.Min(vcoords) - 1)
	n := (slices.Max(hcoords) + 2) - (slices.Min(hcoords) - 1)
	grid := makeGrid(m, n)

	for _, c := range coords {
		i := c.i + 1 - slices.Min(vcoords)
		j := c.j + 1 - slices.Min(hcoords)
		grid[i][j] = true
	}
	return grid

}

type FloodFill struct {
	grid  [][]bool
	flood [][]bool
}

func NewFloodFill(grid [][]bool) *FloodFill {
	fill := makeGrid(len(grid), len(grid[0]))
	return &FloodFill{grid, fill}
}

func (f *FloodFill) String() string {
	repr := ""
	for _, line := range f.flood {
		for _, val := range line {
			if val {
				repr += "█"
			} else {
				repr += "_"
			}
		}
		repr += "\n"
	}
	return repr
}

func (f *FloodFill) Area() int {
	var total int
	for _, line := range f.flood {
		for _, val := range line {
			if val {
				total += 1
			}
		}
	}
	return total
}

func (f *FloodFill) canFlood(i, j int) bool {
	if i >= 0 && i < len(f.grid) && j >= 0 && j < len(f.grid[0]) {
		if !f.grid[i][j] && !f.flood[i][j] {
			return true
		} else {
			return false
		}

	}
	return false
}

func (f *FloodFill) Visit(curr Coord) {
	f.flood[curr.i][curr.j] = true
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if f.canFlood(curr.i+i, curr.j+j) {
				f.Visit(Coord{curr.i + i, curr.j + j})
			}
		}
	}
}
