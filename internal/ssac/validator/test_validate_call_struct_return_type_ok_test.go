//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what @call 결과가 struct 타입이면 에러 없음 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateCallStructReturnTypeOK(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Login", FileName: "login.go",
		Sequences: []parser.Sequence{{
			Type:   parser.SeqCall,
			Model:  "auth.IssueToken",
			Inputs: map[string]string{"UserID": "user.ID"},
			Result: &parser.Result{Type: "TokenResponse", Var: "token"},
		}},
	}}
	errs := Validate(funcs)
	for _, e := range errs {
		if contains(e.Message, "기본 타입") {
			t.Errorf("unexpected primitive type error: %s", e.Message)
		}
	}
}
