//ff:func feature=gen-react type=generator control=sequence
//ff:what GET 엔드포인트의 fetch 함수 본문을 생성한다

package react

import (
	"fmt"
	"strings"
)

// writeGetBody writes the body of a GET fetch function.
func writeGetBody(b *strings.Builder, ep endpoint, fetchPath string) {
	b.WriteString("  const query = new URLSearchParams()\n")
	b.WriteString("  if (params) {\n")
	if len(ep.pathParams) > 0 {
		excluded := make([]string, len(ep.pathParams))
		for i, pp := range ep.pathParams {
			excluded[i] = fmt.Sprintf("'%s'", pp)
		}
		b.WriteString(fmt.Sprintf("    const exclude = new Set([%s])\n", strings.Join(excluded, ", ")))
		b.WriteString("    for (const [k, v] of Object.entries(params)) {\n")
		b.WriteString("      if (v != null && !exclude.has(k)) query.set(k, String(v))\n")
		b.WriteString("    }\n")
	} else {
		b.WriteString("    for (const [k, v] of Object.entries(params)) {\n")
		b.WriteString("      if (v != null) query.set(k, String(v))\n")
		b.WriteString("    }\n")
	}
	b.WriteString("  }\n")
	b.WriteString("  const qs = query.toString()\n")
	b.WriteString(fmt.Sprintf("  const res = await fetch(`${BASE}%s${qs ? '?' + qs : ''}`)\n", fetchPath))
	b.WriteString("  return res.json()\n")
}
