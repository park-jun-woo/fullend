//ff:func feature=ssac-gen type=test control=sequence
//ff:what 전체 시퀀스 조합의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateFullExample(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CancelReservation", FileName: "cancel_reservation.go",
		Imports: []string{"myapp/billing"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqAuth, Action: "cancel", Resource: "reservation", Inputs: map[string]string{"id": "request.ReservationID"}, Message: "권한 없음"},
			{Type: ssacparser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ReservationID": "request.ReservationID"}, Result: &ssacparser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: ssacparser.SeqEmpty, Target: "reservation", Message: "예약을 찾을 수 없습니다"},
			{Type: ssacparser.SeqState, DiagramID: "reservation", Inputs: map[string]string{"status": "reservation.Status"}, Transition: "cancel", Message: "취소할 수 없습니다"},
			{Type: ssacparser.SeqCall, Model: "billing.CalculateRefund", Inputs: map[string]string{"ID": "reservation.ID", "StartAt": "reservation.StartAt", "EndAt": "reservation.EndAt"}, Result: &ssacparser.Result{Type: "Refund", Var: "refund"}},
			{Type: ssacparser.SeqPut, Model: "Reservation.UpdateStatus", Inputs: map[string]string{"ReservationID": "request.ReservationID", "Status": `"cancelled"`}},
			{Type: ssacparser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ReservationID": "request.ReservationID"}, Result: &ssacparser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"reservation": "reservation", "refund": "refund"}},
		},
	}
	code := mustGenerate(t, sf, nil)

	// auth
	assertContains(t, code, `authz.Check(authz.CheckRequest{`)
	// get
	assertContains(t, code, `reservation, err := h.ReservationModel.WithTx(tx).FindByID`)
	// empty
	assertContains(t, code, `if reservation == nil`)
	// state
	assertContains(t, code, `reservationstate.CanTransition`)
	// call
	assertContains(t, code, `billing.CalculateRefund`)
	// put
	assertContains(t, code, `h.ReservationModel.WithTx(tx).UpdateStatus`)
	// response
	assertContains(t, code, `"reservation":`)
	assertContains(t, code, `"refund":`)
}
