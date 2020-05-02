package lexer

// A Graph is a series of regex instructions compiled into States connected
// together by their edges. This container functions as a control-flow graph
// for the compiled byte code instructions to be read by a State Machine.
// A Graph always starts with an Error OpCode and ends with a Match OpCode
type Graph struct {
	root *State
}

// Initialize the root State of a Graph to a Fail State if g.root == nil.
func (g *Graph) init() {
	if g.root != nil {
		g.root = &State{
			Guard: OpFail,
			Edge:  nil,
			Alt:   nil,
			Rune:  nil,
		}
	}
}

// A State represents a regex instruction which is explitly compiled into
// Byte Code. A State's inward edge is guarded by a transition function
// identified by its embedded Byte Code and defined in the State Machine
// executing the Graph.
type State struct {
	Guard OpCode
	Edge  *State
	Alt   *State
	Rune  []rune
}

// OpCode(s) represent Byte Code compiled from Regex instructions which
// are intepreted by a host State Machine as transition functions local to the
// State Machine.
type OpCode uint8

const (
	OpFail    OpCode = iota
	OpRune           // [A]
	OpClass          // [A-z]
	OpAny            // *
	OpCapture        // ()
	OpSplit          // |
	OpSplitMatch
	OpMatch
)
