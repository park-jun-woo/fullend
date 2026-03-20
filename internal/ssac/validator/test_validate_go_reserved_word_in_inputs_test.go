//ff:func feature=ssac-validate type=test control=sequence
//ff:what DDL column이 Go 예약어이면 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateGoReservedWordInInputs(t *testing.T) {
	st := &SymbolTable{
		Models: map[string]ModelSymbol{"Transaction": {Methods: map[string]MethodInfo{"Create": {Cardinality: "exec"}}}},
		Operations: map[string]OperationSymbol{},
		DDLTables: map[string]DDLTable{"transactions": {Columns: map[string]string{"type": "string", "amount": "int64", "gig_id": "int64"}}},
	}
	funcs := []parser.ServiceFunc{{Name: "CreateTransaction", FileName: "create_transaction.go", Sequences: []parser.Sequence{
		{Type: parser.SeqPost, Model: "Transaction.Create", Inputs: map[string]string{"amount": "request.Amount", "gigID": "request.GigID", "type": "request.Type"}, Result: &parser.Result{Type: "Transaction", Var: "tx"}},
	}}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, `DDL column "type" in table "transactions" is a Go reserved word`)
}
