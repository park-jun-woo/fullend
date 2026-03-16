//ff:func feature=stml-gen type=generator control=iteration dimension=1
//ff:what 필터 입력 컨트롤 JSX를 생성한다
package generator

import "fmt"

// renderFilterUI generates filter input controls.
func renderFilterUI(filters []string, indent int) []string {
	ind := indentStr(indent)
	var lines []string
	lines = append(lines, fmt.Sprintf(`%s<div className="flex gap-2 mb-4">`, ind))
	for _, col := range filters {
		lines = append(lines, fmt.Sprintf(`%s  <input placeholder="%s" value={filters.%s ?? ''} className="px-3 py-2 border rounded" onChange={(e) => setFilters(f => ({ ...f, %s: e.target.value }))} />`, ind, col, col, col))
	}
	lines = append(lines, fmt.Sprintf(`%s</div>`, ind))
	return lines
}
