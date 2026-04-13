//ff:func feature=gen-hurl type=decider control=sequence topic=scenario-order
//ff:what DecideMidStepClass — StepFacts → MidDecision (depth 2 early-return chain)

package hurl

// DecideMidStepClass selects the StepClass and emits the order key.
// Depth 2: outer stateOp check → nested branchSkip check.
func DecideMidStepClass(f StepFacts) MidDecision {
	if f.IsStateOp {
		if f.IsBranchSkip {
			return MidDecision{Class: ClassExcluded, Include: false}
		}
		return MidDecision{
			Class:   ClassStateTransition,
			Order:   float64(f.TransitionOrder),
			Include: true,
		}
	}
	if f.Step.Method != "POST" {
		return MidDecision{Class: ClassUpdate, Order: 900.0, Include: true}
	}
	if f.ParentResource == "" {
		return MidDecision{Class: ClassTopLevelCreate, Order: -1.0, Include: true}
	}
	if f.HasFirstTransition {
		return MidDecision{
			Class:   ClassNestedUnderTransition,
			Order:   float64(f.FirstTransition) + 0.5,
			Include: true,
		}
	}
	return MidDecision{Class: ClassNestedOrphan, Order: -0.5, Include: true}
}
