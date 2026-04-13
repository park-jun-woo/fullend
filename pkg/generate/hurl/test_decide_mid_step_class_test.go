//ff:func feature=gen-hurl type=test control=iteration dimension=1 topic=scenario-order
//ff:what DecideMidStepClass 6 StepClass 전수 테이블 테스트

package hurl

import "testing"

func TestDecideMidStepClass(t *testing.T) {
	cases := []struct {
		name  string
		facts StepFacts
		want  MidDecision
	}{
		{
			"state op with branch skip → excluded",
			StepFacts{IsStateOp: true, IsBranchSkip: true},
			MidDecision{Class: ClassExcluded, Include: false},
		},
		{
			"state transition → order = transitionOrder",
			StepFacts{IsStateOp: true, TransitionOrder: 3},
			MidDecision{Class: ClassStateTransition, Order: 3.0, Include: true},
		},
		{
			"non-POST without @state → update (900)",
			StepFacts{Step: scenarioStep{Method: "PUT"}},
			MidDecision{Class: ClassUpdate, Order: 900.0, Include: true},
		},
		{
			"POST top-level (no parent) → -1",
			StepFacts{Step: scenarioStep{Method: "POST", Path: "/users"}, ParentResource: ""},
			MidDecision{Class: ClassTopLevelCreate, Order: -1.0, Include: true},
		},
		{
			"POST nested under transition → firstTransition + 0.5",
			StepFacts{Step: scenarioStep{Method: "POST", Path: "/users/{id}/posts"}, ParentResource: "users", FirstTransition: 2, HasFirstTransition: true},
			MidDecision{Class: ClassNestedUnderTransition, Order: 2.5, Include: true},
		},
		{
			"POST nested orphan (parent has no first transition) → -0.5",
			StepFacts{Step: scenarioStep{Method: "POST", Path: "/users/{id}/tags"}, ParentResource: "users", HasFirstTransition: false},
			MidDecision{Class: ClassNestedOrphan, Order: -0.5, Include: true},
		},
		{
			"state op takes priority over POST semantics",
			StepFacts{Step: scenarioStep{Method: "POST", Path: "/users"}, IsStateOp: true, TransitionOrder: 1},
			MidDecision{Class: ClassStateTransition, Order: 1.0, Include: true},
		},
	}
	for _, tc := range cases {
		if got := DecideMidStepClass(tc.facts); got != tc.want {
			t.Errorf("%s: want %+v got %+v", tc.name, tc.want, got)
		}
	}
}
