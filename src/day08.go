package main

import (
	"fmt"
	// "math"
	"os"
	// "sort"
	// "errors"
	"strings"
)

const testInputP1 = `
LLR

AAA = (BBB, BBB)
BBB = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)
`

const testInputP2 = `
LR

11A = (11B, XXX)
11B = (XXX, 11Z)
11Z = (11B, XXX)
22A = (22B, XXX)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22B, 22B)
XXX = (XXX, XXX)`

func main() {

	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day08.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 6 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 6)
	}

	testResultP2 := solveP2(testInputP2)
	if testResultP2 != 6 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 6)
	}

	resultP1 := solveP1(string(content))
	fmt.Println(resultP1)

	resultP2 := solveP2(string(content))
	fmt.Println(resultP2)
}

func solveP1(input string) int {
	input = strings.TrimSpace(input)

	insts := input[:strings.Index(input, "\n\n")]
	nodes := input[strings.Index(input, "\n\n"):]
	nodes = strings.TrimSpace(nodes)

	nodeMap := parseNodes(nodes)

	currNode, n := "AAA", 0
	for currNode != "ZZZ" {
		inst := insts[n%len(insts)]
		if inst == 'L' {
			currNode = nodeMap[currNode][0]
		} else {
			currNode = nodeMap[currNode][1]
		}
		n++
	}
	return n
}

func solveP2(input string) int {
	input = strings.TrimSpace(input)

	insts := input[:strings.Index(input, "\n\n")]
	nodes := input[strings.Index(input, "\n\n"):]
	nodes = strings.TrimSpace(nodes)

	nodeMap := parseNodes(nodes)

	var currNodes []string
	for n := range nodeMap {
		if strings.HasSuffix(n, "A") {
			currNodes = append(currNodes, n)
		}
	}

	walker := Walker{nodes: nodeMap, insts: insts}
	var counts []int

	for _, n := range currNodes {
		var count int
		for !strings.HasSuffix(n, "Z") {
			n = walker.Step(n, count)
			count++
		}
		counts, count = append(counts, count), 0
	}

	return LCM(counts...)
}

type Walker struct {
	nodes map[string][]string
	insts string
}

func (w *Walker) Step(node string, n int) string {
	inst := w.insts[n%len(w.insts)]
	if inst == 'L' {
		return w.nodes[node][0]
	} else {
		return w.nodes[node][1]
	}
}

func parseNodes(input string) map[string][]string {
	nodes := map[string][]string{}
	for _, line := range strings.Split(input, "\n") {
		node, conn := parseNode(line)
		nodes[node] = conn
	}
	return nodes
}

func parseNode(input string) (string, []string) {
	node := strings.TrimSpace(input[:strings.Index(input, "=")])
	connStr := strings.Trim(input[strings.Index(input, "=")+1:], " ()")

	var conn []string
	for _, n := range strings.Split(connStr, ", ") {
		conn = append(conn, n)

	}
	return node, conn
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(integers ...int) int {
	result := integers[0] * integers[1] / GCD(integers[0], integers[1])
	for i := 2; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}
	return result
}
