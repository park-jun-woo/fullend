//ff:func feature=stml-validate type=test control=sequence
//ff:what x-filter.allowed에 없는 필터 사용 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateFilterNotAllowed(t *testing.T) {
	root := setupTestProject(t, infraOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "ListItems", Filters: []string{"status", "bad_col"}, Eaches: []parser.EachBlock{{Field: "items"}}}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "x-filter.allowed")
	assertHasError(t, errs, "bad_col")
}
