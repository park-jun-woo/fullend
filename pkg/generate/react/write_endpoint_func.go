//ff:func feature=gen-react type=generator control=iteration dimension=1
//ff:what 단일 엔드포인트의 async fetch 함수를 생성한다

package react

import (
	"fmt"
	"strings"
)

// writeEndpointFunc writes a single async function for an endpoint.
func writeEndpointFunc(b *strings.Builder, ep endpoint) {
	funcName := lcFirst(ep.opID)
	method := strings.ToUpper(ep.method)

	b.WriteString(fmt.Sprintf("async function %s(params?: Record<string, any>) {\n", funcName))

	// Build URL with path param substitution.
	fetchPath := openAPIPathToTemplateLiteral(ep.path)
	for _, pp := range ep.pathParams {
		b.WriteString(fmt.Sprintf("  const %s = params?.%s\n", pp, pp))
	}

	if method == "GET" {
		writeGetBody(b, ep, fetchPath)
	} else {
		writeMutationBody(b, ep, method, fetchPath)
	}
	b.WriteString("}\n\n")
}
