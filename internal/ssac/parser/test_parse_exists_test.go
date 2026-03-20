//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseExists: @exists 가드 어노테이션 파싱 후 타겟·메시지 검증
package parser

import "testing"

func TestParseExists(t *testing.T) {
	src := `package service

// @exists existing "이미 존재합니다"
func CreateCourse(c *gin.Context) {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Type", seq.Type, SeqExists)
	assertEqual(t, "Target", seq.Target, "existing")
	assertEqual(t, "Message", seq.Message, "이미 존재합니다")
}
