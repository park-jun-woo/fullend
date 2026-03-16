//ff:func feature=gen-gogin type=util control=iteration dimension=2 topic=states
//ff:what returns (trueState, falseState) for boolean state diagrams

package gogin

import "github.com/geul-org/fullend/internal/statemachine"

// inferBoolStates returns (trueState, falseState) for boolean diagrams.
func inferBoolStates(d *statemachine.StateDiagram) (string, string) {
	for _, s := range d.States {
		negated := "un" + s
		for _, s2 := range d.States {
			if s2 == negated {
				return s, s2 // s = true, "un"+s = false
			}
		}
	}
	// Fallback: initial state = false.
	if d.InitialState != "" {
		for _, s := range d.States {
			if s != d.InitialState && s != "deleted" {
				return s, d.InitialState
			}
		}
	}
	return d.States[0], d.States[1]
}
