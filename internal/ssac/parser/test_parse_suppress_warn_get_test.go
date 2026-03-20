//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseSuppressWarnGet: @get! SuppressWarn 플래그 파싱 검증
package parser

import "testing"

func TestParseSuppressWarnGet(t *testing.T) {
	src := `package service

// @get! Course course = Course.FindByID({CourseID: request.CourseID})
func GetCourse() {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Type", seq.Type, SeqGet)
	assertEqual(t, "Model", seq.Model, "Course.FindByID")
	if !seq.SuppressWarn {
		t.Error("expected SuppressWarn=true")
	}
}
