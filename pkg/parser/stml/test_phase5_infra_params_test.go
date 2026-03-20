//ff:func feature=stml-parse type=parser control=sequence
//ff:what TestPhase5_InfraParams — data-paginate, data-sort, data-filter extraction

package parser

import (
	"strings"
	"testing"
)

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
	if err != nil {
		t.Fatal(err)
	}

	fetch := page.Fetches[0]

	// data-paginate
	if !fetch.Paginate {
		t.Error("Paginate = false, want true")
	}

	// data-sort
	if fetch.Sort == nil {
		t.Fatal("Sort = nil, want non-nil")
	}
	if fetch.Sort.Column != "StartAt" {
		t.Errorf("Sort.Column = %q, want %q", fetch.Sort.Column, "StartAt")
	}
	if fetch.Sort.Direction != "desc" {
		t.Errorf("Sort.Direction = %q, want %q", fetch.Sort.Direction, "desc")
	}

	// data-filter
	if len(fetch.Filters) != 2 {
		t.Fatalf("Filters = %d, want 2", len(fetch.Filters))
	}
	if fetch.Filters[0] != "Status" || fetch.Filters[1] != "RoomID" {
		t.Errorf("Filters = %v, want [Status RoomID]", fetch.Filters)
	}

}
