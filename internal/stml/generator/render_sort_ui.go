//ff:func feature=stml-gen type=generator control=sequence
//ff:what 정렬 토글 컨트롤 JSX를 생성한다
package generator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/stml/parser"
)

// renderSortUI generates sort toggle controls.
func renderSortUI(sort *parser.SortDecl, indent int) []string {
	ind := indentStr(indent)
	var lines []string
	lines = append(lines, fmt.Sprintf(`%s<div className="flex gap-2 mb-4">`, ind))
	lines = append(lines, fmt.Sprintf(`%s  <button onClick={() => { setSortBy('%s'); setSortDir(d => d === 'asc' ? 'desc' : 'asc') }}>`, ind, sort.Column))
	lines = append(lines, fmt.Sprintf(`%s    %s {sortBy === '%s' ? (sortDir === 'asc' ? '↑' : '↓') : ''}`, ind, sort.Column, sort.Column))
	lines = append(lines, fmt.Sprintf(`%s  </button>`, ind))
	lines = append(lines, fmt.Sprintf(`%s</div>`, ind))
	return lines
}
