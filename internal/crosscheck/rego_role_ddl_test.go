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
