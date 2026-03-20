//ff:func feature=ssac-gen type=test control=sequence
//ff:what Page[T] wrapper 사용 시 total 없이 단일 반환하는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGeneratePageNoHasTotal(t *testing.T) {
	st := &validator.SymbolTable{
		Models:     map[string]validator.ModelSymbol{},
		DDLTables:  map[string]validator.DDLTable{},
		Operations: map[string]validator.OperationSymbol{
			"ListGigs": {XPagination: &validator.XPagination{Style: "offset", DefaultLimit: 20, MaxLimit: 100}},
		},
	}
	sf := parser.ServiceFunc{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Gig.List", Inputs: map[string]string{"Query": "query"}, Result: &parser.Result{Type: "Gig", Var: "gigPage", Wrapper: "Page"}},
			{Type: parser.SeqResponse, Target: "gigPage"},
		},
	}
	code := mustGenerate(t, sf, st)
	// Page[T]이면 3-tuple 아니라 단일 반환
	assertNotContains(t, code, "total")
	assertContains(t, code, `gigPage, err :=`)
}
