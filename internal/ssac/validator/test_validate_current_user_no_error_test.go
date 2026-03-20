//ff:func feature=ssac-validate type=test control=sequence
//ff:what currentUser 참조 시 에러 없음 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateCurrentUserNoError(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "GetMy", FileName: "get_my.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Item.ListByUser", Inputs: map[string]string{"ID": "currentUser.ID"}, Result: &parser.Result{Type: "[]Item", Var: "items"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"items": "items"}},
		},
	}}
	errs := Validate(funcs)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}
