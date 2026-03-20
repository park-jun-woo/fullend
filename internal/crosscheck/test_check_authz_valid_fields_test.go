//ff:func feature=crosscheck type=rule control=sequence
//ff:what CheckAuthzValidFields: 유효한 authz 필드에 대해 에러가 없는지 테스트
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckAuthzValidFields(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "AcceptProposal",
			Sequences: []ssacparser.Sequence{
				{
					Type: "auth",
					Inputs: map[string]string{
						"UserID":     "currentUser.ID",
						"ResourceID": "gig.ClientID",
					},
				},
			},
		},
	}

	errs := CheckAuthz(funcs, "")
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %d: %v", len(errs), errs)
	}
}
