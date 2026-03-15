//ff:func feature=gen-hurl type=util
//ff:what Finds branching state transitions and marks all but the first as skippable in smoke tests.
package hurl

import "github.com/geul-org/fullend/internal/statemachine"

// buildBranchSkipSet finds state transitions that share a 'from' state (branching)
// and marks all but the first (by transitionOrder) as skippable in smoke tests.
func buildBranchSkipSet(diagrams []*statemachine.StateDiagram, transitionOrder map[string]int) map[string]bool {
	skip := make(map[string]bool)
	for _, d := range diagrams {
		// Group transitions by 'from' state.
		fromGroups := make(map[string][]string) // from -> []event
		for _, t := range d.Transitions {
			fromGroups[t.From] = append(fromGroups[t.From], t.Event)
		}
		// For groups with >1 transitions, keep only the one with lowest order.
		for _, events := range fromGroups {
			if len(events) <= 1 {
				continue
			}
			// Find best (lowest order).
			bestEvent := ""
			bestOrder := 999999
			for _, e := range events {
				if ord, ok := transitionOrder[e]; ok && ord < bestOrder {
					bestOrder = ord
					bestEvent = e
				}
			}
			// Mark all others as skip.
			for _, e := range events {
				if e != bestEvent {
					skip[e] = true
				}
			}
		}
	}
	return skip
}
