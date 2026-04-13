//ff:func feature=ssac-gen type=test control=sequence
//ff:what @put 시퀀스의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGeneratePut(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "UpdateCourse", FileName: "update_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPut, Model: "Course.Update", Inputs: map[string]string{"Title": "request.Title", "ID": "course.ID"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `err = h.CourseModel.WithTx(tx).Update(course.ID, title)`)
	assertContains(t, code, `h.DB.BeginTx`)
}
