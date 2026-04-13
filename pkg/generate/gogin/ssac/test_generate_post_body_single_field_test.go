//ff:func feature=ssac-gen type=test control=sequence
//ff:what path param + body field 혼합 시 JSON body 필드만 request struct에 포함되는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGeneratePostBodySingleField(t *testing.T) {
	st := &rule.Ground{
		Models:    map[string]rule.ModelInfo{},
		Tables: map[string]rule.TableInfo{
			"proposals": {Columns: map[string]string{"bid_amount": "int64", "gig_id": "int64", "freelancer_id": "int64"}},
		},
		Ops: map[string]rule.OperationInfo{
			"SubmitProposal": {
				PathParams:     []rule.PathParam{{Name: "ID", GoType: "int64"}},
				HasRequestBody: true,
			},
		},
	}
	sf := ssacparser.ServiceFunc{
		Name: "SubmitProposal", FileName: "submit_proposal.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Gig.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &ssacparser.Result{Type: "Gig", Var: "gig"}},
			{Type: ssacparser.SeqPost, Model: "Proposal.Create", Inputs: map[string]string{"GigID": "gig.ID", "FreelancerID": "currentUser.ID", "BidAmount": "request.BidAmount"}, Result: &ssacparser.Result{Type: "Proposal", Var: "proposal"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"proposal": "proposal"}},
		},
	}
	code := mustGenerate(t, sf, st)
	// BidAmount는 path param이 아니므로 JSON body에서 읽어야 함
	assertContains(t, code, `ShouldBindJSON(&req)`)
	assertContains(t, code, `BidAmount int64`)
	assertNotContains(t, code, `c.Query("BidAmount")`)
}
