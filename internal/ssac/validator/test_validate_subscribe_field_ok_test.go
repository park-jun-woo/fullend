//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what @subscribe 메시지 필드 매칭 성공 시 에러 없음 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateSubscribeFieldOK(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "OnOrder", FileName: "on_order.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.completed", MessageType: "Msg"},
		Param: &parser.ParamInfo{TypeName: "Msg", VarName: "message"},
		Structs: []parser.StructInfo{{Name: "Msg", Fields: []parser.StructField{{Name: "OrderID", Type: "int64"}, {Name: "Email", Type: "string"}}}},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
			{Type: parser.SeqPut, Model: "Order.Update", Inputs: map[string]string{"Email": "message.Email"}},
		},
	}}
	errs := Validate(funcs)
	for _, e := range errs {
		if contains(e.Message, "필드가 없습니다") {
			t.Errorf("unexpected field error: %s", e.Message)
		}
	}
}
