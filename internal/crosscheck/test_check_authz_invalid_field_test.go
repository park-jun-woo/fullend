//ff:func feature=crosscheck type=rule control=sequence
//ff:what CheckAuthzInvalidField: 유효하지 않은 authz 필드에 대해 에러를 생성하는지 테스트
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckAuthzInvalidField(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "DoSomething",
			Sequences: []ssacparser.Sequence{
				{
					Type: "auth",
					Inputs: map[string]string{
						"UserID":   "currentUser.ID",
						"BadField": "gig.ClientID",
					},
				},
			},
		},
	}

	errs := CheckAuthz(funcs, "")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if errs[0].Level != "ERROR" {
		t.Fatalf("expected ERROR level, got %s", errs[0].Level)
	}
}
