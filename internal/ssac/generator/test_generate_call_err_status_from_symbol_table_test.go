//ff:func feature=ssac-gen type=test control=sequence
//ff:what @call 대상 함수의 @error 어노테이션에서 ErrStatus를 가져오는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateCallErrStatusFromSymbolTable(t *testing.T) {
	// @call 대상 함수에 @error 401 어노테이션 → .ssac 명시 없으면 401 사용
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
			{Type: parser.SeqCall, Model: "auth.VerifyPassword", Inputs: map[string]string{"Email": "request.Email", "Password": "request.Password"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `http.StatusUnauthorized`)
	assertNotContains(t, code, `http.StatusInternalServerError`)
}
