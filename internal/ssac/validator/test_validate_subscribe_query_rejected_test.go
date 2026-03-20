//ff:func feature=ssac-validate type=test control=sequence
//ff:what @subscribe에서 query 사용 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateSubscribeQueryRejected(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "OnOrder", FileName: "on_order.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.created", MessageType: "OrderMsg"},
		Param: &parser.ParamInfo{TypeName: "OrderMsg", VarName: "message"},
		Structs: []parser.StructInfo{{Name: "OrderMsg", Fields: []parser.StructField{{Name: "ID", Type: "string"}}}},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.ID", "Filter": "query"}, Result: &parser.Result{Type: "Order", Var: "order"}},
		},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "query는 HTTP 전용입니다")
}
