//ff:func feature=ssac-gen type=test control=sequence
//ff:what snake_case DDL 컬럼이 PascalCase struct 필드로 변환되는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGenerateRequestStructSnakeCase(t *testing.T) {
	st := &rule.Ground{
		Tables: map[string]rule.TableInfo{
			"bids": {Columns: map[string]string{"bid_amount": "int32", "id": "int64"}},
		},
		Ops: map[string]rule.OperationInfo{},
		Models:     map[string]rule.ModelInfo{},
	}
	sf := ssacparser.ServiceFunc{
		Name: "PlaceBid", FileName: "place_bid.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPost, Model: "Bid.Place", Inputs: map[string]string{"bid_amount": "request.bid_amount", "id": "request.id"}, Result: &ssacparser.Result{Type: "Bid", Var: "bid"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"bid": "bid"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, "`json:\"bid_amount\"`")
	assertContains(t, code, "BidAmount int32")
	assertContains(t, code, "`json:\"id\"`")
	assertContains(t, code, "ID ")
	assertContains(t, code, "bidAmount := req.BidAmount")
	assertContains(t, code, "id := req.ID")
}
