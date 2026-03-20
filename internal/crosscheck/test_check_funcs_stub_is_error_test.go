//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_StubIsError: func spec이 stub(HasBody=false)이면 ERROR 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_StubIsError(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package: "auth",
		Name:    "verifyPassword",
		RequestFields: []funcspec.Field{
			{Name: "PasswordHash", Type: "string"},
			{Name: "Password", Type: "string"},
		},
		HasBody: false, // stub
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
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "TODO") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected stub TODO ERROR, got: %+v", errs)
	}
	// Ensure it's NOT WARNING.
	for _, e := range errs {
		if e.Level == "WARNING" && contains(e.Message, "TODO") {
			t.Errorf("stub should be ERROR not WARNING: %+v", e)
		}
	}
}
