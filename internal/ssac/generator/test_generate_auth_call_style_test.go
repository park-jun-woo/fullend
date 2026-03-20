//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth CheckRequest 스타일 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateAuthCallStyle(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "AcceptProposal", FileName: "accept_proposal.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Gig.FindByID", Inputs: map[string]string{"ID": "request.GigID"}, Result: &parser.Result{Type: "Gig", Var: "gig"}},
			{Type: parser.SeqAuth, Action: "AcceptProposal", Resource: "gig", Inputs: map[string]string{"UserID": "currentUser.ID", "ResourceID": "gig.ClientID"}, Message: "Not authorized"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `authz.Check(authz.CheckRequest{Action: "AcceptProposal", Resource: "gig"`)
	assertContains(t, code, `ResourceID: gig.ClientID`)
	assertContains(t, code, `Role: currentUser.Role`)
	assertContains(t, code, `UserID: currentUser.ID`)
	assertContains(t, code, `http.StatusForbidden`)
	assertNotContains(t, code, `authz.Input{`)
}
