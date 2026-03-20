//ff:func feature=stml-validate type=test control=sequence
//ff:what x-sort.allowed에 없는 컬럼 사용 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateSortNotAllowed(t *testing.T) {
	root := setupTestProject(t, infraOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "ListItems", Sort: &parser.SortDecl{Column: "invalid_col", Direction: "asc"}, Eaches: []parser.EachBlock{{Field: "items"}}}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "x-sort.allowed")
	assertHasError(t, errs, "invalid_col")
}
