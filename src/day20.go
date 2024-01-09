package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	HighPulse = true
	LowPulse  = false
	StateOn   = true
	StateOff  = false
)

type ConjStates map[string]map[string]bool

const testInputP1 = `
broadcaster -> a, b, c
%a -> b
%b -> c
%c -> inv
&inv -> a`

const testInputP2 = `
broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day20.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 32000000 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 32000000)
	}

	testResultP2 := solveP1(testInputP2)
	if testResultP2 != 11687500 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 11687500)
	}

	resultP1 := solveP1(string(content))
	fmt.Println(resultP1)

	// resultP2 := solveP2(string(content))
	// fmt.Println(resultP2)
}

func solveP1(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	conjStates := parseConjunctionsState(lines)
	modules := make(map[string]Module)
	for _, line := range lines {
		l, m := parseModule(line, conjStates)
		if m != nil {
			modules[l] = m
		}
	}
	l, h := runButtonPressNth(modules, 1000)
	return l * h
}

func solveP2(input string) int {
	return 0
}

func parseConjunctionsState(lines []string) ConjStates {
	connections := make(map[string]string)
	var conjunctions []string

	for _, line := range lines {
		connections[parseLabel(line)] = strings.Join(parseDst(line), ",")
		if strings.HasPrefix(line, "&") {
			conjunctions = append(conjunctions, parseLabel(line))
		}
	}

	conjStates := make(map[string]map[string]bool)
	for _, c := range conjunctions {
		conjStates[c] = make(map[string]bool)
	}

	for k, v := range connections {
		for _, c := range conjunctions {
			if strings.Contains(v, c) {
				conjStates[c][k] = false
			}
		}
	}
	return conjStates
}

func parseModule(line string, states ConjStates) (string, Module) {
	if strings.HasPrefix(line, "broadcaster") {
		return "broadcaster", &BroadCasterModule{"broadcaster", parseDst(line)}
	} else if strings.HasPrefix(line, "%") {
		return parseLabel(line), &FlipFlopModule{parseLabel(line), parseDst(line), StateOff}
	} else if strings.HasPrefix(line, "&") {
		state := states[parseLabel(line)]
		return parseLabel(line), &ConjunctionModule{parseLabel(line), parseDst(line), state}
	} else {
		return "", nil
	}
}

func parseLabel(line string) string {
	labelStr := strings.TrimSpace(line[:strings.Index(line, "->")])
	return strings.Trim(labelStr, "%&")
}

func parseDst(line string) []string {
	dstStr := strings.TrimSpace(line[strings.Index(line, "->")+2:])
	return strings.Split(dstStr, ", ")
}

type Message struct {
	src   string
	dst   string
	pulse bool
}

func (m *Message) String() string {
	pulse := "low"
	if m.pulse {
		pulse = "high"
	}
	return fmt.Sprintf("%s -%s-> %s", m.src, pulse, m.dst)
}

type Module interface {
	send(*Message) []*Message
}

type FlipFlopModule struct {
	label   string
	outputs []string
	state   bool
}

func (ff *FlipFlopModule) send(m *Message) []*Message {
	var messages []*Message
	if m.pulse != HighPulse {
		ff.state = !ff.state
		for _, output := range ff.outputs {
			m := &Message{ff.label, output, ff.state}
			messages = append(messages, m)
		}
	}
	return messages
}

type ConjunctionModule struct {
	label   string
	outputs []string
	state   map[string]bool
}

func (cj *ConjunctionModule) send(m *Message) []*Message {

	cj.state[m.src] = m.pulse
	pulse := allHigh(cj.state)

	var messages []*Message
	for _, output := range cj.outputs {
		m := &Message{cj.label, output, !pulse}
		messages = append(messages, m)
	}
	return messages
}

func allHigh(m map[string]bool) bool {
	for _, v := range m {
		if !v {
			return false
		}
	}
	return true
}

type BroadCasterModule struct {
	label   string
	outputs []string
}

func (bc *BroadCasterModule) send(m *Message) []*Message {
	var messages []*Message
	for _, output := range bc.outputs {
		m := &Message{bc.label, output, LowPulse}
		messages = append(messages, m)
	}
	return messages
}

func runIter(modules map[string]Module, messages []*Message) []*Message {
	var newMessages []*Message
	for _, m := range messages {
		module, ok := modules[m.dst]
		if !ok {
			continue
		}
		msgs := module.send(m)
		newMessages = append(newMessages, msgs...)
	}
	return newMessages
}

func runButtonPressNth(modules map[string]Module, n int) (int, int) {
	var totalLow, totalHigh int
	for i := 0; i < n; i++ {
		msg := []*Message{&Message{src: "button", dst: "broadcaster"}}
		for len(msg) > 0 {
			msg = runIter(modules, msg)
			low, high := sumPulses(msg)
			totalLow += low
			totalHigh += high
		}
	}
	return totalLow + n, totalHigh
}

func sumPulses(msg []*Message) (int, int) {
	var low, high int
	for _, m := range msg {
		if m.pulse {
			high++
		} else {
			low++
		}
	}
	return low, high
}
