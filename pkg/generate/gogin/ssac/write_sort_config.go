//ff:func feature=ssac-gen type=generator control=sequence topic=query-opts
//ff:what QueryOpts의 Sort 설정 코드를 버퍼에 출력
package ssac

import (
	"bytes"
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func writeSortConfig(buf *bytes.Buffer, op rule.OperationInfo) {
	if op.Sort == nil {
		return
	}
	buf.WriteString("\t\tSort: &model.SortConfig{")
	writeSortAllowed(buf, op.Sort.Allowed)
	if op.Sort.Default != "" {
		fmt.Fprintf(buf, ", Default: %q", op.Sort.Default)
	}
	if op.Sort.Direction != "" {
		fmt.Fprintf(buf, ", Direction: %q", op.Sort.Direction)
	}
	buf.WriteString("},\n")
}
