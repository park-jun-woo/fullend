//ff:func feature=gen-hurl type=util control=selection
//ff:what classifies a non-auth, non-read, non-delete step into ordered mid-step

package hurl

// classifyMidStep classifies a step into an ordered mid-step.
// Returns false if the step should be skipped.
func classifyMidStep(s scenarioStep, stateOps, branchSkip map[string]bool, transitionOrder map[string]int, resourceFirstTransition map[string]int) (orderedStep, bool) {
	if stateOps[s.OperationID] && branchSkip[s.OperationID] {
		return orderedStep{}, false
	}
	if stateOps[s.OperationID] {
		return orderedStep{s, float64(transitionOrder[s.OperationID])}, true
	}
	if s.Method != "POST" {
		return orderedStep{s, 900.0}, true
	}
	parentResource := findParentResource(s.Path)
	switch {
	case parentResource == "":
		return orderedStep{s, -1.0}, true
	default:
		firstOrd, ok := resourceFirstTransition[parentResource]
		if ok {
			return orderedStep{s, float64(firstOrd) + 0.5}, true
		}
		return orderedStep{s, -0.5}, true
	}
}
