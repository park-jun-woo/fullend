//ff:type feature=gen-hurl type=model topic=scenario-order
//ff:what StepClass — DecideMidStepClass 판정 결과 열거

package hurl

// StepClass identifies how a Mid-phase step should be ordered.
type StepClass int

const (
	// ClassExcluded: state op on a branching event not selected.
	ClassExcluded StepClass = iota
	// ClassStateTransition: @state sequence op, order = transitionOrder.
	ClassStateTransition
	// ClassUpdate: PUT/PATCH without @state, order = 900.
	ClassUpdate
	// ClassTopLevelCreate: POST with no parent resource, order = -1.
	ClassTopLevelCreate
	// ClassNestedUnderTransition: POST nested under a state-bearing parent, order = firstTransition + 0.5.
	ClassNestedUnderTransition
	// ClassNestedOrphan: POST nested but parent has no first transition, order = -0.5.
	ClassNestedOrphan
)
