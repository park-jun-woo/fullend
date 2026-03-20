//ff:func feature=ssac-gen type=test control=sequence
//ff:what 미사용 변수 + err already declared 시 _, err = 패턴을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateUnusedVarErrAlreadyDeclared(t *testing.T) {
	// 2번째 @get에서 Unused + err already declared → _, err =
	sf := parser.ServiceFunc{
		Name: "DoSomething", FileName: "do_something.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"ID": "request.UserID"}, Result: &parser.Result{Type: "User", Var: "user"}},
			{Type: parser.SeqGet, Model: "Token.Generate", Inputs: map[string]string{"UserID": "user.ID"}, Result: &parser.Result{Type: "Token", Var: "token"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"user": "user"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// token은 미참조 + err already declared → _, err =
	assertContains(t, code, `_, err = h.TokenModel.Generate(user.ID)`)
	// user는 참조됨 → user, err :=
	assertContains(t, code, `user, err := h.UserModel.FindByID`)
}
