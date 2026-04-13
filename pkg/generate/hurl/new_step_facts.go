//ff:func feature=gen-hurl type=util control=sequence topic=scenario-order
//ff:what newStepFacts — scenarioStep + lookup 맵 → StepFacts 프로젝션

package hurl

// newStepFacts projects a scenarioStep + lookup maps into StepFacts.
func newStepFacts(
	s scenarioStep,
	stateOps, branchSkip map[string]bool,
	transitionOrder map[string]int,
	resourceFirstTransition map[string]int,
) StepFacts {
	parent := findParentResource(s.Path)
	ft, okFT := resourceFirstTransition[parent]
	return StepFacts{
		Step:               s,
		IsStateOp:          stateOps[s.OperationID],
		IsBranchSkip:       branchSkip[s.OperationID],
		TransitionOrder:    transitionOrder[s.OperationID],
		ParentResource:     parent,
		FirstTransition:    ft,
		HasFirstTransition: okFT,
	}
}
