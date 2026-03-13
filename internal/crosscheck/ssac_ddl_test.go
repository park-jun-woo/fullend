package crosscheck

import (
	"testing"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func TestCheckSSaCDDL_PackageModelSkipped(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {Columns: map[string]string{"id": "int64"}},
		},
	}

	// Package model should be skipped — no DDL warning for "Session" type.
	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetSession",
		FileName: "service.go",
		Sequences: []ssacparser.Sequence{{
			Type:    "get",
			Package: "session",
			Model:   "Session.Get",
			Inputs:  map[string]string{"token": "request.Token"},
			Result:  &ssacparser.Result{Var: "session", Type: "Session"},
		}},
	}}

	errs := CheckSSaCDDL(funcs, st, nil)
	for _, e := range errs {
		if contains(e.Message, "Session") || contains(e.Message, "sessions") {
			t.Errorf("package model should skip DDL check: %+v", e)
		}
	}
}

func TestCheckSSaCDDL_DBModelChecked(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {Columns: map[string]string{"id": "int64", "email": "string"}},
		},
	}

	// Non-package model with unknown result type should get WARNING.
	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetGig",
		FileName: "service.go",
		Sequences: []ssacparser.Sequence{{
			Type:   "get",
			Model:  "Gig.FindByID",
			Inputs: map[string]string{"ID": "request.ID"},
			Result: &ssacparser.Result{Var: "gig", Type: "Gig"},
		}},
	}}

	errs := CheckSSaCDDL(funcs, st, nil)
	found := false
	for _, e := range errs {
		if e.Rule == "SSaC @result ↔ DDL" && contains(e.Message, "gigs") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected DDL table warning for Gig, got: %+v", errs)
	}
}

func TestCheckSSaCDDL_MixedPackageAndDB(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {Columns: map[string]string{"id": "int64", "email": "string"}},
		},
	}

	funcs := []ssacparser.ServiceFunc{{
		Name:     "ComplexHandler",
		FileName: "service.go",
		Sequences: []ssacparser.Sequence{
			{
				// Package model — should be skipped.
				Type:    "get",
				Package: "cache",
				Model:   "Cache.Get",
				Inputs:  map[string]string{"key": "request.Key"},
				Result:  &ssacparser.Result{Var: "cached", Type: "CachedGig"},
			},
			{
				// DB model — should be checked.
				Type:   "get",
				Model:  "User.FindByID",
				Inputs: map[string]string{"ID": "request.ID"},
				Result: &ssacparser.Result{Var: "user", Type: "User"},
			},
		},
	}}

	errs := CheckSSaCDDL(funcs, st, nil)
	// No error for CachedGig (package model skipped).
	for _, e := range errs {
		if contains(e.Message, "CachedGig") || contains(e.Message, "cached_gigs") {
			t.Errorf("package model result should skip DDL check: %+v", e)
		}
	}
	// No error for User (DDL table exists).
	for _, e := range errs {
		if contains(e.Message, "users") {
			t.Errorf("User table exists, should not get DDL error: %+v", e)
		}
	}
}
