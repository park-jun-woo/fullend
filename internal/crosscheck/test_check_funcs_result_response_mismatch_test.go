//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_ResultResponseMismatch: @call에 result 있지만 func spec에 Response 필드 없으면 ERROR 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_ResultResponseMismatch(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package:        "auth",
		Name:           "issueToken",
		RequestFields:  []funcspec.Field{{Name: "UserID", Type: "int64"}},
		ResponseFields: []funcspec.Field{}, // no response fields
		HasBody:        true,
	}}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Login",
		Sequences: []ssacparser.Sequence{{
			Type:  "call",
			Model: "auth.IssueToken",
			Inputs: map[string]string{
				"UserID": "user.ID",
			},
			Result: &ssacparser.Result{Var: "token", Type: "Token"}, // has result
		}},
	}}

	errs := CheckFuncs(sfs, specs, nil, nil, nil)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "Response 필드 없음") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected result/response mismatch ERROR, got: %+v", errs)
	}
}
