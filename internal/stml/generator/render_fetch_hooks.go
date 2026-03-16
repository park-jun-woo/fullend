//ff:func feature=stml-gen type=generator control=iteration dimension=1
//ff:what FetchBlock의 useState + useQuery 훅 선언을 렌더링한다
package generator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/stml/parser"
)

// renderFetchHooks writes useState + useQuery hook declarations.
func renderFetchHooks(f parser.FetchBlock, sb *strings.Builder) {
	if f.Paginate {
		defaultLimit := 20
		sb.WriteString(fmt.Sprintf("  const [page, setPage] = useState(1)\n"))
		sb.WriteString(fmt.Sprintf("  const [limit] = useState(%d)\n", defaultLimit))
	}
	if f.Sort != nil {
		sb.WriteString(fmt.Sprintf("  const [sortBy, setSortBy] = useState('%s')\n", f.Sort.Column))
		sb.WriteString(fmt.Sprintf("  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('%s')\n", f.Sort.Direction))
	}
	if len(f.Filters) > 0 {
		sb.WriteString("  const [filters, setFilters] = useState<Record<string, string>>({})\n")
	}
	if f.Paginate || f.Sort != nil || len(f.Filters) > 0 {
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("  %s\n\n", renderUseQuery(f)))
	for _, child := range f.NestedFetches {
		renderFetchHooks(child, sb)
	}
}
