//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=policy-check
//ff:what TestCheckRegoRoleDDL_CaseMismatch: Rego 역할 값이 DDL CHECK enum과 대소문자 불일치 시 ERROR 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckRegoRoleDDL_CaseMismatch(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {
				Columns:    map[string]string{"id": "int64", "role": "string"},
				CheckEnums: map[string][]string{"role": {"admin", "member"}},
			},
		},
	}

	policies := []*policy.Policy{{
		File: "authz.rego",
		Rules: []policy.AllowRule{
			{RoleValue: "Admin", SourceLine: 5},
		},
	}}

	errs := CheckRegoRoleDDL(policies, st)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "Admin") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ERROR for role Admin not in DDL CHECK, got: %+v", errs)
	}
}
