//ff:func feature=ssac-gen type=test control=sequence
//ff:what @state ErrStatus 커스텀 상태 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateStateErrStatus(t *testing.T) {
	sf := parser.ServiceFunc{
		Name:     "ActivateWorkflow",
		FileName: "activate_workflow.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Workflow.FindByID", Inputs: map[string]string{"ID": "request.WorkflowID"}, Result: &parser.Result{Type: "Workflow", Var: "workflow"}},
			{Type: parser.SeqState, DiagramID: "workflow", Inputs: map[string]string{"Status": "workflow.Status"}, Transition: "Activate", Message: "Cannot transition", ErrStatus: 422},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "http.StatusUnprocessableEntity")
	assertNotContains(t, code, "http.StatusConflict")
}
