//ff:func feature=ssac-gen type=test control=sequence
//ff:what @post 시퀀스의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGeneratePost(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CreateSession", FileName: "create_session.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPost, Model: "Session.Create", Inputs: map[string]string{"ProjectID": "request.ProjectID", "Command": "request.Command"}, Result: &ssacparser.Result{Type: "Session", Var: "session"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"session": "session"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `session, err := h.SessionModel.WithTx(tx).Create(command, projectID)`)
	assertContains(t, code, `"session": session`)
	assertContains(t, code, `h.DB.BeginTx`)
	assertContains(t, code, `tx.Commit()`)
}
