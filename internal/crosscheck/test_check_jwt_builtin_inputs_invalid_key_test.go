//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what TestCheckJWTBuiltinInputs_InvalidKey: JWT 내장 함수에 claims에 없는 키를 전달하면 ERROR 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/projectconfig"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckJWTBuiltinInputs_InvalidKey(t *testing.T) {
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
				"OrgId": "user.OrgID", // OrgId instead of OrgID
			},
		}},
	}}

	errs := CheckJWTBuiltinInputs(sfs, claims)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "OrgId") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ERROR for OrgId not in claims, got: %+v", errs)
	}
}
