//ff:func feature=crosscheck type=rule control=sequence topic=policy-check
//ff:what TestCheckRegoRoleDDL_NoCheckEnum: DDL에 CHECK enum이 없으면 Rego 역할 검증을 건너뜀 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckRegoRoleDDL_NoCheckEnum(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {
				Columns: map[string]string{"id": "int64", "role": "string"},
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
	if len(errs) != 0 {
		t.Errorf("expected no errors when no CHECK enum, got: %+v", errs)
	}
}
