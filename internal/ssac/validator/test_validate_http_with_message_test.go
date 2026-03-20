//ff:func feature=ssac-validate type=test control=sequence
//ff:what HTTP 함수에서 message 사용 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateHTTPWithMessage(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "GetOrder", FileName: "get_order.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "HTTP 함수에서 message를 사용할 수 없습니다")
}
