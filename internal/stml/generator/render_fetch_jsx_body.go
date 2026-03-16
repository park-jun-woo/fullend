//ff:func feature=stml-gen type=generator control=sequence
//ff:what FetchBlock의 필터/정렬/자식/페이지네이션 UI를 생성한다
package generator

import "github.com/geul-org/fullend/internal/stml/parser"

// renderFetchJSXBody generates the inner content of a fetch JSX block.
func renderFetchJSXBody(f parser.FetchBlock, alias string, indent int) []string {
	var lines []string

	// Phase 5: filter UI
	if len(f.Filters) > 0 {
		lines = append(lines, renderFilterUI(f.Filters, indent)...)
	}

	// Phase 5: sort UI
	if f.Sort != nil {
		lines = append(lines, renderSortUI(f.Sort, indent)...)
	}

	if len(f.Children) > 0 {
		lines = append(lines, renderChildNodes(f.Children, alias, "item", indent)...)
	} else {
		lines = append(lines, renderFetchJSXFlatChildren(f, alias, indent)...)
	}

	// Phase 5: pagination UI
	if f.Paginate {
		lines = append(lines, renderPaginationUI(alias, indent)...)
	}

	return lines
}
