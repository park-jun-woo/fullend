//ff:func feature=stml-validate type=test control=sequence
//ff:what x-pagination 없는 endpoint에 data-paginate 사용 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidatePaginateNoExt(t *testing.T) {
	root := setupTestProject(t, infraOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "ListSimple", Paginate: true, Eaches: []parser.EachBlock{{Field: "items"}}}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "x-pagination이 선언되지 않았습니다")
}
