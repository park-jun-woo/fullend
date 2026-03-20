//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-ddl
//ff:what CheckDDLCoveragePackageModelSkipped: 패키지 모델만 참조할 때 DDL 테이블 미참조 에러를 생성하는지 테스트
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
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

	// Only a package model references -- no DDL model references.
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
