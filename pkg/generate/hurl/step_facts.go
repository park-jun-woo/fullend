//ff:type feature=gen-hurl type=model topic=scenario-order
//ff:what StepFacts — DecideMidStepClass 입력 축 값 (per-step lookup 결과)

package hurl

// StepFacts carries the per-step axis values used by DecideMidStepClass.
// Precomputed by newStepFacts to keep the decider side-effect-free and unit-testable.
type StepFacts struct {
	Step               scenarioStep
	IsStateOp          bool
	IsBranchSkip       bool
	TransitionOrder    int
	ParentResource     string
	FirstTransition    int
	HasFirstTransition bool
}
