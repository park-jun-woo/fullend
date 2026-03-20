//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=ssac-ddl
//ff:what CheckSSaCDDL: 복수형 결과 타입에 대한 단수형 WARNING 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckSSaCDDL_PluralResultType(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"gigs": {Columns: map[string]string{"id": "int64"}},
		},
	}

	// Plural result type "Gigs" should trigger singular WARNING.
	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetGig",
		FileName: "service.go",
		Sequences: []ssacparser.Sequence{{
			Type:   "get",
			Model:  "Gig.FindByID",
			Inputs: map[string]string{"ID": "request.ID"},
			Result: &ssacparser.Result{Var: "gig", Type: "Gigs"},
		}},
	}}

	errs := CheckSSaCDDL(funcs, st, nil)
	found := false
	for _, e := range errs {
		if e.Level == "WARNING" && contains(e.Message, "singular") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected singular WARNING for plural type Gigs, got: %+v", errs)
	}
}
