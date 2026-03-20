//ff:func feature=ssac-gen type=test control=sequence
//ff:what query 입력 시 ParseQueryOpts 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateQueryArg(t *testing.T) {
	st := &validator.SymbolTable{
		Models:     map[string]validator.ModelSymbol{},
		Operations: map[string]validator.OperationSymbol{
			"ListMyReservations": {
				XPagination: &validator.XPagination{Style: "offset", DefaultLimit: 20, MaxLimit: 100},
			},
		},
		DDLTables: map[string]validator.DDLTable{},
	}
	sf := parser.ServiceFunc{
		Name: "ListMyReservations", FileName: "list_my_reservations.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Reservation.ListByUserID", Inputs: map[string]string{"UserID": "currentUser.ID", "Opts": "query"}, Result: &parser.Result{Type: "[]Reservation", Var: "reservations"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"reservations": "reservations"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `opts := model.ParseQueryOpts(c, model.QueryOptsConfig{`)
	assertContains(t, code, `PaginationConfig{Style: "offset"`)
	assertContains(t, code, `h.ReservationModel.ListByUserID(currentUser.ID, opts)`)
	assertContains(t, code, `reservations, total, err`)
	assertContains(t, code, `"total":`)
}
