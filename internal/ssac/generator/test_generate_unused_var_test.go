//ff:func feature=ssac-gen type=test control=sequence
//ff:what response에서 미참조 변수를 _ 처리하는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateUnusedVar(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "ProcessOrder", FileName: "process_order.go",
		Imports: []string{"myapp/billing"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "request.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
			{Type: parser.SeqCall, Model: "billing.HoldEscrow", Inputs: map[string]string{"Amount": "order.Budget"}, Result: &parser.Result{Type: "Escrow", Var: "escrow"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"order": "order"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// escrow는 response에서 미참조 → _, err already declared → =
	assertContains(t, code, `_, err = billing.HoldEscrow(billing.HoldEscrowRequest{`)
	// order는 response에서 참조 → 변수명 유지
	assertContains(t, code, `order, err := h.OrderModel.FindByID`)
}
