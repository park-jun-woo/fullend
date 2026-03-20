//ff:func feature=ssac-validate type=test control=sequence
//ff:what FK 참조 조회 후 @empty 누락 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateFKReferenceGetMissingEmpty(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "UpdateSchedule", FileName: "update_schedule.go", Sequences: []parser.Sequence{
		{Type: parser.SeqGet, Model: "Schedule.FindByID", Inputs: map[string]string{"ID": "request.ScheduleID"}, Result: &parser.Result{Type: "Schedule", Var: "schedule"}},
		{Type: parser.SeqGet, Model: "Project.FindByID", Inputs: map[string]string{"ID": "schedule.ProjectID"}, Result: &parser.Result{Type: "Project", Var: "project"}},
		{Type: parser.SeqPut, Model: "Schedule.Update", Inputs: map[string]string{"ID": "request.ScheduleID", "Name": "request.Name"}},
	}}}
	errs := Validate(funcs)
	assertHasError(t, errs, "FK 참조 조회 후 @empty 가드가 필요합니다")
}
