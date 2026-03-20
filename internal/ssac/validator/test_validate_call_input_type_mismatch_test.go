//ff:func feature=ssac-validate type=test control=sequence
//ff:what @call 입력 타입 불일치 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateCallInputTypeMismatch(t *testing.T) {
	st := &SymbolTable{
		Models: map[string]ModelSymbol{"billing.Billing": {Methods: map[string]MethodInfo{"HoldEscrow": {Params: []string{"req"}, ParamTypes: map[string]string{"Amount": "int", "GigID": "int64"}}}}},
		Operations: map[string]OperationSymbol{},
		DDLTables: map[string]DDLTable{"gigs": {Columns: map[string]string{"id": "int64", "budget": "int64"}}},
	}
	funcs := []parser.ServiceFunc{{Name: "ProcessGig", FileName: "process_gig.go", Sequences: []parser.Sequence{
		{Type: parser.SeqGet, Model: "Gig.FindByID", Inputs: map[string]string{"ID": "request.GigID"}, Result: &parser.Result{Type: "Gig", Var: "gig"}},
		{Type: parser.SeqCall, Model: "billing.HoldEscrow", Inputs: map[string]string{"Amount": "gig.Budget", "GigID": "gig.ID"}},
	}}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, `타입 불일치`)
	assertHasError(t, errs, `int64`)
}
