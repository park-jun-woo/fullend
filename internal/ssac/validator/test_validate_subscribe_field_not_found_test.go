//ff:func feature=ssac-validate type=test control=sequence
//ff:what @subscribe 메시지 타입에 없는 필드 참조 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateSubscribeFieldNotFound(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "OnOrder", FileName: "on_order.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.completed", MessageType: "Msg"},
		Param: &parser.ParamInfo{TypeName: "Msg", VarName: "message"},
		Structs: []parser.StructInfo{{Name: "Msg", Fields: []parser.StructField{{Name: "OrderID", Type: "int64"}}}},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.Email"}, Result: &parser.Result{Type: "Order", Var: "order"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, `메시지 타입 "Msg"에 "Email" 필드가 없습니다`)
}
