//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=states
//ff:what checkFuncGuardStates: 전이 이벤트명 일치 시 에러 없음 검증

package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/statemachine"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

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
