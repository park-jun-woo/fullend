//ff:func feature=ssac-gen type=test control=sequence
//ff:what 리터럴 값 인자의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateLiteralArg(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "Cancel", FileName: "cancel.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPut, Model: "Reservation.UpdateStatus", Inputs: map[string]string{"ID": "request.ID", "Status": `"cancelled"`}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `h.ReservationModel.WithTx(tx).UpdateStatus(id, "cancelled")`)
}
