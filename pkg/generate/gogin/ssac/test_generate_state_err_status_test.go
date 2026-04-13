//ff:func feature=ssac-gen type=test control=sequence
//ff:what @state ErrStatus 커스텀 상태 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateStateErrStatus(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name:     "ActivateWorkflow",
		FileName: "activate_workflow.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Workflow.FindByID", Inputs: map[string]string{"ID": "request.WorkflowID"}, Result: &ssacparser.Result{Type: "Workflow", Var: "workflow"}},
			{Type: ssacparser.SeqState, DiagramID: "workflow", Inputs: map[string]string{"Status": "workflow.Status"}, Transition: "Activate", Message: "Cannot transition", ErrStatus: 422},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "http.StatusUnprocessableEntity")
	assertNotContains(t, code, "http.StatusConflict")
}
