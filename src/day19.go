package main

import (
	"fmt"
	"os"
	// "slices"
	"strconv"
	"strings"
)

const testInputP1 = `
px{a<2006:qkq,m>2090:A,rfg}
pv{a>1716:R,A}
lnx{m>1548:A,A}
rfg{s<537:gd,x>2440:R,A}
qs{s>3448:A,lnx}
qkq{x<1416:A,crn}
crn{x>2662:A,R}
in{s<1351:px,qqz}
qqz{s>2770:qs,m<1801:hdj,R}
gd{a>3333:R,R}
hdj{m>838:A,pv}

{x=787,m=2655,a=1222,s=2876}
{x=1679,m=44,a=2067,s=496}
{x=2036,m=264,a=79,s=2244}
{x=2461,m=1339,a=466,s=291}
{x=2127,m=1623,a=2188,s=1013}`

func main() {
	var fileName string
	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		fileName = "dat/day19.txt"
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

	testResultP1 := solveP1(testInputP1)
	if testResultP1 != 19114 {
		fmt.Fprintf(os.Stderr, "Test input P1 fail: %d, expected: %d\n", testResultP1, 19114)
	}

	testResultP2 := solveP2(testInputP1)
	if testResultP2 != 167409079868000 {
		fmt.Fprintf(os.Stderr, "Test input P2 fail: %d, expected: %d\n", testResultP2, 167409079868000)
	}

	resultP1 := solveP1(string(content))
	fmt.Println(resultP1)

	// resultP2 := solveP2(string(content))
	// fmt.Println(resultP2)
}

func solveP1(input string) int {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	workflowsRaw, objsRaw := strings.Fields(parts[0]), strings.Fields(parts[1])

	workflows := make(map[string]*WorkFlow)
	for _, line := range workflowsRaw {
		wf := parseWorkFlow(line)
		workflows[wf.Name] = wf
	}

	var objs []map[string]int
	for _, line := range objsRaw {
		m := stringToMap(line)
		objs = append(objs, m)
	}

	var total int
	for _, obj := range objs {
		label := "in"
		for label != "R" && label != "A" {
			label = workflows[label].Run(obj)
		}
		if label == "A" {
			total += SumRates(obj)
		}
	}
	return total
}

func solveP2(input string) int {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	workflowsRaw, _ := strings.Fields(parts[0]), strings.Fields(parts[1])

	workflows := make(map[string]*WorkFlow)
	for _, line := range workflowsRaw {
		wf := parseWorkFlow(line)
		workflows[wf.Name] = wf
	}

	var paths [][]*Rule
	ww := WorkFlowWalker{workflows, paths}
	ww.Run("in", []*Rule{})
	fmt.Println()

	var total int
	for _, p := range ww.paths {
		total += runPath(p)
		fmt.Println("\n", pathToStr(p))
		fmt.Println(total, "\n")
	}
	return total
}

type Cond func(map[string]int) bool

type Rule struct {
	Cond  Cond
	Label string
	Raw   string
}

type WorkFlow struct {
	Rules []*Rule
	Name  string
}

func parseRule(rule string) *Rule {
	// default case, rule is just a label
	if !strings.Contains(rule, ":") {
		return &Rule{func(s map[string]int) bool { return true }, rule, "END:" + rule}
	}

	parts := strings.SplitN(rule, ":", 2)
	cond, label := parts[0], parts[1]

	// check what function is used in condition
	op := func(a, b int) bool { return a < b }
	opstr := "<"
	if strings.Contains(cond, ">") {
		op = func(a, b int) bool { return a > b }
		opstr = ">"
	}

	// parse condition e.g. a>3333
	condParts := strings.SplitN(strings.Trim(cond, "{}"), opstr, 2)
	key, valueStr := condParts[0], condParts[1]
	value, _ := strconv.Atoi(valueStr)

	fn := func(m map[string]int) bool {
		if val, ok := m[key]; ok {
			return op(val, value)
		}
		return false

	}
	return &Rule{fn, label, strings.Trim(cond, "{}")}
}

func parseWorkFlow(input string) *WorkFlow {
	name := input[:strings.Index(input, "{")]
	rules := input[strings.Index(input, "{"):strings.Index(input, "}")]

	var parsed []*Rule
	for _, rule := range strings.Split(rules, ",") {
		parsed = append(parsed, parseRule(rule))
	}
	return &WorkFlow{parsed, name}
}

func (wf *WorkFlow) Run(m map[string]int) string {
	for _, r := range wf.Rules {
		if r.Cond(m) {
			return r.Label
		}
	}
	return ""
}

func stringToMap(input string) map[string]int {
	input = strings.Trim(input, "{}")

	result := make(map[string]int)
	for _, pair := range strings.Split(input, ",") {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			continue
		}

		val, err := strconv.Atoi(kv[1])
		if err != nil {
			continue
		}

		result[kv[0]] = val
	}
	return result
}

func SumRates(m map[string]int) int {
	var total int
	for _, v := range m {
		total += v
	}
	return total
}

type WorkFlowWalker struct {
	workflows map[string]*WorkFlow
	paths     [][]*Rule
}

func (ww *WorkFlowWalker) Run(curr string, path []*Rule) {
	wf := ww.workflows[curr]
	for _, r := range wf.Rules {

		if r.Label == "A" {
			// fmt.Println(pathToStr(path) + "->A")
			ww.paths = append(ww.paths, append(path, r))
		} else if r.Label == "R" {
			// fmt.Println(pathToStr(path) + "->R")
		} else {
			ww.Run(r.Label, append(path, r))
		}

		r := rev(r)
		path = append(path, r)
	}

}

func rev(r *Rule) *Rule {
	raw := r.Raw
	if strings.Contains(raw, ">") {
		raw = strings.Replace(raw, ">", "<", -1)
	} else if strings.Contains(raw, "<") {
		raw = strings.Replace(raw, "<", ">", -1)
	}
	return &Rule{r.Cond, r.Label, raw}
}

func pathToStr(rules []*Rule) string {
	path := "in"
	for _, r := range rules {
		path += fmt.Sprintf("->%s(%s)", r.Label, r.Raw)
	}
	return path
}

type Interval struct {
	lb, ub int
}

func (i *Interval) String() string {
	return fmt.Sprintf(" [%d , %d] ", i.lb, i.ub)
}

func (i *Interval) less(n int) {
	i.ub = min(n, i.ub)
}

func (i *Interval) grether(n int) {
	i.lb = max(n, i.lb)
}

func (i *Interval) len() int {
	lenInt := i.ub - i.lb - 1
	if lenInt > 0 {
		return lenInt
	}
	return 0
}

func runPath(path []*Rule) int {
	obj := map[string]*Interval{
		"x": &Interval{1, 4000},
		"m": &Interval{1, 4000},
		"a": &Interval{1, 4000},
		"s": &Interval{1, 4000},
	}

	for _, rule := range path {

		var opstr string
		if strings.Contains(rule.Raw, "END") {
			continue
		} else if strings.Contains(rule.Raw, ">") {
			opstr = ">"
		} else if strings.Contains(rule.Raw, "<") {
			opstr = "<"
		}

		condParts := strings.SplitN(rule.Raw, opstr, 2)
		key, valueStr := condParts[0], condParts[1]
		value, _ := strconv.Atoi(valueStr)

		interval := obj[key]
		if opstr == "<" {
			interval.less(value)
		} else if opstr == ">" {
			interval.grether(value)
		}
		obj[key] = interval
	}

	total := 1
	for _, v := range obj {
		total *= v.len()
	}
	fmt.Println(obj)
	return total
}
