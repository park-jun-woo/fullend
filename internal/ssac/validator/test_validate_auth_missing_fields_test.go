//ff:func feature=ssac-validate type=test control=sequence
//ff:what @auth 필수 필드 전부 누락 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateAuthMissingFields(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Delete", FileName: "delete.go",
		Sequences: []parser.Sequence{{Type: parser.SeqAuth}},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "Action 누락")
	assertHasError(t, errs, "Resource 누락")
	assertHasError(t, errs, "Message 누락")
}
