//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what @call 결과 없이 호출만 하면 에러 없음 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateCallNoResultOK(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Notify", FileName: "notify.go",
		Sequences: []parser.Sequence{{
			Type:   parser.SeqCall,
			Model:  "notification.Send",
			Inputs: map[string]string{"ID": "reservation.ID"},
		}},
	}}
	errs := Validate(funcs)
	for _, e := range errs {
		if contains(e.Message, "기본 타입") {
			t.Errorf("unexpected primitive type error: %s", e.Message)
		}
	}
}
