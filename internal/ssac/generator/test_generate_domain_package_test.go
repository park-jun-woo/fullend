//ff:func feature=ssac-gen type=test control=sequence
//ff:what Domain 필드에 따른 패키지명 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateDomainPackage(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "Login", FileName: "login.go", Domain: "auth",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByEmail", Inputs: map[string]string{"Email": "request.Email"}, Result: &parser.Result{Type: "User", Var: "user"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"user": "user"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "package auth")
}
