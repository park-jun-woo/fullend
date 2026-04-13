//ff:func feature=ssac-gen type=test control=sequence
//ff:what @exists ErrStatus 커스텀 상태 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateExistsErrStatus(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name:     "Register",
		FileName: "register.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "User.FindByEmail", Inputs: map[string]string{"Email": "request.Email"}, Result: &ssacparser.Result{Type: "*User", Var: "user"}},
			{Type: ssacparser.SeqExists, Target: "user", Message: "Already registered", ErrStatus: 422},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "http.StatusUnprocessableEntity")
	assertNotContains(t, code, "http.StatusConflict")
}
