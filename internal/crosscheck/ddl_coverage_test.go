package crosscheck

import (
	"testing"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func TestCheckDDLCoverage_PackageModelSkipped(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {
				Columns:     map[string]string{"id": "int64", "email": "string"},
				ColumnOrder: []string{"id", "email"},
			},
		},
	}

	// Only a package model references — no DDL model references.
	// "users" table should get a WARNING because it's not referenced by any non-package @model.
	funcs := []ssacparser.ServiceFunc{{
		Name: "GetSession",
		Sequences: []ssacparser.Sequence{{
			Type:    "get",
			Package: "session",
			Model:   "Session.Get",
			Inputs:  map[string]string{"token": "request.Token"},
			Result:  &ssacparser.Result{Var: "session", Type: "Session"},
		}},
	}}

	errs := CheckDDLCoverage(st, funcs, nil)
	found := false
	for _, e := range errs {
		if e.Rule == "DDL → SSaC" && e.Level == "ERROR" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected DDL table unreferenced ERROR, got: %+v", errs)
	}
}

func TestCheckDDLCoverage_DBModelReferenced(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {
				Columns:     map[string]string{"id": "int64", "email": "string"},
				ColumnOrder: []string{"id", "email"},
			},
		},
	}

	funcs := []ssacparser.ServiceFunc{{
		Name: "GetUser",
		Sequences: []ssacparser.Sequence{{
			Type:   "get",
			Model:  "User.FindByID",
			Inputs: map[string]string{"ID": "request.ID"},
			Result: &ssacparser.Result{Var: "user", Type: "User"},
		}},
	}}

	errs := CheckDDLCoverage(st, funcs, nil)
	for _, e := range errs {
		if e.Rule == "DDL → SSaC" && e.Level == "WARNING" && contains(e.Message, "users") {
			t.Errorf("unexpected unreferenced WARNING for users table: %+v", e)
		}
	}
}
