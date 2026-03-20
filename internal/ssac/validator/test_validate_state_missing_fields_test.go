//ff:func feature=ssac-validate type=test control=sequence
//ff:what @state 필수 필드 전부 누락 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateStateMissingFields(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Cancel", FileName: "cancel.go",
		Sequences: []parser.Sequence{{Type: parser.SeqState}},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "DiagramID 누락")
	assertHasError(t, errs, "Inputs 누락")
	assertHasError(t, errs, "Transition 누락")
	assertHasError(t, errs, "Message 누락")
}
