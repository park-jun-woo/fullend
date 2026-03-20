//ff:func feature=ssac-gen type=test control=sequence
//ff:what JSON body request 시 ShouldBindJSON 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateWithJSONBody(t *testing.T) {
	st := &validator.SymbolTable{
		Models: map[string]validator.ModelSymbol{},
		DDLTables: map[string]validator.DDLTable{
			"sessions": {Columns: map[string]string{"project_id": "int64", "command": "string"}},
		},
		Operations: map[string]validator.OperationSymbol{},
	}
	sf := parser.ServiceFunc{
		Name: "CreateSession", FileName: "create_session.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqPost, Model: "Session.Create", Inputs: map[string]string{"ProjectID": "request.ProjectID", "Command": "request.Command"}, Result: &parser.Result{Type: "Session", Var: "session"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"session": "session"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `ShouldBindJSON(&req)`)
	assertContains(t, code, `ProjectID int64`)
}
