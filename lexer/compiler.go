package lexer

// Compile parses a series of instructions from a Regex into Byte Code
// which is embedded into State(s). Additional Regex instructions define
// how the dangling edges of new State(s) are connected to form the final
// shape of the Graph.
func Compile(re string) *Graph {
	c := compiler{
		ptr:   0,
		stack: make([]fragment, len(re)),
		re:    re,
	}
	return c.compile()
}

// A fragment represents a logical section of unconnected State edges
// visible to the Compiler which can still be manipulated by future
// instructions.
type fragment struct {
	start *State
	out   *edgeList
}

// An edgeList references a chain of dangling State edges starting from
// a fragment.
type edgeList struct {
	edge **State
	next *edgeList
}

// The calling edgeList has all of its dangling State pointers patch(ed) to
// the past State*.
func (l1 *edgeList) patch(s *State) {
	for el := l1; el != nil; el = el.next {
		*el.edge = s
	}
}

// Calling edgeList has the passed edgeList connect(ed) to the end
// of its list.
func (l1 *edgeList) connect(l2 *edgeList) {
	l1End := l1
	for l1End.next != nil {
		l1End = l1End.next
	}
	l1End.next = l2
}

// A compiler holds the state required to manipulate fragments and connect
// State edges into a Graph.
type compiler struct {
	ptr   int
	stack []fragment
	re    string
	graph *Graph
}

func (c *compiler) init(re string) {
	c.ptr = 0
	c.stack = make([]fragment, len(re))
	c.re = re
	c.graph = &Graph{}
	c.graph.init()
}

// Loops over each instruction in the compiler's Regex string,
// compiling the instructions into States and their constituent fragments
// in order to connect State edges and form a complete Graph.
func (c *compiler) compile() *Graph {
	for _, r := range c.re {
		switch r {
		case '?': // Zero or One
			e1 := c.pop()
			s := &State{
				Guard: OpSplit,
				Edge:  e1.start,
				Alt:   nil,
				Rune:  nil,
			}
			e1.out.connect(&edgeList{edge: &s.Alt})
			f := fragment{
				start: s,
				out:   e1.out,
			}
			c.push(f)
		case '*': // Zero or Many
			e1 := c.pop()
			s := &State{
				Guard: OpSplit,
				Edge:  e1.start,
				Alt:   nil,
				Rune:  nil,
			}
			e1.out.patch(s)
			f := fragment{
				start: s,
				out:   &edgeList{edge: &s.Alt},
			}
			c.push(f)
			break
		case '+': // One or Many
			e1 := c.pop()
			s := &State{
				Guard: OpSplit,
				Edge:  e1.start,
				Alt:   nil,
				Rune:  nil,
			}
			e1.out.patch(s)
			f := fragment{
				start: s,
				out:   &edgeList{edge: &s.Alt},
			}
			c.push(f)
		default: // Literal
			s := &State{
				Guard: OpRune,
				Edge:  nil,
				Alt:   nil,
				Rune:  []rune{r},
			}
			f := fragment{
				start: s,
				out:   nil,
			}
			c.push(f)
			break
		}
	}
	c.cat()
	f := c.pop()
	f.out.patch(&State{Guard: OpMatch, Edge: nil, Alt: nil, Rune: nil})
	return c.graph
}

// Pushes the fragment into the stack at the current stack pointer
// and then increments the stack pointer.
// Overwrites fragment at stack pointer when pushed!
func (c *compiler) push(f fragment) {
	c.stack[c.ptr] = f
	c.ptr++
}

// First decrements the stack counter, and then returns a copy of the,
// fragment located at the new stack pointer location.
// Does not remove fragment at stack pointer from the stack!
func (c *compiler) pop() fragment {
	c.ptr--
	return c.stack[c.ptr]
}

// Linearly patches together all framents in the stack in reverse order.
func (c *compiler) cat() {
	var e2, e1 fragment
	for c.ptr != 1 {
		e2, e1 = c.pop(), c.pop()
		e1.out.patch(e2.start)
		c.push(fragment{e1.start, e2.out})
	}
}
