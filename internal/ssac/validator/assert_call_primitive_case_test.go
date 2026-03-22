//ff:func feature=ssac-validate type=test control=sequence
//ff:what 단일 기본 타입에 대해 @call 결과 타입 검증을 수행하는 헬퍼

package validator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// assertCallPrimitiveCase validates that a @call with a primitive return type produces an error.
func assertCallPrimitiveCase(t *testing.T, typ string) {
	t.Helper()
	funcs := []parser.ServiceFunc{{
		Name: "Login", FileName: "login.go",
		Sequences: []parser.Sequence{{
			Type:   parser.SeqCall,
			Model:  "auth.IssueToken",
			Inputs: map[string]string{"UserID": "user.ID"},
			Result: &parser.Result{Type: typ, Var: "token"},
		}},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "기본 타입")
	assertHasError(t, errs, "Response struct 타입을 지정하세요")
}
