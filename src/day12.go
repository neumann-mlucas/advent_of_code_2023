package main

import (
	"fmt"
	"os"
	"strings"
    "strconv"
)

const testInputP1 = `
???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1
`

func main() {
	// var fileName string
	// if len(os.Args) >= 2 {
	// 	fileName = os.Args[1]
	// } else {
	// 	fileName = "dat/day12.txt"
	// }
	//
	// content, err := os.ReadFile(fileName)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "error: %s", err)
	// }

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 21 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 21)
	}

	// testResultP2 := solveP2(testInputP1)
	// if testResultP2 != 64 {
	// 	fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 64)
	// }

	// resultP1 := solveP1(string(content))
	// fmt.Println(resultP1)
	//
	// resultP2 := solveP1(string(content), 26501365)
	// fmt.Println(resultP2)
}

func solveP1(input string) int {
    total := 0
    for _, line := range strings.Split(strings.TrimSpace(input), "\n"){
        s, c := NewSprings(line), &Counter{}
        c.resolve(s, -1)
        total += c.count
    }
    return total
}

type Status int

const (
    Operational Status = iota
    Dameged 
    Unknown
)

type Springs struct {
    states []Status
    constraints []int
}

type SpringsPos struct {
    idx, size int
}

func NewSprings(line string) Springs {
    input := strings.Split(line, " ")

    status := []Status{}
    for _, c := range input[0] {
        switch c {
        case '#': status = append(status, Operational)
        case '.': status = append(status, Dameged)
        case '?': status = append(status, Unknown)
    }
    }

    constraints := []int{}
    for _, c := range strings.Split(strings.TrimSpace(input[1]),",")  {
        num, _ := strconv.Atoi(string(c))
        constraints = append(constraints, num)
    }

    return Springs{status, constraints}
}

func (s Springs) String() string {
    out := ""
    for _, c := range s.states {
        switch c {
        case Operational: out += " S "
        case Dameged: out += " x "
        case Unknown: out += " ? "
        }
    }
    return fmt.Sprintf("<Springs: [%s] | %v>", out, s.constraints)
}

func (sp Springs) getGroups(state Status) []SpringsPos {
    groups := []SpringsPos{}

    curr, gidx := 0, 0
    for idx, st := range sp.states {
        if st == state {
            if curr == 0 {
                gidx = idx
            }
            curr += 1
        } else if curr != 0 {
            groups = append(groups, SpringsPos{gidx, curr})
            curr = 0
        }
    }
    if curr != 0 {
        groups = append(groups, SpringsPos{len(sp.states)-curr, curr})
    }
    return groups
}

func (sp Springs) satisfies() bool {
    if sp.numSatisfied() == len(sp.constraints) {
        return true
    }
    return false
}

func (sp Springs) numSatisfied() int {
    opSprings := sp.getGroups(Operational)
    total, lc := 0, min(len(opSprings), len(sp.constraints))
    for idx := 0; idx < lc; idx++ {
        if opSprings[idx].size == sp.constraints[idx] {
            total += 1
        } else {
            break
        }
    }
    return total
}

type Counter struct {
    count int
}

func (c *Counter) resolve(sp Springs, start int) {
    if sp.satisfies() {
        c.count += 1
        return
    }

    for idx := max(0, start); idx < len(sp.states); idx++ {
        if sp.states[idx] != Unknown { 
            continue // only Unknowns can change
        }
        // heuristic, only look at states where the score equal or greeter score
        oldScore := sp.numSatisfied()
        sp.states[idx] = Operational
        newScore := sp.numSatisfied()
        if newScore >= oldScore {
            newState := make([]Status, len(sp.states))
            copy(newState, sp.states)
            c.resolve(Springs{newState, sp.constraints}, idx) // branch out
        } 
        sp.states[idx] = Unknown // restore previous state
    }
}


