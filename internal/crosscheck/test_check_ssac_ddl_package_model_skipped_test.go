//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=ssac-ddl
//ff:what CheckSSaCDDL: 패키지 모델은 DDL 검사 스킵 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
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
