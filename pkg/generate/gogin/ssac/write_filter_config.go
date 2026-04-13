//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=query-opts
//ff:what QueryOpts의 Filter 설정 코드를 버퍼에 출력
package ssac

import (
	"bytes"
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func writeFilterConfig(buf *bytes.Buffer, op rule.OperationInfo) {
	if op.Filter == nil || len(op.Filter.Allowed) == 0 {
		return
	}
	buf.WriteString("\t\tFilter: &model.FilterConfig{Allowed: []string{")
	for i, col := range op.Filter.Allowed {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(buf, "%q", col)
	}
	buf.WriteString("}},\n")
}
