//ff:func feature=ssac-gen type=test control=sequence
//ff:what response에서 참조된 변수는 변수명을 유지하는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateUsedVar(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "ProcessOrder", FileName: "process_order.go",
		Imports: []string{"myapp/billing"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "request.OrderID"}, Result: &ssacparser.Result{Type: "Order", Var: "order"}},
			{Type: ssacparser.SeqCall, Model: "billing.HoldEscrow", Inputs: map[string]string{"Amount": "order.Budget"}, Result: &ssacparser.Result{Type: "Escrow", Var: "escrow"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"order": "order", "escrow": "escrow"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// escrow는 response에서 참조됨 → 변수명 유지
	assertContains(t, code, `escrow, err := billing.HoldEscrow(billing.HoldEscrowRequest{`)
	assertContains(t, code, `order, err := h.OrderModel.FindByID`)
}
