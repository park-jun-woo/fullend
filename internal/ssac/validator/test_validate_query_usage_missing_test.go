//ff:func feature=ssac-validate type=test control=sequence
//ff:what x-pagination 있지만 query 미사용 시 WARNING 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateQueryUsageMissing(t *testing.T) {
	st := &SymbolTable{
		Models: map[string]ModelSymbol{}, DDLTables: map[string]DDLTable{},
		Operations: map[string]OperationSymbol{"ListReservations": {XPagination: &XPagination{Style: "offset", DefaultLimit: 20, MaxLimit: 100}}},
	}
	funcs := []parser.ServiceFunc{{
		Name: "ListReservations", FileName: "list_reservations.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Reservation.List", Result: &parser.Result{Type: "[]Reservation", Var: "reservations"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasWarning(t, errs, "query가 사용되지 않았습니다")
}
