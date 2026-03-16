//ff:func feature=stml-gen type=util control=sequence
//ff:what 인프라 파라미터(페이지네이션, 정렬, 필터)를 포함한 API 호출 인자를 생성한다
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/stml/parser"
)

// renderInfraApiArgs builds the API call arguments with infra params.
func renderInfraApiArgs(f parser.FetchBlock, paramArgs string) string {
	var parts []string

	// Spread existing params
	if paramArgs != "" {
		// paramArgs is "{ key: val, ... }" — strip braces and spread
		inner := strings.TrimPrefix(paramArgs, "{ ")
		inner = strings.TrimSuffix(inner, " }")
		parts = append(parts, inner)
	}

	if f.Paginate {
		parts = append(parts, "page", "limit")
	}
	if f.Sort != nil {
		parts = append(parts, "sortBy", "sortDir")
	}
	if len(f.Filters) > 0 {
		parts = append(parts, "...filters")
	}

	return "{ " + strings.Join(parts, ", ") + " }"
}
