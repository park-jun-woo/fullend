//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what FK 참조 조회 후 @empty가 있으면 에러 없음 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateFKReferenceGetWithEmpty(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "UpdateSchedule", FileName: "update_schedule.go", Sequences: []parser.Sequence{
		{Type: parser.SeqGet, Model: "Schedule.FindByID", Inputs: map[string]string{"ID": "request.ScheduleID"}, Result: &parser.Result{Type: "Schedule", Var: "schedule"}},
		{Type: parser.SeqGet, Model: "Project.FindByID", Inputs: map[string]string{"ID": "schedule.ProjectID"}, Result: &parser.Result{Type: "Project", Var: "project"}},
		{Type: parser.SeqEmpty, Target: "project", Message: "project not found"},
		{Type: parser.SeqPut, Model: "Schedule.Update", Inputs: map[string]string{"ID": "request.ScheduleID"}},
	}}}
	errs := Validate(funcs)
	for _, e := range errs {
		if contains(e.Message, "FK 참조 조회") { t.Errorf("unexpected FK reference error: %s", e.Message) }
	}
}
