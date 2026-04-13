//ff:func feature=ssac-gen type=test control=sequence
//ff:what @call 결과를 bare variable로 후속 @post에 전달하는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateCallBareVariable(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "Register", FileName: "register.go",
		Imports: []string{"myapp/auth"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqCall, Model: "auth.HashPassword", Inputs: map[string]string{"Password": "request.Password"}, Result: &ssacparser.Result{Type: "string", Var: "hashedPassword"}},
			{Type: ssacparser.SeqPost, Model: "User.Create", Inputs: map[string]string{"Email": "request.Email", "HashedPassword": "hashedPassword", "Role": "request.Role"}, Result: &ssacparser.Result{Type: "User", Var: "user"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"user": "user"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// @call named field
	assertContains(t, code, `auth.HashPassword(auth.HashPasswordRequest{Password: password})`)
	// bare variable: no trailing dot
	assertContains(t, code, `h.UserModel.WithTx(tx).Create(email, hashedPassword, role)`)
	assertNotContains(t, code, `hashedPassword.`)
}
