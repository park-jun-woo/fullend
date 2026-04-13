//ff:func feature=ssac-gen type=test control=sequence
//ff:what @call ErrStatus 명시 시 해당 HTTP 상태 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateCallErrStatus(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "Login", FileName: "login.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqCall, Model: "auth.VerifyPassword", Inputs: map[string]string{"Email": "request.Email", "Password": "request.Password"}, ErrStatus: 401},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `http.StatusUnauthorized`)
	assertNotContains(t, code, `http.StatusInternalServerError`)
}
