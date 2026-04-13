//ff:func feature=gen-hurl type=util control=iteration dimension=3
//ff:what Walks stateDiagrams BFS from initial state, returning event -> order index for sorting.
package hurl

import "github.com/park-jun-woo/fullend/pkg/parser/statemachine"

// buildTransitionOrder walks stateDiagrams BFS from initial state,
// returning event -> order index for sorting.
func buildTransitionOrder(diagrams []*statemachine.StateDiagram) map[string]int {
	order := make(map[string]int)
	idx := 0
	for _, d := range diagrams {
		// BFS from initial state following transitions.
		visited := make(map[string]bool)
		queue := []string{d.InitialState}
		visited[d.InitialState] = true
		for len(queue) > 0 {
			state := queue[0]
			queue = queue[1:]
			for _, t := range d.Transitions {
				if t.From != state || visited[t.To] {
					continue
				}
				if _, exists := order[t.Event]; !exists {
					order[t.Event] = idx
					idx++
				}
				visited[t.To] = true
				queue = append(queue, t.To)
			}
		}
	}
	return order
}
