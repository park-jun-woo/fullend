//ff:func feature=ssac-gen type=test control=sequence
//ff:what 첫 시퀀스에서 미사용 변수 시 _, err := 패턴을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateUnusedVarFirstErr(t *testing.T) {
	// 첫 시퀀스에서 Unused → _, err := (err 첫 선언)
	sf := ssacparser.ServiceFunc{
		Name: "DoSomething", FileName: "do_something.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Token.Generate", Inputs: map[string]string{"Key": "request.Key"}, Result: &ssacparser.Result{Type: "Token", Var: "token"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// token은 미참조 + err 첫 선언 → _, err :=
	assertContains(t, code, `_, err := h.TokenModel.Generate`)
}
