//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what @call은 외부 패키지이므로 모델 체크를 스킵하는지 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateCallSkipsModelCheck(t *testing.T) {
	st := &SymbolTable{Models: map[string]ModelSymbol{}, Operations: map[string]OperationSymbol{}}
	funcs := []parser.ServiceFunc{{
		Name: "Login", FileName: "login.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqCall, Model: "auth.VerifyPassword", Inputs: map[string]string{"Password": "request.Password"}, Result: &parser.Result{Type: "Token", Var: "token"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"token": "token"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	for _, e := range errs {
		if !e.IsWarning() && (contains(e.Message, "모델") || contains(e.Message, "메서드")) {
			t.Errorf("unexpected model error for @call: %s", e.Message)
		}
	}
}
