//ff:func feature=ssac-validate type=test control=sequence
//ff:what 존재하지 않는 모델 참조 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateModelNotFound(t *testing.T) {
	st := &SymbolTable{Models: map[string]ModelSymbol{}, Operations: map[string]OperationSymbol{}}
	funcs := []parser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, `"Course" 모델을 찾을 수 없습니다`)
}
