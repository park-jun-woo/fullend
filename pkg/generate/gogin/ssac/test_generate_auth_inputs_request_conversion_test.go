//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth inputs에서 request.* → 로컬 변수 변환을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateAuthInputsRequestConversion(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CreateReservation", FileName: "create_reservation.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqAuth, Action: "create", Resource: "reservation", Inputs: map[string]string{"id": "request.RoomID"}, Message: "권한 없음"},
		},
	}
	code := mustGenerate(t, sf, nil)
	// request.RoomID → roomID
	assertContains(t, code, `ID: roomID`)
	assertNotContains(t, code, `request.RoomID`)
	assertContains(t, code, `authz.CheckRequest{`)
}
