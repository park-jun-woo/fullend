//ff:func feature=ssac-gen type=test control=sequence
//ff:what .ssac 명시값이 @error 어노테이션보다 우선하는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateCallErrStatusSsacOverridesAnnotation(t *testing.T) {
	// .ssac 파일 명시값(500)이 @error 어노테이션(401)보다 우선
	st := &validator.SymbolTable{
		Models: map[string]validator.ModelSymbol{
			"auth._func": {Methods: map[string]validator.MethodInfo{
				"VerifyPassword": {ErrStatus: 401},
			}},
		},
		Operations: map[string]validator.OperationSymbol{},
		DDLTables:  map[string]validator.DDLTable{},
	}
	sf := parser.ServiceFunc{
		Name: "Login", FileName: "login.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqCall, Model: "auth.VerifyPassword", Inputs: map[string]string{"Email": "request.Email"}, ErrStatus: 500},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `http.StatusInternalServerError`)
	assertNotContains(t, code, `http.StatusUnauthorized`)
}
