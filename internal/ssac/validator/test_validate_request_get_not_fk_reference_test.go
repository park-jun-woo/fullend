//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what request.* 소스는 FK 참조가 아님 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateRequestGetNotFKReference(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "GetCourse", FileName: "get_course.go", Sequences: []parser.Sequence{
		{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.CourseID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
		{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
	}}}
	errs := Validate(funcs)
	for _, e := range errs {
		if contains(e.Message, "FK 참조 조회") { t.Errorf("request source should not trigger FK reference error: %s", e.Message) }
	}
}
