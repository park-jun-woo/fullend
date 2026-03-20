//ff:func feature=stml-parse type=parser control=sequence
//ff:what TestPhase5_NoInfraParams — fetch without paginate/sort/filter yields zero defaults

package parser

import (
	"strings"
	"testing"
)

func TestPhase5_NoInfraParams(t *testing.T) {
	input := `<section data-fetch="GetItem">
  <span data-bind="name"></span>
</section>`

	page, diags := ParseReader("test.html", strings.NewReader(input))
	if len(diags) > 0 {
		t.Fatal(diags)
	}

	fetch := page.Fetches[0]
	if fetch.Paginate {
		t.Error("Paginate = true, want false")
	}
	if fetch.Sort != nil {
		t.Error("Sort != nil, want nil")
	}
	if len(fetch.Filters) != 0 {
		t.Errorf("Filters = %d, want 0", len(fetch.Filters))
	}
}
