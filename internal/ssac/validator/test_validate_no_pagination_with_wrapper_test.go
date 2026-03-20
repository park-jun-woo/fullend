//ff:func feature=ssac-validate type=test control=sequence
//ff:what x-pagination 없이 Page wrapper 사용 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateNoPaginationWithWrapper(t *testing.T) {
	st := &SymbolTable{
		Models: map[string]ModelSymbol{"Gig": {Methods: map[string]MethodInfo{"List": {Cardinality: "many"}}}},
		DDLTables: map[string]DDLTable{}, Operations: map[string]OperationSymbol{"ListGigs": {}},
	}
	funcs := []parser.ServiceFunc{{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []parser.Sequence{{Type: parser.SeqGet, Model: "Gig.List", Result: &parser.Result{Type: "Gig", Var: "gigPage", Wrapper: "Page"}}},
	}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, "x-pagination이 없지만")
}
