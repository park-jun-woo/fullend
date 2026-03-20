//ff:func feature=ssac-validate type=test control=sequence
//ff:what OpenAPI에 있지만 SSaC에서 미사용 필드 WARNING 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateReverseRequestMissing(t *testing.T) {
	st := &SymbolTable{
		Models:     map[string]ModelSymbol{"Course": {Methods: map[string]MethodInfo{"FindByID": {Cardinality: "one"}}}},
		Operations: map[string]OperationSymbol{"GetCourse": {RequestFields: map[string]bool{"CourseID": true, "Description": true}}},
	}
	funcs := []parser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"CourseID": "request.CourseID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasWarning(t, errs, `OpenAPI request에 "Description" 필드가 있지만 SSaC에서 사용하지 않습니다`)
}
