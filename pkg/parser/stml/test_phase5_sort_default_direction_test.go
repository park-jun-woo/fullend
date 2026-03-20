//ff:func feature=stml-parse type=parser control=sequence
//ff:what TestPhase5_SortDefaultDirection — data-sort without direction defaults to asc

package parser

import (
	"strings"
	"testing"
)

func TestPhase5_SortDefaultDirection(t *testing.T) {
	input := `<section data-fetch="ListItems" data-sort="name">
  <ul data-each="items"><li><span data-bind="name"></span></li></ul>
</section>`

	page, err := ParseReader("test.html", strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	fetch := page.Fetches[0]
	if fetch.Sort == nil {
		t.Fatal("Sort = nil")
	}
	if fetch.Sort.Direction != "asc" {
		t.Errorf("Sort.Direction = %q, want %q", fetch.Sort.Direction, "asc")
	}
}
