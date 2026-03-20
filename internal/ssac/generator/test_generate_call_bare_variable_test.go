//ff:func feature=ssac-gen type=test control=sequence
//ff:what @call 결과를 bare variable로 후속 @post에 전달하는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateCallBareVariable(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "Register", FileName: "register.go",
		Imports: []string{"myapp/auth"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqCall, Model: "auth.HashPassword", Inputs: map[string]string{"Password": "request.Password"}, Result: &parser.Result{Type: "string", Var: "hashedPassword"}},
			{Type: parser.SeqPost, Model: "User.Create", Inputs: map[string]string{"Email": "request.Email", "HashedPassword": "hashedPassword", "Role": "request.Role"}, Result: &parser.Result{Type: "User", Var: "user"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"user": "user"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// @call named field
	assertContains(t, code, `auth.HashPassword(auth.HashPasswordRequest{Password: password})`)
	// bare variable: no trailing dot
	assertContains(t, code, `h.UserModel.WithTx(tx).Create(email, hashedPassword, role)`)
	assertNotContains(t, code, `hashedPassword.`)
}
