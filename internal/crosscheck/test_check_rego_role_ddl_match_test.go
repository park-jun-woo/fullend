//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=policy-check
//ff:what TestCheckRegoRoleDDL_Match: Rego 역할 값이 DDL CHECK enum과 일치하면 에러 없음 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckRegoRoleDDL_Match(t *testing.T) {
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
			{RoleValue: "admin", SourceLine: 5},
		},
	}}

	errs := CheckRegoRoleDDL(policies, st)
	for _, e := range errs {
		if e.Level == "ERROR" {
			t.Errorf("expected no error for matching role, got: %+v", e)
		}
	}
}
