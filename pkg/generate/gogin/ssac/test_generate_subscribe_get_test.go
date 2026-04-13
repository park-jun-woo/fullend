//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe에서 @get + @empty의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateSubscribeGet(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "OnTest", FileName: "on_test.go",
		Subscribe: &ssacparser.SubscribeInfo{Topic: "test", MessageType: "TestMsg"},
		Param:     &ssacparser.ParamInfo{TypeName: "TestMsg", VarName: "message"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.OrderID"}, Result: &ssacparser.Result{Type: "Order", Var: "order"}},
			{Type: ssacparser.SeqEmpty, Target: "order", Message: "주문을 찾을 수 없습니다"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "h.OrderModel.FindByID(message.OrderID)")
	assertContains(t, code, `return fmt.Errorf("Order 조회 실패: %w", err)`)
	assertContains(t, code, `return fmt.Errorf("주문을 찾을 수 없습니다")`)
}
