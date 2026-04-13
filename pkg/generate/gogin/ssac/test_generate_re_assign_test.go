//ff:func feature=ssac-gen type=test control=sequence
//ff:what 동일 변수 재할당 시 := → = 전환을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateReAssign(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CancelReservation", FileName: "cancel_reservation.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &ssacparser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: ssacparser.SeqPut, Model: "Reservation.UpdateStatus", Inputs: map[string]string{"ID": "request.ID", "Status": `"cancelled"`}},
			{Type: ssacparser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &ssacparser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"reservation": "reservation"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// 첫 번째 @get: :=
	assertContains(t, code, `reservation, err := h.ReservationModel.WithTx(tx).FindByID`)
	// 두 번째 @get: = (재선언)
	assertContains(t, code, `reservation, err = h.ReservationModel.WithTx(tx).FindByID`)
}
