//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=query-opts
//ff:what Sort 허용 컬럼 목록을 Go 코드로 버퍼에 출력
package generator

import (
	"bytes"
	"fmt"
)

func writeSortAllowed(buf *bytes.Buffer, allowed []string) {
	if len(allowed) == 0 {
		return
	}
	buf.WriteString("Allowed: []string{")
	for i, col := range allowed {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(buf, "%q", col)
	}
	buf.WriteString("}")
}
