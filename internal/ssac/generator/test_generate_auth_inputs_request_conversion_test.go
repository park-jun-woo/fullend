//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth inputs에서 request.* → 로컬 변수 변환을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateAuthInputsRequestConversion(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "CreateReservation", FileName: "create_reservation.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqAuth, Action: "create", Resource: "reservation", Inputs: map[string]string{"id": "request.RoomID"}, Message: "권한 없음"},
		},
	}
	code := mustGenerate(t, sf, nil)
	// request.RoomID → roomID
	assertContains(t, code, `ID: roomID`)
	assertNotContains(t, code, `request.RoomID`)
	assertContains(t, code, `authz.CheckRequest{`)
}
