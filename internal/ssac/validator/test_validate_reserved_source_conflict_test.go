//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what 예약 소스명을 result 변수로 사용 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateReservedSourceConflict(t *testing.T) {
	for _, name := range []string{"request", "currentUser", "config"} {
		t.Run(name, func(t *testing.T) {
			funcs := []parser.ServiceFunc{{
				Name: "Test", FileName: "test.go",
				Sequences: []parser.Sequence{
					{Type: parser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "User", Var: name}},
				},
			}}
			errs := Validate(funcs)
			assertHasError(t, errs, "예약 소스이므로 result 변수명으로 사용할 수 없습니다")
		})
	}
}
