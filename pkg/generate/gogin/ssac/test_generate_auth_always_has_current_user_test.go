//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth는 항상 currentUser 추출 + UserID, Role을 포함하는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateAuthAlwaysHasCurrentUser(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CheckAccess", FileName: "check_access.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqAuth, Action: "read", Resource: "public", Inputs: map[string]string{"Key": "request.APIKey"}, Message: "Forbidden"},
		},
	}
	code := mustGenerate(t, sf, nil)
	// @auth는 항상 currentUser 추출 + UserID, Role 포함
	assertContains(t, code, `c.MustGet("currentUser")`)
	assertContains(t, code, `UserID: currentUser.ID`)
	assertContains(t, code, `Role: currentUser.Role`)
	assertContains(t, code, `authz.Check(authz.CheckRequest{`)
}
