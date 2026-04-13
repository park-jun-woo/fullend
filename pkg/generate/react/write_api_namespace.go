//ff:func feature=gen-react type=generator control=iteration dimension=1
//ff:what export const api 객체를 생성한다

package react

import (
	"fmt"
	"strings"
)

// writeApiNamespace writes the export const api = { ... } object.
func writeApiNamespace(b *strings.Builder, endpoints []endpoint) {
	b.WriteString("export const api = {\n")
	for i, ep := range endpoints {
		funcName := lcFirst(ep.opID)
		b.WriteString(fmt.Sprintf("  %s: %s", ep.opID, funcName))
		if i < len(endpoints)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}
	b.WriteString("}\n")
}
