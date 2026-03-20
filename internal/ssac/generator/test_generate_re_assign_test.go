//ff:func feature=ssac-gen type=test control=sequence
//ff:what 동일 변수 재할당 시 := → = 전환을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateReAssign(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "CancelReservation", FileName: "cancel_reservation.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: parser.SeqPut, Model: "Reservation.UpdateStatus", Inputs: map[string]string{"ID": "request.ID", "Status": `"cancelled"`}},
			{Type: parser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"reservation": "reservation"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// 첫 번째 @get: :=
	assertContains(t, code, `reservation, err := h.ReservationModel.WithTx(tx).FindByID`)
	// 두 번째 @get: = (재선언)
	assertContains(t, code, `reservation, err = h.ReservationModel.WithTx(tx).FindByID`)
}
