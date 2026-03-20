//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what SuppressWarn으로 전체 삭제 WARNING 억제 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateDeleteNoArgsSuppressed(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "DeleteAll", FileName: "delete_all.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqDelete, Model: "Room.DeleteAll", SuppressWarn: true},
		},
	}}
	errs := Validate(funcs)
	for _, e := range errs {
		if e.IsWarning() && contains(e.Message, "전체 삭제") {
			t.Errorf("expected WARNING to be suppressed: %s", e.Message)
		}
	}
}
