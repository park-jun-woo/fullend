//ff:func feature=ssac-gen type=test control=sequence
//ff:what JSON body request 시 ShouldBindJSON 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGenerateWithJSONBody(t *testing.T) {
	st := &rule.Ground{
		Models: map[string]rule.ModelInfo{},
		Tables: map[string]rule.TableInfo{
			"sessions": {Columns: map[string]string{"project_id": "int64", "command": "string"}},
		},
		Ops: map[string]rule.OperationInfo{},
	}
	sf := ssacparser.ServiceFunc{
		Name: "CreateSession", FileName: "create_session.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPost, Model: "Session.Create", Inputs: map[string]string{"ProjectID": "request.ProjectID", "Command": "request.Command"}, Result: &ssacparser.Result{Type: "Session", Var: "session"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"session": "session"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `ShouldBindJSON(&req)`)
	assertContains(t, code, `ProjectID int64`)
}
