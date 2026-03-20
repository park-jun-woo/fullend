//ff:func feature=ssac-validate type=test control=sequence
//ff:what @response에서 미선언 변수 사용 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateUndeclaredInResponse(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, `"course" 변수가 아직 선언되지 않았습니다`)
}
