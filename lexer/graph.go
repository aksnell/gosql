// Graph represents the structure of byte code control flow graph powering
// an FSM able to lex tokens.

package main

// A Graph is the compiled control flow diagram for a State Machine.
type Graph struct {
	root *State
}

// makeGraph returns an initilized Graph.
func makeGraph() *Graph {
	return &Graph{
		root: &State{
			guard: IError,
			edge:  nil,
			alt:   nil,
			value: 0,
		},
	}
}

// State represents a node in a control flow graph whose edge is guarded
// by a transition fuction identified by its contained ByteCode.
type State struct {
	guard ByteCode
	edge  *State
	alt   *State
	value rune
}

// makeState returns an initialized State.
func makeState(g ByteCode, e *State, a *State, v rune) *State {
	return &State{
		guard: g,
		edge:  e,
		alt:   a,
		value: v,
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
