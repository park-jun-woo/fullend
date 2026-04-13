//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe 함수의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateSubscribeFunc(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "OnOrderCompleted", FileName: "on_order_completed.go",
		Subscribe: &ssacparser.SubscribeInfo{Topic: "order.completed", MessageType: "OnOrderCompletedMessage"},
		Param:     &ssacparser.ParamInfo{TypeName: "OnOrderCompletedMessage", VarName: "message"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.OrderID"}, Result: &ssacparser.Result{Type: "Order", Var: "order"}},
			{Type: ssacparser.SeqPut, Model: "Order.UpdateNotified", Inputs: map[string]string{"ID": "order.ID", "Notified": `"true"`}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "func (h *Handler) OnOrderCompleted(ctx context.Context, message OnOrderCompletedMessage) error {")
	assertContains(t, code, "return nil")
	assertContains(t, code, `return fmt.Errorf(`)
	assertContains(t, code, `"context"`)
	assertContains(t, code, `"fmt"`)
	assertNotContains(t, code, "gin.Context")
	assertNotContains(t, code, "c.JSON")
}
