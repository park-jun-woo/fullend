//ff:func feature=stml-validate type=test control=sequence
//ff:what infra params 검증 통과 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateInfraParamsPass(t *testing.T) {
	root := setupTestProject(t, infraOpenAPI, nil, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "ListItems", Paginate: true, Sort: &parser.SortDecl{Column: "name", Direction: "desc"}, Filters: []string{"status"}, Eaches: []parser.EachBlock{{Field: "items"}}}}}}
	errs := Validate(pages, root)
	if len(errs) > 0 { for _, e := range errs { t.Error(e.Error()) } }
}
