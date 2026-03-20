//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what @subscribe에서 message는 선언 없이 사용 가능한지 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateSubscribeMessageVariable(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "OnOrder", FileName: "on_order.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.completed", MessageType: "Msg"},
		Param: &parser.ParamInfo{TypeName: "Msg", VarName: "message"},
		Structs: []parser.StructInfo{{Name: "Msg", Fields: []parser.StructField{{Name: "OrderID", Type: "int64"}}}},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
		},
	}}
	errs := Validate(funcs)
	for _, e := range errs {
		if contains(e.Message, `"message" 변수가 아직 선언되지 않았습니다`) {
			t.Errorf("message should be pre-declared in subscribe func: %s", e.Message)
		}
	}
}
