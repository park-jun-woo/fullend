//ff:func feature=ssac-gen type=test control=sequence
//ff:what snake_case DDL 컬럼이 PascalCase struct 필드로 변환되는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateRequestStructSnakeCase(t *testing.T) {
	st := &validator.SymbolTable{
		DDLTables: map[string]validator.DDLTable{
			"bids": {Columns: map[string]string{"bid_amount": "int32", "id": "int64"}},
		},
		Operations: map[string]validator.OperationSymbol{},
		Models:     map[string]validator.ModelSymbol{},
	}
	sf := parser.ServiceFunc{
		Name: "PlaceBid", FileName: "place_bid.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqPost, Model: "Bid.Place", Inputs: map[string]string{"bid_amount": "request.bid_amount", "id": "request.id"}, Result: &parser.Result{Type: "Bid", Var: "bid"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"bid": "bid"}},
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
