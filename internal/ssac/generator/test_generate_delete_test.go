//ff:func feature=ssac-gen type=test control=sequence
//ff:what @delete 시퀀스의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateDelete(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "CancelReservation", FileName: "cancel_reservation.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: parser.SeqDelete, Model: "Reservation.Cancel", Inputs: map[string]string{"ID": "reservation.ID"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `err = h.ReservationModel.WithTx(tx).Cancel(reservation.ID)`)
}
