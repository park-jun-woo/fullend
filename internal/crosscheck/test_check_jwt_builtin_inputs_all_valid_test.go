//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what TestCheckJWTBuiltinInputs_AllValid: JWT 내장 함수에 모든 키가 claims에 있으면 에러 없음 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/projectconfig"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckJWTBuiltinInputs_AllValid(t *testing.T) {
	claims := map[string]projectconfig.ClaimDef{
		"ID":    {Key: "user_id", GoType: "int64"},
		"Email": {Key: "email", GoType: "string"},
		"Role":  {Key: "role", GoType: "string"},
		"OrgID": {Key: "org_id", GoType: "int64"},
	}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Login",
		Sequences: []ssacparser.Sequence{{
			Type:  "call",
			Model: "auth.IssueToken",
			Inputs: map[string]string{
				"ID":    "user.ID",
				"Email": "user.Email",
				"Role":  "user.Role",
				"OrgID": "user.OrgID",
			},
		}},
	}}

	errs := CheckJWTBuiltinInputs(sfs, claims)
	for _, e := range errs {
		if e.Level == "ERROR" {
			t.Errorf("expected no error for valid keys, got: %+v", e)
		}
	}
}
