//ff:func feature=ssac-validate type=test control=sequence
//ff:what @response 빈 Fields 허용 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateResponseEmptyFieldsAllowed(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "DeleteRoom", FileName: "delete_room.go",
		Sequences: []parser.Sequence{{Type: parser.SeqResponse}},
	}}
	errs := Validate(funcs)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}
