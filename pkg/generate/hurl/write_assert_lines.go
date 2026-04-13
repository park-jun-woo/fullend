//ff:func feature=gen-hurl type=generator control=iteration dimension=1
//ff:what writes [Asserts] block with assertion lines

package hurl

import "strings"

// writeAssertLines writes [Asserts] block if there are any assertions.
func writeAssertLines(buf *strings.Builder, asserts []string) {
	if len(asserts) == 0 {
		return
	}
	buf.WriteString("[Asserts]\n")
	for _, a := range asserts {
		buf.WriteString(a + "\n")
	}
}
