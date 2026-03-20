//ff:func feature=ssac-validate type=test control=sequence
//ff:what Inputs에서 미선언 변수.필드 참조 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateArgVarRef(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "GetDetail", FileName: "get_detail.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"InstructorID": "course.InstructorID"}, Result: &parser.Result{Type: "User", Var: "user"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, `"course" 변수가 아직 선언되지 않았습니다`)
}
