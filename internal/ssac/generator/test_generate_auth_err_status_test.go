//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth ErrStatus 커스텀 상태 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateAuthErrStatus(t *testing.T) {
	sf := parser.ServiceFunc{
		Name:     "Execute",
		FileName: "execute.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqAuth, Action: "Execute", Resource: "workflow", Inputs: map[string]string{"UserID": "currentUser.ID"}, Message: "Token expired", ErrStatus: 401},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "http.StatusUnauthorized")
	assertNotContains(t, code, "http.StatusForbidden")
}
