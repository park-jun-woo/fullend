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

func TestCheckFuncGuardStates_TransitionMatch(t *testing.T) {
	d := &statemachine.StateDiagram{
		ID:     "reservation",
		States: []string{"active", "cancelled"},
		Transitions: []statemachine.Transition{
			{From: "active", To: "cancelled", Event: "cancel"},
		},
	}
	diagramByID := map[string]*statemachine.StateDiagram{"reservation": d}

	// transition="cancel", diagram event="cancel" — should pass.
	fn := ssacparser.ServiceFunc{
		Name: "CancelReservation",
		Sequences: []ssacparser.Sequence{{
			Type:       "state",
			DiagramID:  "reservation",
			Transition: "cancel",
			Message:    "취소할 수 없습니다",
			Inputs:     map[string]string{"status": "reservation.Status"},
		}},
	}

	errs := checkFuncGuardStates(fn, diagramByID)
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "cancel") {
			t.Errorf("transition=cancel should match diagram event=cancel, got: %+v", e)
		}
	}
}
