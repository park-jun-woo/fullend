//ff:func feature=stml-validate type=test control=sequence
//ff:what 존재하지 않는 operationId 참조 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateOperationNotFound(t *testing.T) {
	root := setupTestProject(t, dummyOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "NonExistent"}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "NonExistent")
	assertHasError(t, errs, "operationId가 없습니다")
}
