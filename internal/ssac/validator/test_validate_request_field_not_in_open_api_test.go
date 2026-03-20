//ff:func feature=ssac-validate type=test control=sequence
//ff:what OpenAPI에 없는 request 필드 사용 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateRequestFieldNotInOpenAPI(t *testing.T) {
	st := &SymbolTable{
		Models:     map[string]ModelSymbol{"Course": {Methods: map[string]MethodInfo{"FindByID": {Cardinality: "one"}}}},
		Operations: map[string]OperationSymbol{"GetCourse": {RequestFields: map[string]bool{"CourseID": true}}},
	}
	funcs := []parser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"UnknownField": "request.UnknownField"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, `OpenAPI request에 "UnknownField" 필드가 없습니다`)
}
