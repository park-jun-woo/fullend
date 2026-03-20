//ff:func feature=ssac-gen type=test control=sequence
//ff:what x-sort 없을 때 SortConfig가 생성되지 않는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateSortNoXSort(t *testing.T) {
	st := &validator.SymbolTable{
		Models:     map[string]validator.ModelSymbol{},
		DDLTables:  map[string]validator.DDLTable{},
		Operations: map[string]validator.OperationSymbol{
			"ListGigs": {
				XPagination: &validator.XPagination{Style: "offset", DefaultLimit: 20, MaxLimit: 100},
			},
		},
	}
	sf := parser.ServiceFunc{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Gig.List", Inputs: map[string]string{"Query": "query"}, Result: &parser.Result{Type: "[]Gig", Var: "gigs"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"gigs": "gigs"}},
		},
	}
	code := mustGenerate(t, sf, st)
	// x-sort 없으면 SortConfig 없음
	assertNotContains(t, code, `Sort:`)
	assertNotContains(t, code, `SortConfig`)
}
