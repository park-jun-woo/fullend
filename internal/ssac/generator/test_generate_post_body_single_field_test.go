//ff:func feature=ssac-gen type=test control=sequence
//ff:what path param + body field 혼합 시 JSON body 필드만 request struct에 포함되는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGeneratePostBodySingleField(t *testing.T) {
	st := &validator.SymbolTable{
		Models:    map[string]validator.ModelSymbol{},
		DDLTables: map[string]validator.DDLTable{
			"proposals": {Columns: map[string]string{"bid_amount": "int64", "gig_id": "int64", "freelancer_id": "int64"}},
		},
		Operations: map[string]validator.OperationSymbol{
			"SubmitProposal": {
				PathParams:     []validator.PathParam{{Name: "ID", GoType: "int64"}},
				HasRequestBody: true,
			},
		},
	}
	sf := parser.ServiceFunc{
		Name: "SubmitProposal", FileName: "submit_proposal.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Gig.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Gig", Var: "gig"}},
			{Type: parser.SeqPost, Model: "Proposal.Create", Inputs: map[string]string{"GigID": "gig.ID", "FreelancerID": "currentUser.ID", "BidAmount": "request.BidAmount"}, Result: &parser.Result{Type: "Proposal", Var: "proposal"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"proposal": "proposal"}},
		},
	}
	code := mustGenerate(t, sf, st)
	// BidAmount는 path param이 아니므로 JSON body에서 읽어야 함
	assertContains(t, code, `ShouldBindJSON(&req)`)
	assertContains(t, code, `BidAmount int64`)
	assertNotContains(t, code, `c.Query("BidAmount")`)
}
