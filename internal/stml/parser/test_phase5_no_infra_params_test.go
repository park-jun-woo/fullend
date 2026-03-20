//ff:func feature=stml-parse type=test control=sequence
//ff:what Phase5 infra params 없을 때 기본값 검증
package parser

import ("strings"; "testing")

func TestPhase5_NoInfraParams(t *testing.T) {
	input := `<section data-fetch="GetItem">
  <span data-bind="name"></span>
</section>`
	page, err := ParseReader("test.html", strings.NewReader(input))
	if err != nil { t.Fatal(err) }
	fetch := page.Fetches[0]
	if fetch.Paginate { t.Error("Paginate = true") }
	if fetch.Sort != nil { t.Error("Sort != nil") }
	if len(fetch.Filters) != 0 { t.Errorf("Filters = %d", len(fetch.Filters)) }
}
