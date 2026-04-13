//ff:func feature=ssac-gen type=test control=sequence
//ff:what @delete 시퀀스의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateDelete(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CancelReservation", FileName: "cancel_reservation.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &ssacparser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: ssacparser.SeqDelete, Model: "Reservation.Cancel", Inputs: map[string]string{"ID": "reservation.ID"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `err = h.ReservationModel.WithTx(tx).Cancel(reservation.ID)`)
}
