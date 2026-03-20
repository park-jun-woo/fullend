//ff:func feature=stml-validate type=test control=sequence
//ff:what data-each가 배열이 아닌 필드 참조 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateEachNotArray(t *testing.T) {
	root := setupTestProject(t, dummyOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "GetReservation", Eaches: []parser.EachBlock{{Field: "reservation"}}}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "배열이 아닙니다")
}
