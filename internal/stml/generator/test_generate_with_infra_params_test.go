//ff:func feature=stml-gen type=test control=sequence
//ff:what infra params(paginate, sort, filter) TSX 생성을 검증
package generator

import ("strings"; "testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestGenerateWithInfraParams(t *testing.T) {
	page, _ := parser.ParseReader("list-page.html", strings.NewReader(`<main>
  <section data-fetch="ListItems" data-paginate data-sort="name:desc" data-filter="status,category">
    <ul data-each="items">
      <li><span data-bind="name"></span></li>
    </ul>
  </section>
</main>`))
	code := GeneratePage(page, "")
	assertContains(t, code, "import { useState } from 'react'")
	assertContains(t, code, "const [page, setPage] = useState(1)")
	assertContains(t, code, "const [limit] = useState(20)")
	assertContains(t, code, "const [sortBy, setSortBy] = useState('name')")
	assertContains(t, code, "const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc')")
	assertContains(t, code, "const [filters, setFilters] = useState<Record<string, string>>({})")
	assertContains(t, code, "page, limit")
	assertContains(t, code, "sortBy, sortDir")
	assertContains(t, code, "filters")
	assertContains(t, code, `placeholder="status"`)
	assertContains(t, code, `placeholder="category"`)
	assertContains(t, code, "setFilters")
	assertContains(t, code, "setSortBy")
	assertContains(t, code, "setSortDir")
	assertContains(t, code, "setPage")
	assertContains(t, code, ">이전</button>")
	assertContains(t, code, ">다음</button>")
}
