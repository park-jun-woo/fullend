//ff:func feature=ssac-validate type=test control=sequence
//ff:what OpenAPI에 x-pagination 없이 query 사용 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateQueryUsageNoOpenAPI(t *testing.T) {
	st := &SymbolTable{Models: map[string]ModelSymbol{}, Operations: map[string]OperationSymbol{}, DDLTables: map[string]DDLTable{}}
	funcs := []parser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID", "Opts": "query"}, Result: &parser.Result{Type: "Course", Var: "course"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, "OpenAPI에 x-pagination")
}
