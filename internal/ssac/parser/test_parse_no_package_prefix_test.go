//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseNoPackagePrefix: 패키지 접두사 없는 모델에서 Package="" 검증
package parser

import "testing"

func TestParseNoPackagePrefix(t *testing.T) {
	src := `package service

// @get User user = User.FindByID({ID: request.ID})
func GetUser(c *gin.Context) {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Package", seq.Package, "")
	assertEqual(t, "Model", seq.Model, "User.FindByID")
}
