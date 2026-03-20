//ff:func feature=stml-validate type=test control=sequence
//ff:what 존재하지 않는 action field 참조 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateFieldNotFound(t *testing.T) {
	root := setupTestProject(t, dummyOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Actions: []parser.ActionBlock{{OperationID: "Login", Fields: []parser.FieldBind{{Name: "NonExistentField", Tag: "input"}}}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "NonExistentField")
	assertHasError(t, errs, "request schema")
}
