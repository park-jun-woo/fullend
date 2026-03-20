//ff:func feature=stml-validate type=test control=sequence
//ff:what custom.ts fallback으로 bind 검증 통과 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateCustomTSFallback(t *testing.T) {
	customFiles := map[string]string{"test-page.custom.ts": `export function totalPrice(items) { return items.reduce((sum, item) => sum + item.price, 0) }`}
	root := setupTestProject(t, dummyOpenAPI, customFiles, nil)
	pages := []parser.PageSpec{{Name: "test-page", FileName: "test-page.html", Fetches: []parser.FetchBlock{{OperationID: "GetReservation", Binds: []parser.FieldBind{{Name: "totalPrice", Tag: "span"}}}}}}
	errs := Validate(pages, root)
	if len(errs) > 0 { for _, e := range errs { t.Error(e.Error()) }; t.Fatal("expected no errors with custom.ts fallback") }
}
