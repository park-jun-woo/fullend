//ff:func feature=ssac-gen type=test control=sequence
//ff:what database/sql import가 불필요할 때 제거되는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateNoUnusedImportDatabaseSQL(t *testing.T) {
	// @post가 있으면 tx 코드가 생성되지만, database/sql은 handler.go에만 필요
	sf := ssacparser.ServiceFunc{
		Name: "CreateSession", FileName: "create_session.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPost, Model: "Session.Create", Inputs: map[string]string{"UserID": "request.UserID"}, Result: &ssacparser.Result{Type: "Session", Var: "session"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"session": "session"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertNotContains(t, code, `"database/sql"`)
	assertContains(t, code, `h.DB.BeginTx`)
}
