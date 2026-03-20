//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what @call 결과가 기본 타입이면 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateCallPrimitiveReturnType(t *testing.T) {
	for _, typ := range []string{"string", "int", "int64", "bool", "float64", "time.Time"} {
		t.Run(typ, func(t *testing.T) {
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
		})
	}
}
