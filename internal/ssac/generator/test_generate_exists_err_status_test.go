//ff:func feature=ssac-gen type=test control=sequence
//ff:what @exists ErrStatus 커스텀 상태 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateExistsErrStatus(t *testing.T) {
	sf := parser.ServiceFunc{
		Name:     "Register",
		FileName: "register.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByEmail", Inputs: map[string]string{"Email": "request.Email"}, Result: &parser.Result{Type: "*User", Var: "user"}},
			{Type: parser.SeqExists, Target: "user", Message: "Already registered", ErrStatus: 422},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "http.StatusUnprocessableEntity")
	assertNotContains(t, code, "http.StatusConflict")
}
