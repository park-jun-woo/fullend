//ff:func feature=ssac-validate type=test control=sequence
//ff:what DDL column "range"가 Go 예약어 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateGoReservedWordRange(t *testing.T) {
	st := &SymbolTable{
		Models: map[string]ModelSymbol{"Schedule": {Methods: map[string]MethodInfo{"Create": {Cardinality: "exec"}}}},
		Operations: map[string]OperationSymbol{},
		DDLTables: map[string]DDLTable{"schedules": {Columns: map[string]string{"range": "string", "name": "string"}}},
	}
	funcs := []parser.ServiceFunc{{Name: "CreateSchedule", FileName: "create_schedule.go", Sequences: []parser.Sequence{
		{Type: parser.SeqPost, Model: "Schedule.Create", Inputs: map[string]string{"range": "request.Range", "name": "request.Name"}, Result: &parser.Result{Type: "Schedule", Var: "schedule"}},
	}}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, `DDL column "range" in table "schedules" is a Go reserved word`)
}
