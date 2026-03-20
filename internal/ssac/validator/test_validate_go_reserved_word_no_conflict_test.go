//ff:func feature=ssac-validate type=test control=sequence
//ff:what Go 예약어가 아닌 DDL column은 에러 없음 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateGoReservedWordNoConflict(t *testing.T) {
	st := &SymbolTable{
		Models: map[string]ModelSymbol{"Transaction": {Methods: map[string]MethodInfo{"Create": {Cardinality: "exec"}}}},
		Operations: map[string]OperationSymbol{},
		DDLTables: map[string]DDLTable{"transactions": {Columns: map[string]string{"tx_type": "string", "amount": "int64"}}},
	}
	funcs := []parser.ServiceFunc{{Name: "CreateTransaction", FileName: "create_transaction.go", Sequences: []parser.Sequence{
		{Type: parser.SeqPost, Model: "Transaction.Create", Inputs: map[string]string{"amount": "request.Amount", "txType": "request.TxType"}, Result: &parser.Result{Type: "Transaction", Var: "tx"}},
	}}}
	errs := ValidateWithSymbols(funcs, st)
	assertNoErrors(t, errs)
}
