//ff:func feature=ssac-validate type=test control=sequence
//ff:what @put에서 Result가 있으면 ERROR 검증
package validator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidatePutHasResult(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "UpdateCourse", FileName: "update_course.go",
		Sequences: []parser.Sequence{{
			Type:   parser.SeqPut,
			Model:  "Course.Update",
			Inputs: map[string]string{"Title": "request.Title"},
			Result: &parser.Result{Type: "Course", Var: "course"},
		}},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "Result는 nil이어야 함")
}
