//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what counts total transitions across all state diagrams

package orchestrator

import "github.com/park-jun-woo/fullend/internal/statemachine"

// countTransitions counts total transitions across all state diagrams.
func countTransitions(diagrams []*statemachine.StateDiagram) int {
	total := 0
	for _, d := range diagrams {
		total += len(d.Transitions)
	}
	return total
}
