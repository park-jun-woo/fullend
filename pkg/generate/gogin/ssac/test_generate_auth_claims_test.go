//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth claims에서 UserID가 중복 없이 1번만 나오는지 검증
package ssac

import (
	"strings"
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateAuthClaims(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "AcceptProposal", FileName: "accept_proposal.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Gig.FindByID", Inputs: map[string]string{"ID": "request.GigID"}, Result: &ssacparser.Result{Type: "Gig", Var: "gig"}},
			{Type: ssacparser.SeqAuth, Action: "AcceptProposal", Resource: "gig", Inputs: map[string]string{"UserID": "currentUser.ID", "ResourceID": "gig.ClientID"}, Message: "Not authorized"},
		},
	}
	code := mustGenerate(t, sf, nil)
	// 템플릿 고정 UserID, Role — inputs에 명시해도 중복 없음
	assertContains(t, code, `UserID: currentUser.ID`)
	assertContains(t, code, `Role: currentUser.Role`)
	// UserID가 1번만 나오는지 확인 (중복 방지)
	if strings.Count(code, "UserID:") != 1 {
		t.Errorf("expected UserID: to appear exactly once, got %d\n--- code ---\n%s", strings.Count(code, "UserID:"), code)
	}
}
