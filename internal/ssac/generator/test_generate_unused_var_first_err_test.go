//ff:func feature=ssac-gen type=test control=sequence
//ff:what 첫 시퀀스에서 미사용 변수 시 _, err := 패턴을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateUnusedVarFirstErr(t *testing.T) {
	// 첫 시퀀스에서 Unused → _, err := (err 첫 선언)
	sf := parser.ServiceFunc{
		Name: "DoSomething", FileName: "do_something.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Token.Generate", Inputs: map[string]string{"Key": "request.Key"}, Result: &parser.Result{Type: "Token", Var: "token"}},
			{Type: parser.SeqResponse, Fields: map[string]string{}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// token은 미참조 + err 첫 선언 → _, err :=
	assertContains(t, code, `_, err := h.TokenModel.Generate`)
}
