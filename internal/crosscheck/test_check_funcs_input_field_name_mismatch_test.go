//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_InputFieldNameMismatch: @call 입력 키가 func spec Request 필드명에 없으면 ERROR 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_InputFieldNameMismatch(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package: "auth",
		Name:    "issueToken",
		RequestFields: []funcspec.Field{
			{Name: "UserID", Type: "int64"},
			{Name: "Email", Type: "string"},
		},
		HasBody: true,
	}}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Login",
		Sequences: []ssacparser.Sequence{
			{
				Type:   "get",
				Result: &ssacparser.Result{Var: "user", Type: "User"},
			},
			{
				Type:  "call",
				Model: "auth.IssueToken",
				Inputs: map[string]string{
					"ID":    "user.ID",    // wrong: should be UserID
					"Email": "user.Email", // correct
				},
			},
		},
	}}

	errs := CheckFuncs(sfs, specs, nil, nil, nil)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "Request에 없음") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected field name mismatch ERROR, got: %+v", errs)
	}
}
