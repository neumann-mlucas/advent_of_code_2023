package main

import (
	"fmt"
	"os"
	"strings"
)

const testInputP1 = `
.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day16.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 46 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 46)
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
	grid := strings.Split(input, "\n")

	mr := newMirrorRoom(grid)
	b := &Bean{0, 0, 'E'}
	mr.moveBean(b)

	return mr.Sum()
}

func solveP2(input string) int {
	input = strings.TrimSpace(input)
	grid := strings.Split(input, "\n")

	var total int
	for i := 0; i < len(grid); i++ {
		// left to right
		mr := newMirrorRoom(grid)
		b := &Bean{i, 0, 'E'}
		mr.moveBean(b)
		total = max(total, mr.Sum())

		// right to left
		mr = newMirrorRoom(grid)
		b = &Bean{i, len(grid[0]) - 1, 'W'}
		mr.moveBean(b)
		total = max(total, mr.Sum())
	}

	for j := 0; j < len(grid[0]); j++ {
		// top to bottom
		mr := newMirrorRoom(grid)
		b := &Bean{0, j, 'S'}
		mr.moveBean(b)
		total = max(total, mr.Sum())

		// bottom to top
		mr = newMirrorRoom(grid)
		b = &Bean{len(grid) - 1, j, 'N'}
		mr.moveBean(b)
		total = max(total, mr.Sum())
	}

	return total
}

type Bean struct {
	i, j int
	dir  rune
}

type MirrorRoom struct {
	grid    []string
	visited [][]int
	beanSet map[Bean]struct{}
}

func newMirrorRoom(grid []string) *MirrorRoom {
	visited := make([][]int, len(grid))
	for i := range grid {
		visited[i] = make([]int, len(grid[i]))
	}
	set := make(map[Bean]struct{}, 1000)
	return &MirrorRoom{grid, visited, set}
}

func (m *MirrorRoom) String() string {
	mStr := ""
	for i, l := range m.visited {

		// Energized Cells
		for j := range l {
			if l[j] == 1 {
				mStr += "X"
			} else {
				mStr += "_"
			}
		}

		mStr += " | "

		// Mirror Room
		for j := range m.grid[i] {
			mStr += string(m.grid[i][j])
		}
		mStr += "\n"
	}
	return mStr
}

func (m *MirrorRoom) moveBeanFoward(b *Bean) *Bean {
	if b.dir == 'N' {
		return &Bean{b.i - 1, b.j, b.dir}
	} else if b.dir == 'S' {
		return &Bean{b.i + 1, b.j, b.dir}
	} else if b.dir == 'W' {
		return &Bean{b.i, b.j - 1, b.dir}
	} else if b.dir == 'E' {
		return &Bean{b.i, b.j + 1, b.dir}
	}
	return b
}

func (m *MirrorRoom) moveBean(b *Bean) {
	if !m.inBounds(b.i, b.j) {
		return
	}
	if _, ok := m.beanSet[*b]; ok {
		return
	}

	m.beanSet[*b] = struct{}{}
	m.visited[b.i][b.j] = 1

	switch m.grid[b.i][b.j] {
	case '.':
		m.moveBean(m.moveBeanFoward(b))
	case '|':
		if b.dir == 'W' || b.dir == 'E' {
			m.moveBean(&Bean{b.i - 1, b.j, 'N'})
			m.moveBean(&Bean{b.i + 1, b.j, 'S'})
		} else {
			m.moveBean(m.moveBeanFoward(b))
		}
	case '-':
		if b.dir == 'S' || b.dir == 'N' {
			m.moveBean(&Bean{b.i, b.j - 1, 'W'})
			m.moveBean(&Bean{b.i, b.j + 1, 'E'})
		} else {
			m.moveBean(m.moveBeanFoward(b))
		}
	case '/':
		if b.dir == 'N' {
			m.moveBean(&Bean{b.i, b.j + 1, 'E'})
		} else if b.dir == 'S' {
			m.moveBean(&Bean{b.i, b.j - 1, 'W'})
		} else if b.dir == 'W' {
			m.moveBean(&Bean{b.i + 1, b.j, 'S'})
		} else if b.dir == 'E' {
			m.moveBean(&Bean{b.i - 1, b.j, 'N'})
		}
	case '\\':
		if b.dir == 'N' {
			m.moveBean(&Bean{b.i, b.j - 1, 'W'})
		} else if b.dir == 'S' {
			m.moveBean(&Bean{b.i, b.j + 1, 'E'})
		} else if b.dir == 'W' {
			m.moveBean(&Bean{b.i - 1, b.j, 'N'})
		} else if b.dir == 'E' {
			m.moveBean(&Bean{b.i + 1, b.j, 'S'})
		}
	default:
		return
	}
}

func (m *MirrorRoom) inBounds(i, j int) bool {
	return i >= 0 && i < len(m.grid) && j >= 0 && j < len(m.grid[0])
}

func (m *MirrorRoom) Sum() int {
	var total int
	for _, line := range m.visited {
		for j := range line {
			total += line[j]
		}
	}
	return total
}
