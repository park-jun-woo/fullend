//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-ddl
//ff:what CheckDDLCoverageDBModelReferenced: DB 모델이 참조될 때 미참조 경고가 없는지 테스트
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

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
