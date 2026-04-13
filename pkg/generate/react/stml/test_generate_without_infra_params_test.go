//ff:func feature=stml-gen type=test control=sequence
//ff:what infra params 없는 페이지에서 관련 코드 미생성 검증
package stml

import ("strings"; "testing"; stmlparser "github.com/park-jun-woo/fullend/pkg/parser/stml")

func TestGenerateWithoutInfraParams(t *testing.T) {
	page, _ := stmlparser.ParseReader("simple-page.html", strings.NewReader(`<section data-fetch="GetItem">
  <span data-bind="name"></span>
</section>`))
	code := GeneratePage(page, "")
	assertNotContains(t, code, "useState")
	assertNotContains(t, code, "setPage")
	assertNotContains(t, code, "setSortBy")
	assertNotContains(t, code, "setFilters")
	assertNotContains(t, code, "이전")
	assertNotContains(t, code, "다음")
}
