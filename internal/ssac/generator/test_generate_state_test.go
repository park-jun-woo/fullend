//ff:func feature=ssac-gen type=test control=sequence
//ff:what @state 가드의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateState(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "CancelReservation", FileName: "cancel_reservation.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: parser.SeqState, DiagramID: "reservation", Inputs: map[string]string{"status": "reservation.Status"}, Transition: "cancel", Message: "취소할 수 없습니다"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `err := reservationstate.CanTransition(reservationstate.Input{`)
	assertContains(t, code, `Status: reservation.Status`)
	assertContains(t, code, `"cancel"`)
	assertContains(t, code, `err != nil`)
	assertContains(t, code, `err.Error()`)
}
