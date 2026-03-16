//ff:func feature=stml-gen type=generator control=sequence
//ff:what 페이지네이션 이전/다음 컨트롤 JSX를 생성한다
package generator

import "fmt"

// renderPaginationUI generates pagination controls.
func renderPaginationUI(alias string, indent int) []string {
	ind := indentStr(indent)
	var lines []string
	lines = append(lines, fmt.Sprintf(`%s<div className="flex justify-between items-center mt-4">`, ind))
	lines = append(lines, fmt.Sprintf(`%s  <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>이전</button>`, ind))
	lines = append(lines, fmt.Sprintf(`%s  <span>{page} / {Math.ceil((%s?.total ?? 0) / limit)}</span>`, ind, alias))
	lines = append(lines, fmt.Sprintf(`%s  <button disabled={!%s?.total || page * limit >= %s.total} onClick={() => setPage(p => p + 1)}>다음</button>`, ind, alias, alias))
	lines = append(lines, fmt.Sprintf(`%s</div>`, ind))
	return lines
}
