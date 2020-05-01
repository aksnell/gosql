package main

// A fragment references a section of the Graph still being connected by a compiler.
type fragment struct {
	start *State
	out   *edgeList
}

func makeFragment(s *State, l1 *edgeList) fragment {
	if l1 == nil {
		l1 = makeEdgeList(s.edge)
	}
	return fragment{
		start: s,
		out:   l1,
	}
}

// An edgeList references a chain of dangling State edges in a Fragment.
type edgeList struct {
	edge **State
	next *edgeList
}

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

func Compile(re string) *Graph {
	var c compiler
	c.init()
	c.compile(re)
	return c.g
}

// A compiler takes a regex string and parses into a Graph.
type compiler struct {
	ptr   int
	stack []fragment
	g     *Graph
}

func (c *compiler) init() {
	c.ptr = 0
	c.stack = make([]fragment, 64)
	c.g = makeGraph()
	c.push(makeFragment(c.g.root, nil))
}

func (c *compiler) compile(re string) {
	for _, r := range re {
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
}
