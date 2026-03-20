//ff:func feature=ssac-validate type=test control=sequence
//ff:what Inputs에서 미선언 변수 참조 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateUndeclaredInInputs(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Cancel", FileName: "cancel.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqState, DiagramID: "reservation", Inputs: map[string]string{"status": "reservation.Status"}, Transition: "cancel", Message: "fail"},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, `"reservation" 변수가 아직 선언되지 않았습니다`)
}
