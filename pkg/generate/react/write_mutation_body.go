//ff:func feature=gen-react type=generator control=sequence
//ff:what POST/PUT/DELETE 엔드포인트의 fetch 함수 본문을 생성한다

package react

import (
	"fmt"
	"strings"
)

// writeMutationBody writes the body of a POST/PUT/DELETE fetch function.
func writeMutationBody(b *strings.Builder, ep endpoint, method string, fetchPath string) {
	if len(ep.pathParams) > 0 {
		excluded := make([]string, len(ep.pathParams))
		for i, pp := range ep.pathParams {
			excluded[i] = fmt.Sprintf("'%s'", pp)
		}
		b.WriteString(fmt.Sprintf("  const exclude = new Set([%s])\n", strings.Join(excluded, ", ")))
		b.WriteString("  const body: Record<string, any> = {}\n")
		b.WriteString("  if (params) {\n")
		b.WriteString("    for (const [k, v] of Object.entries(params)) {\n")
		b.WriteString("      if (!exclude.has(k)) body[k] = v\n")
		b.WriteString("    }\n")
		b.WriteString("  }\n")
	} else {
		b.WriteString("  const body = params ?? {}\n")
	}
	b.WriteString(fmt.Sprintf("  const res = await fetch(`${BASE}%s`, {\n", fetchPath))
	b.WriteString(fmt.Sprintf("    method: '%s',\n", method))
	b.WriteString("    headers: { 'Content-Type': 'application/json' },\n")
	b.WriteString("    body: JSON.stringify(body),\n")
	b.WriteString("  })\n")
	b.WriteString("  return res.json()\n")
}
