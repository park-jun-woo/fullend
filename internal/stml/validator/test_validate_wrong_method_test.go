//ff:func feature=stml-validate type=test control=sequence
//ff:what fetch에 POST 메서드 사용 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateWrongMethod(t *testing.T) {
	root := setupTestProject(t, dummyOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "Login"}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "POST")
	assertHasError(t, errs, "GET이어야 함")
}
