//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=ssac-ddl
//ff:what CheckSSaCDDL: 패키지 모델과 DB 모델 혼합 시 각각 올바르게 처리 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

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
