//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_InputTypeMismatch: @call 입력 타입이 DDL 컬럼 타입과 func spec 필드 타입 불일치 시 ERROR 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckFuncs_InputTypeMismatch(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package: "billing",
		Name:    "holdEscrow",
		RequestFields: []funcspec.Field{
			{Name: "GigID", Type: "int64"},
			{Name: "Amount", Type: "int"}, // DDL budget is int64
			{Name: "ClientID", Type: "int64"},
		},
		HasBody: true,
	}}

	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"gigs": {
				Columns: map[string]string{
					"id":        "int64",
					"budget":    "int64",
					"client_id": "int64",
				},
			},
		},
	}

	sfs := []ssacparser.ServiceFunc{{
		Name: "AcceptProposal",
		Sequences: []ssacparser.Sequence{
			{
				Type:   "get",
				Result: &ssacparser.Result{Var: "gig", Type: "Gig"},
			},
			{
				Type:  "call",
				Model: "billing.HoldEscrow",
				Inputs: map[string]string{
					"GigID":    "gig.ID",
					"Amount":   "gig.Budget",   // int64 vs int → mismatch
					"ClientID": "gig.ClientID", // int64 vs int64 → ok
				},
			},
		},
	}}

	errs := CheckFuncs(sfs, specs, nil, st, nil)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "타입 불일치") && contains(e.Message, "Amount") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected type mismatch ERROR for Amount, got: %+v", errs)
	}
}
