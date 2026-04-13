//ff:func feature=ssac-gen type=test control=sequence
//ff:what 미사용 변수 + err already declared 시 _, err = 패턴을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateUnusedVarErrAlreadyDeclared(t *testing.T) {
	// 2번째 @get에서 Unused + err already declared → _, err =
	sf := ssacparser.ServiceFunc{
		Name: "DoSomething", FileName: "do_something.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"ID": "request.UserID"}, Result: &ssacparser.Result{Type: "User", Var: "user"}},
			{Type: ssacparser.SeqGet, Model: "Token.Generate", Inputs: map[string]string{"UserID": "user.ID"}, Result: &ssacparser.Result{Type: "Token", Var: "token"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"user": "user"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// token은 미참조 + err already declared → _, err =
	assertContains(t, code, `_, err = h.TokenModel.Generate(user.ID)`)
	// user는 참조됨 → user, err :=
	assertContains(t, code, `user, err := h.UserModel.FindByID`)
}
