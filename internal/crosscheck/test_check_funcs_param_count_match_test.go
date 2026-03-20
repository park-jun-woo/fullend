//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_ParamCountMatch: @call 입력 개수와 func spec Request 필드 수 일치 시 에러 없음 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_ParamCountMatch(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package: "auth",
		Name:    "verifyPassword",
		RequestFields: []funcspec.Field{
			{Name: "PasswordHash", Type: "string"},
			{Name: "Password", Type: "string"},
		},
		HasBody: true,
	}}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Login",
		Sequences: []ssacparser.Sequence{{
			Type:  "call",
			Model: "auth.VerifyPassword",
			Inputs: map[string]string{
				"PasswordHash": "user.PasswordHash",
				"Password":     "request.Password",
			},
		}},
	}}

	errs := CheckFuncs(sfs, specs, nil, nil, nil)
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "불일치") {
			t.Errorf("unexpected param count ERROR: %s", e.Message)
		}
	}
}
