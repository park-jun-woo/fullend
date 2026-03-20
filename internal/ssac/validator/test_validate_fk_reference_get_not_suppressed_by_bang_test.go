//ff:func feature=ssac-validate type=test control=sequence
//ff:what SuppressWarn으로도 FK 참조 ERROR는 억제 불가 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateFKReferenceGetNotSuppressedByBang(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "UpdateSchedule", FileName: "update_schedule.go", Sequences: []parser.Sequence{
		{Type: parser.SeqGet, Model: "Schedule.FindByID", Inputs: map[string]string{"ID": "request.ScheduleID"}, Result: &parser.Result{Type: "Schedule", Var: "schedule"}},
		{Type: parser.SeqGet, Model: "Project.FindByID", Inputs: map[string]string{"ID": "schedule.ProjectID"}, Result: &parser.Result{Type: "Project", Var: "project"}, SuppressWarn: true},
	}}}
	errs := Validate(funcs)
	assertHasError(t, errs, "FK 참조 조회 후 @empty 가드가 필요합니다")
}
