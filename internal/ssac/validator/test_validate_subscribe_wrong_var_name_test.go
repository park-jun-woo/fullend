//ff:func feature=ssac-validate type=test control=sequence
//ff:what @subscribe 파라미터 변수명이 message가 아니면 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateSubscribeWrongVarName(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "OnOrder", FileName: "on_order.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.completed", MessageType: "Msg"},
		Param: &parser.ParamInfo{TypeName: "Msg", VarName: "msg"},
		Structs: []parser.StructInfo{{Name: "Msg", Fields: []parser.StructField{{Name: "OrderID", Type: "int64"}}}},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, `파라미터 변수명은 "message"여야 합니다`)
}
