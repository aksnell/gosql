package main

import "fmt"

// A Graph is the compiled control flow diagram for a State Machine.
type Graph struct {
	start    *State
	maxdepth int
}

// A State is a single node in a Graph.
type State struct {
	Guard ByteCode
	edge  *State
	alt   *State
	value []rune
}

// A fragment references a section of the Graph still being connected by a compiler.
type fragment struct {
	start *State
	out   *edgeList
}

// An edgeList references a chain of dangling State edges in a Fragment.
type edgeList struct {
	edge **State
	next *edgeList
}

// patch connects every dangling State edge in an edgeList to the passed State.
func (l1 *edgeList) patch(s *State) {
	for el := l1; el != nil; el = el.next {
		*el.edge = s
	}
}

// ByteCode represents the ID of the transition function guarding the edge,
// of a State in a State graph.
type ByteCode uint8

const (
	IRune ByteCode = iota
	IA
	ISplit
	IError
	IMatch
)

// A compiler takes a regex string and parses into a Graph.
type Compiler struct {
	stack []fragment
}

func (c *Compiler) Compile(re string) *Graph {
	stateID := 0
	for _, r := range re {
		switch r {
		case '*':
			e1 := c.Pop()
			s := State{ISplit, e1.start, nil, r, stateID}
			stateID++
			e1.out.patch(&s)
			c.Push(fragment{&s, &edgeList{&s.Alt, nil}})
			break
		default:
			s := State{IRune, nil, nil, r, stateID}
			stateID++
			c.Push(fragment{&s, &edgeList{&s.Edge, nil}})
			break
		}
	}
	c.Cat()
	f := c.Pop()
	f.out.patch(&State{BC: IMatch})
	return &Graph{f.start}
}

func (c *Compiler) Push(f fragment) {
	c.stack[c.fPtr] = f
	c.fPtr++
}

func (c *Compiler) Pop() fragment {
	c.fPtr--
	return c.stack[c.fPtr]
}

func (c *Compiler) Cat() {
	var e2, e1 fragment
	for c.fPtr != 1 {
		e2, e1 = c.Pop(), c.Pop()
		e1.out.patch(e2.start)
		c.Push(fragment{e1.start, e2.out})
	}
}

func main() {
	test := regexp.Compile("aab+")
	fmt.Println(test.Match("aabb"))
}
