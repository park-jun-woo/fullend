//ff:func feature=ssac-gen type=test control=sequence
//ff:what Domain 필드에 따른 패키지명 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateFeaturePackage(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "Login", FileName: "login.go", Feature: "auth",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "User.FindByEmail", Inputs: map[string]string{"Email": "request.Email"}, Result: &ssacparser.Result{Type: "User", Var: "user"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"user": "user"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "package auth")
}
