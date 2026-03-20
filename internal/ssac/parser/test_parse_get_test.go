//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseGet: @get 어노테이션 파싱 후 모델·결과·입력 검증
package parser

import "testing"

func TestParseGet(t *testing.T) {
	src := `package service

// @get Course course = Course.FindByID({CourseID: request.CourseID})
func GetCourse(c *gin.Context) {}
`
	sfs := parseTestFile(t, src)
	if len(sfs) != 1 {
		t.Fatalf("expected 1 func, got %d", len(sfs))
	}
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Type", seq.Type, SeqGet)
	assertEqual(t, "Model", seq.Model, "Course.FindByID")
	if seq.Result == nil {
		t.Fatal("expected result")
	}
	assertEqual(t, "Result.Type", seq.Result.Type, "Course")
	assertEqual(t, "Result.Var", seq.Result.Var, "course")
	if len(seq.Inputs) != 1 {
		t.Fatalf("expected 1 input, got %d", len(seq.Inputs))
	}
	assertEqual(t, "Inputs.CourseID", seq.Inputs["CourseID"], "request.CourseID")
}
