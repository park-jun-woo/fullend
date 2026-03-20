//ff:func feature=stml-validate type=test control=sequence
//ff:what 존재하지 않는 parameter 참조 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateParamNotFound(t *testing.T) {
	root := setupTestProject(t, dummyOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "GetReservation", Params: []parser.ParamBind{{Name: "nonExistentParam", Source: "route.foo"}}}}}}
	errs := Validate(pages, root)
	assertHasError(t, errs, "nonExistentParam")
}
