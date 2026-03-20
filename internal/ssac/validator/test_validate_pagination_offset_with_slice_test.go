//ff:func feature=ssac-validate type=test control=sequence
//ff:what offset pagination + slice 결과 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidatePaginationOffsetWithSlice(t *testing.T) {
	st := &SymbolTable{
		Models: map[string]ModelSymbol{"Gig": {Methods: map[string]MethodInfo{"List": {Cardinality: "many"}}}},
		DDLTables: map[string]DDLTable{}, Operations: map[string]OperationSymbol{"ListGigs": {XPagination: &XPagination{Style: "offset"}}},
	}
	funcs := []parser.ServiceFunc{{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []parser.Sequence{{Type: parser.SeqGet, Model: "Gig.List", Result: &parser.Result{Type: "[]Gig", Var: "gigs"}}},
	}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, "Page[T]가 아닙니다")
}
