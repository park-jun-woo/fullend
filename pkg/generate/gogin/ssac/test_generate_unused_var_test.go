//ff:func feature=ssac-gen type=test control=sequence
//ff:what response에서 미참조 변수를 _ 처리하는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateUnusedVar(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "ProcessOrder", FileName: "process_order.go",
		Imports: []string{"myapp/billing"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "request.OrderID"}, Result: &ssacparser.Result{Type: "Order", Var: "order"}},
			{Type: ssacparser.SeqCall, Model: "billing.HoldEscrow", Inputs: map[string]string{"Amount": "order.Budget"}, Result: &ssacparser.Result{Type: "Escrow", Var: "escrow"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"order": "order"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// escrow는 response에서 미참조 → _, err already declared → =
	assertContains(t, code, `_, err = billing.HoldEscrow(billing.HoldEscrowRequest{`)
	// order는 response에서 참조 → 변수명 유지
	assertContains(t, code, `order, err := h.OrderModel.FindByID`)
}
