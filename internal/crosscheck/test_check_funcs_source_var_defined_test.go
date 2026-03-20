//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_SourceVarDefined: @call 입력의 소스 변수가 이전 @result로 정의되면 WARNING 없음 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_SourceVarDefined(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package: "auth",
		Name:    "verifyPassword",
		RequestFields: []funcspec.Field{
			{Name: "PasswordHash", Type: "string"},
			{Name: "Password", Type: "string"},
		},
		HasBody: true,
	}}

	// Prior @result defines "user" variable.
	sfs := []ssacparser.ServiceFunc{{
		Name: "Login",
		Sequences: []ssacparser.Sequence{
			{
				Type:   "get",
				Result: &ssacparser.Result{Var: "user", Type: "User"},
			},
			{
				Type:  "call",
				Model: "auth.VerifyPassword",
				Inputs: map[string]string{
					"PasswordHash": "user.PasswordHash",
					"Password":     "request.Password",
				},
			},
		},
	}}

	errs := CheckFuncs(sfs, specs, nil, nil, nil)
	for _, e := range errs {
		if e.Level == "WARNING" && contains(e.Message, "미정의") {
			t.Errorf("unexpected source var WARNING: %s", e.Message)
		}
	}
}
