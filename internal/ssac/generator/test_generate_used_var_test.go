//ff:func feature=ssac-gen type=test control=sequence
//ff:what response에서 참조된 변수는 변수명을 유지하는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateUsedVar(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "ProcessOrder", FileName: "process_order.go",
		Imports: []string{"myapp/billing"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Order.FindByID", Inputs: map[string]string{"ID": "request.OrderID"}, Result: &parser.Result{Type: "Order", Var: "order"}},
			{Type: parser.SeqCall, Model: "billing.HoldEscrow", Inputs: map[string]string{"Amount": "order.Budget"}, Result: &parser.Result{Type: "Escrow", Var: "escrow"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"order": "order", "escrow": "escrow"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// escrow는 response에서 참조됨 → 변수명 유지
	assertContains(t, code, `escrow, err := billing.HoldEscrow(billing.HoldEscrowRequest{`)
	assertContains(t, code, `order, err := h.OrderModel.FindByID`)
}
