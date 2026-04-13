//ff:func feature=ssac-gen type=test control=sequence
//ff:what @call 시퀀스(결과 없음)의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateCallWithoutResult(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "Cancel", FileName: "cancel.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqCall, Model: "notification.Send", Inputs: map[string]string{"ID": "reservation.ID", "Status": `"cancelled"`}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `notification.Send(notification.SendRequest{`)
	assertContains(t, code, `ID: reservation.ID`)
	assertContains(t, code, `Status: "cancelled"`)
	assertContains(t, code, `_, err :=`)
	assertContains(t, code, `http.StatusInternalServerError`)
}
