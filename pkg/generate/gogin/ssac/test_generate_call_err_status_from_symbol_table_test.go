//ff:func feature=ssac-gen type=test control=sequence
//ff:what @call 대상 함수의 @error 어노테이션에서 ErrStatus를 가져오는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGenerateCallErrStatusFromSymbolTable(t *testing.T) {
	// @call 대상 함수에 @error 401 어노테이션 → .ssac 명시 없으면 401 사용
	st := &rule.Ground{
		Models: map[string]rule.ModelInfo{
			"auth._func": {Methods: map[string]rule.MethodInfo{
				"VerifyPassword": {ErrStatus: 401},
			}},
		},
		Ops: map[string]rule.OperationInfo{},
		Tables: map[string]rule.TableInfo{},
	}
	sf := ssacparser.ServiceFunc{
		Name: "Login", FileName: "login.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqCall, Model: "auth.VerifyPassword", Inputs: map[string]string{"Email": "request.Email", "Password": "request.Password"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `http.StatusUnauthorized`)
	assertNotContains(t, code, `http.StatusInternalServerError`)
}
