//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth CheckRequest 스타일 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateAuthCallStyle(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "AcceptProposal", FileName: "accept_proposal.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Gig.FindByID", Inputs: map[string]string{"ID": "request.GigID"}, Result: &ssacparser.Result{Type: "Gig", Var: "gig"}},
			{Type: ssacparser.SeqAuth, Action: "AcceptProposal", Resource: "gig", Inputs: map[string]string{"UserID": "currentUser.ID", "ResourceID": "gig.ClientID"}, Message: "Not authorized"},
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
