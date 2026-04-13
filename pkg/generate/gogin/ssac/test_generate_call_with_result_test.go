//ff:func feature=ssac-gen type=test control=sequence
//ff:what @call 시퀀스(결과 있음)의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateCallWithResult(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "Login", FileName: "login.go",
		Imports: []string{"myapp/auth"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "User.FindByEmail", Inputs: map[string]string{"Email": "request.Email"}, Result: &ssacparser.Result{Type: "User", Var: "user"}},
			{Type: ssacparser.SeqCall, Model: "auth.VerifyPassword", Inputs: map[string]string{"Email": "user.Email", "Password": "request.Password"}, Result: &ssacparser.Result{Type: "Token", Var: "token"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"token": "token"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `auth.VerifyPassword(auth.VerifyPasswordRequest{`)
	assertContains(t, code, `Email: user.Email`)
	assertContains(t, code, `Password: password`)
	assertContains(t, code, `http.StatusInternalServerError`)
}
