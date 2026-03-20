//ff:func feature=ssac-validate type=test control=sequence
//ff:what @call에서 Model 누락 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateCallMissingModel(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Login", FileName: "login.go",
		Sequences: []parser.Sequence{{Type: parser.SeqCall}},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "Model 누락")
}
