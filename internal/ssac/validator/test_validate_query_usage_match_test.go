//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what x-pagination + query 일치 시 에러 없음 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateQueryUsageMatch(t *testing.T) {
	st := &SymbolTable{
		Models: map[string]ModelSymbol{}, DDLTables: map[string]DDLTable{},
		Operations: map[string]OperationSymbol{"ListReservations": {XPagination: &XPagination{Style: "offset", DefaultLimit: 20, MaxLimit: 100}}},
	}
	funcs := []parser.ServiceFunc{{
		Name: "ListReservations", FileName: "list_reservations.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Reservation.List", Inputs: map[string]string{"Opts": "query"}, Result: &parser.Result{Type: "[]Reservation", Var: "reservations"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	for _, e := range errs {
		if contains(e.Message, "query") {
			t.Errorf("unexpected query validation error: %s", e.Message)
		}
	}
}
