//ff:func feature=ssac-validate type=test control=sequence
//ff:what @get에서 Result 누락 시 ERROR 검증
package validator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateGetMissingResult(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{{
			Type:   parser.SeqGet,
			Model:  "Course.FindByID",
			Inputs: map[string]string{"CourseID": "request.CourseID"},
		}},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "Result 누락")
}
