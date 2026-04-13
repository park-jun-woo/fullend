//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth ErrStatus 커스텀 상태 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateAuthErrStatus(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name:     "Execute",
		FileName: "execute.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqAuth, Action: "Execute", Resource: "workflow", Inputs: map[string]string{"UserID": "currentUser.ID"}, Message: "Token expired", ErrStatus: 401},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "http.StatusUnauthorized")
	assertNotContains(t, code, "http.StatusForbidden")
}
