//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what 슬라이스 결과는 FK 참조 @empty 불필요 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateFKReferenceGetSliceSkipped(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "ListTasks", FileName: "list_tasks.go", Sequences: []parser.Sequence{
		{Type: parser.SeqGet, Model: "Project.FindByID", Inputs: map[string]string{"ID": "request.ProjectID"}, Result: &parser.Result{Type: "Project", Var: "project"}},
		{Type: parser.SeqEmpty, Target: "project", Message: "not found"},
		{Type: parser.SeqGet, Model: "Task.ListByProject", Inputs: map[string]string{"ProjectID": "project.ID"}, Result: &parser.Result{Type: "[]Task", Var: "tasks"}},
		{Type: parser.SeqResponse, Fields: map[string]string{"tasks": "tasks"}},
	}}}
	errs := Validate(funcs)
	for _, e := range errs {
		if contains(e.Message, "FK 참조 조회") { t.Errorf("slice result should not require @empty: %s", e.Message) }
	}
}
