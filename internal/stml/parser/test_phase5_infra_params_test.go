//ff:func feature=stml-parse type=test control=sequence
//ff:what Phase5 infra params(paginate, sort, filter) 파싱 검증
package parser

import ("strings"; "testing")

func TestPhase5_InfraParams(t *testing.T) {
	input := `<main>
  <section data-fetch="ListMyReservations"
           data-paginate
           data-sort="StartAt:desc"
           data-filter="Status,RoomID"
    <ul data-each="reservations">
      <li><span data-bind="RoomID"></span></li>
    </ul>
  </section>
</main>`
	page, err := ParseReader("test.html", strings.NewReader(input))
	if err != nil { t.Fatal(err) }
	fetch := page.Fetches[0]
	if !fetch.Paginate { t.Error("Paginate = false, want true") }
	if fetch.Sort == nil { t.Fatal("Sort = nil") }
	if fetch.Sort.Column != "StartAt" { t.Errorf("Sort.Column = %q", fetch.Sort.Column) }
	if fetch.Sort.Direction != "desc" { t.Errorf("Sort.Direction = %q", fetch.Sort.Direction) }
	if len(fetch.Filters) != 2 { t.Fatalf("Filters = %d", len(fetch.Filters)) }
	if fetch.Filters[0] != "Status" || fetch.Filters[1] != "RoomID" { t.Errorf("Filters = %v", fetch.Filters) }
}
