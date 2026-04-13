//ff:func feature=ssac-gen type=test control=sequence
//ff:what .ssac 명시값이 @error 어노테이션보다 우선하는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGenerateCallErrStatusSsacOverridesAnnotation(t *testing.T) {
	// .ssac 파일 명시값(500)이 @error 어노테이션(401)보다 우선
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
			{Type: ssacparser.SeqCall, Model: "auth.VerifyPassword", Inputs: map[string]string{"Email": "request.Email"}, ErrStatus: 500},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `http.StatusInternalServerError`)
	assertNotContains(t, code, `http.StatusUnauthorized`)
}
