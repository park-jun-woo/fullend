//ff:func feature=stml-validate type=test control=sequence
//ff:what x-sort 없는 endpoint에 data-sort 사용 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateSortNoExtOnEndpoint(t *testing.T) {
	root := setupTestProject(t, infraOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "ListSimple", Sort: &parser.SortDecl{Column: "name", Direction: "asc"}, Eaches: []parser.EachBlock{{Field: "items"}}}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "x-sort.allowed")
}
