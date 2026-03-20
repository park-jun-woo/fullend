//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe 함수의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateSubscribeFunc(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "OnOrderCompleted", FileName: "on_order_completed.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.completed", MessageType: "OnOrderCompletedMessage"},
		Param:     &parser.ParamInfo{TypeName: "OnOrderCompletedMessage", VarName: "message"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "message.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
			{Type: parser.SeqPut, Model: "Order.UpdateNotified", Inputs: map[string]string{"ID": "order.ID", "Notified": `"true"`}},
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
