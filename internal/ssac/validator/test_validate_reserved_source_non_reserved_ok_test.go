//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what 예약 소스가 아닌 변수명은 에러 없음 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateReservedSourceNonReservedOK(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Test", FileName: "test.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "User", Var: "user"}},
		},
	}}
	errs := Validate(funcs)
	for _, e := range errs {
		if contains(e.Message, "예약 소스") {
			t.Errorf("unexpected reserved source error: %s", e.Message)
		}
	}
}
