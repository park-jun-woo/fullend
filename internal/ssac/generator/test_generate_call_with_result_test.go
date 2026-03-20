//ff:func feature=ssac-gen type=test control=sequence
//ff:what @call 시퀀스(결과 있음)의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateCallWithResult(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "Login", FileName: "login.go",
		Imports: []string{"myapp/auth"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByEmail", Inputs: map[string]string{"Email": "request.Email"}, Result: &parser.Result{Type: "User", Var: "user"}},
			{Type: parser.SeqCall, Model: "auth.VerifyPassword", Inputs: map[string]string{"Email": "user.Email", "Password": "request.Password"}, Result: &parser.Result{Type: "Token", Var: "token"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"token": "token"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `auth.VerifyPassword(auth.VerifyPasswordRequest{`)
	assertContains(t, code, `Email: user.Email`)
	assertContains(t, code, `Password: password`)
	assertContains(t, code, `http.StatusInternalServerError`)
}
