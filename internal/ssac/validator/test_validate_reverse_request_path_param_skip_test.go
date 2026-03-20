//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what path param은 역방향 request 체크에서 스킵하는지 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateReverseRequestPathParamSkip(t *testing.T) {
	st := &SymbolTable{
		Models:     map[string]ModelSymbol{"Course": {Methods: map[string]MethodInfo{"FindByID": {Cardinality: "one"}}}},
		Operations: map[string]OperationSymbol{"GetCourse": {RequestFields: map[string]bool{"CourseID": true}, PathParams: []PathParam{{Name: "CourseID", GoType: "int64"}}}},
	}
	funcs := []parser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"CourseID": "request.CourseID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	for _, e := range errs {
		if e.IsWarning() && contains(e.Message, "CourseID") {
			t.Errorf("path param should be skipped in reverse check: %s", e.Message)
		}
	}
}
