//ff:func feature=stml-gen type=generator control=iteration dimension=1
//ff:what 페이지의 useParams, useQueryClient, useQuery 훅을 렌더링한다
package generator

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/internal/stml/parser"
)

func renderPageHooks(page parser.PageSpec, is importSet, sb *strings.Builder) {
	allParams := collectAllParams(page)
	if up := renderUseParams(allParams); up != "" {
		sb.WriteString(fmt.Sprintf("  %s\n", up))
	}

	if is.useQueryClient {
		sb.WriteString("  const queryClient = useQueryClient()\n")
	}

	sb.WriteString("\n")

	for _, f := range page.Fetches {
		renderFetchHooks(f, sb)
	}
}
