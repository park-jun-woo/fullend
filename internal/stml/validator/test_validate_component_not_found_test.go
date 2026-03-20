//ff:func feature=stml-validate type=test control=sequence
//ff:what 존재하지 않는 컴포넌트 참조 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateComponentNotFound(t *testing.T) {
	root := setupTestProject(t, dummyOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "ListMyReservations", Components: []parser.ComponentRef{{Name: "MissingComponent"}}}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "MissingComponent")
	assertHasError(t, errs, "파일이 없습니다")
}
