//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=ssac-ddl
//ff:what CheckSSaCDDL: DB 모델의 DDL 테이블 미존재 시 WARNING 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

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
