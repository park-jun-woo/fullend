//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=states
//ff:what checkFuncGuardStates: 전이 이벤트명 대소문자 불일치 시 에러 검증

package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/statemachine"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncGuardStates_TransitionMismatch(t *testing.T) {
	d := &statemachine.StateDiagram{
		ID:     "gig",
		States: []string{"draft", "open"},
		Transitions: []statemachine.Transition{
			{From: "draft", To: "open", Event: "PublishGig"},
		},
	}
	diagramByID := map[string]*statemachine.StateDiagram{"gig": d}

	// Mutation: transition "publishGig" (lowercase p) vs diagram event "PublishGig".
	fn := ssacparser.ServiceFunc{
		Name: "PublishGig",
		Sequences: []ssacparser.Sequence{{
			Type:       "state",
			DiagramID:  "gig",
			Transition: "publishGig",
			Message:    "Cannot transition",
			Inputs:     map[string]string{"status": "gig.Status"},
		}},
	}

	errs := checkFuncGuardStates(fn, diagramByID)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "publishGig") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ERROR for transition case mismatch, got: %+v", errs)
	}
}
