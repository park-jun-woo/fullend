//ff:func feature=ssac-validate type=test control=sequence
//ff:what @subscribe 메시지 타입이 struct로 선언되지 않으면 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateSubscribeTypeNotFound(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "OnOrder", FileName: "on_order.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.completed", MessageType: "NonExistentMsg"},
		Param: &parser.ParamInfo{TypeName: "NonExistentMsg", VarName: "message"},
		Structs: []parser.StructInfo{},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "struct로 선언되지 않았습니다")
}
