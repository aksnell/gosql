// Compiler produces the byte code control flow graph able to describe the
// execution of a tokenizing FSM from a passed Regex string.

package main

import (
	"fmt"
	"regexp"
)

// A fragment references a Graph section being connected by a compiler.
type fragment struct {
	start *State
	out   *edgeList
}

// makeFragment returns an initilized Fragment.
// if l1 is nil, a new edgeList refernecing s.edge is set for its 'out' field.
func makeFragment(s *State, l1 *edgeList) fragment {
	if l1 == nil {
		l1 = makeEdgeList(s.edge)
	}
	return fragment{
		start: s,
		out:   l1,
	}
}

// An edgeList references a chain of dangling State edges in a fragment.
type edgeList struct {
	edge **State
	next *edgeList
}

// makeEdgeList returns an initialized edgeList.
func makeEdgeList(s *State) *edgeList {
	return &edgeList{
		edge: &s,
		next: nil,
	}
}

// patch connects every dangling State edge in an edgeList to the passed State.
func (l1 *edgeList) patch(s *State) {
	for el := l1; el != nil; el = el.next {
		*el.edge = s
	}
}

// Compile takes a regex string and returns a connected State Graph.
func Compile(re string) *Graph {
	c := makeCompiler(re)
	return c.compile()
}

// A compiler holds the state required to connect the nodes in State Graph.
type compiler struct {
	ptr   int
	graph *Graph
	stack []fragment
	re    string
}

// makeCompiler returns an initilized compiler.
func makeCompiler(re string) *compiler {
	return &Compiler{
		ptr:   0,
		graph: makeGraph(),
		stack: make([]Fragment, len(re)),
		re:    re,
	}
}

func (c *compiler) compile() *Graph {
	for _, r := range c.re {
		switch r {
		case '*':
			e1 := c.pop()
			s := makeState(ISplit, e1.start, nil, r)
			e1.out.patch(s)
			c.push(makeFragment(s, makeEdgeList(s.alt)))
			break
		default:
			s := makeState(IRune, nil, nil, r)
			c.push(makeFragment(s, nil))
			break
		}
	}
	c.cat()
	f := c.pop()
	f.out.patch(makeState(IMatch, nil, nil, 0))
	return c.graph
}

func (c *compiler) push(f fragment) {
	c.stack[c.ptr] = f
	c.ptr++
}

func (c *compiler) pop() fragment {
	c.ptr--
	return c.stack[c.ptr]
}

func (c *compiler) cat() {
	var e2, e1 fragment
	for c.ptr != 1 {
		e2, e1 = c.pop(), c.pop()
		e1.out.patch(e2.start)
		c.push(fragment{e1.start, e2.out})
	}
}

func main() {
	t, _ := regexp.Compile("aab+")
	fmt.Println(t.Match([]byte{'a', 'a', 'b', 'b'}))
}
