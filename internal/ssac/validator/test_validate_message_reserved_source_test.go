//ff:func feature=ssac-validate type=test control=sequence
//ff:what message를 result 변수명으로 사용 시 예약 소스 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateMessageReservedSource(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Test", FileName: "test.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "User", Var: "message"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "예약 소스이므로 result 변수명으로 사용할 수 없습니다")
}
