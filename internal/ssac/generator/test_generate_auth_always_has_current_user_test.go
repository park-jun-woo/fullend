//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth는 항상 currentUser 추출 + UserID, Role을 포함하는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateAuthAlwaysHasCurrentUser(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "CheckAccess", FileName: "check_access.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqAuth, Action: "read", Resource: "public", Inputs: map[string]string{"Key": "request.APIKey"}, Message: "Forbidden"},
		},
	}
	code := mustGenerate(t, sf, nil)
	// @auth는 항상 currentUser 추출 + UserID, Role 포함
	assertContains(t, code, `c.MustGet("currentUser")`)
	assertContains(t, code, `UserID: currentUser.ID`)
	assertContains(t, code, `Role: currentUser.Role`)
	assertContains(t, code, `authz.Check(authz.CheckRequest{`)
}
