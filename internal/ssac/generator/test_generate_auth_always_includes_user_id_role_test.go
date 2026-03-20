//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth 템플릿에 항상 UserID, Role이 포함되는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateAuthAlwaysIncludesUserIDRole(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "CheckAccess", FileName: "check_access.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqAuth, Action: "read", Resource: "public", Inputs: map[string]string{"Key": "request.APIKey"}, Message: "Forbidden"},
		},
	}
	code := mustGenerate(t, sf, nil)
	// @auth 템플릿에 항상 UserID, Role 포함
	assertContains(t, code, `UserID: currentUser.ID`)
	assertContains(t, code, `Role: currentUser.Role`)
}
