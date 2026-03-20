//ff:func feature=ssac-validate type=test control=sequence
//ff:what @empty에서 Message 누락 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateEmptyMissingMessage(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqEmpty, Target: "course"},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "Message 누락")
}
