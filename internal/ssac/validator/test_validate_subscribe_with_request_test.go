//ff:func feature=ssac-validate type=test control=sequence
//ff:what @subscribe 함수에서 request 사용 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateSubscribeWithRequest(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "OnOrder", FileName: "on_order.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.completed", MessageType: "Msg"},
		Param: &parser.ParamInfo{TypeName: "Msg", VarName: "message"},
		Structs: []parser.StructInfo{{Name: "Msg", Fields: []parser.StructField{{Name: "OrderID", Type: "int64"}}}},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "request.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "@subscribe 함수에서 request를 사용할 수 없습니다")
}
