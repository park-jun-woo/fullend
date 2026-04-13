//ff:type feature=gen-hurl type=model topic=scenario-order
//ff:what MidDecision — Mid-phase step 분류 + 순서 키 + 포함 여부

package hurl

// MidDecision carries the Class, the computed order key, and whether the step is emitted.
type MidDecision struct {
	Class   StepClass
	Order   float64
	Include bool
}
