//ff:func feature=ssac-gen type=test control=sequence
//ff:what @publish 시퀀스의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGeneratePublish(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CompleteOrder", FileName: "complete_order.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "request.OrderID"}, Result: &ssacparser.Result{Type: "Order", Var: "order"}},
			{Type: ssacparser.SeqPublish, Topic: "order.completed", Inputs: map[string]string{"OrderID": "order.ID", "Email": "order.Email"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"order": "order"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `queue.Publish(c.Request.Context(), "order.completed"`)
	assertContains(t, code, `"OrderID": order.ID`)
	assertContains(t, code, `order.Email`)
	assertContains(t, code, `"queue"`)
}
