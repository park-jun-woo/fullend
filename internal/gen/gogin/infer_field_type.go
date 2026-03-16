//ff:func feature=gen-gogin type=util control=iteration dimension=2
//ff:what determines if the status field is boolean or string-based from state diagram

package gogin

import "github.com/geul-org/fullend/internal/statemachine"

// inferFieldType determines if the status field is boolean or string-based.
// Heuristic: if exactly 2 live states (excluding terminal states like "deleted")
// and one has "un" prefix of the other, it's boolean.
func inferFieldType(d *statemachine.StateDiagram) string {
	// Collect non-terminal states (states that have outgoing transitions).
	hasOutgoing := make(map[string]bool)
	for _, t := range d.Transitions {
		hasOutgoing[t.From] = true
	}

	// If we can find a pair like published/unpublished, it's boolean.
	for _, s := range d.States {
		negated := "un" + s
		for _, s2 := range d.States {
			if s2 == negated {
				return "bool"
			}
		}
	}

	return "string"
}
