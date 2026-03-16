//ff:func feature=ssac-gen type=generator control=sequence topic=query-opts
//ff:what QueryOpts의 Sort 설정 코드를 버퍼에 출력
package generator

import (
	"bytes"
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/validator"
)

func writeSortConfig(buf *bytes.Buffer, op validator.OperationSymbol) {
	if op.XSort == nil {
		return
	}
	buf.WriteString("\t\tSort: &model.SortConfig{")
	writeSortAllowed(buf, op.XSort.Allowed)
	if op.XSort.Default != "" {
		fmt.Fprintf(buf, ", Default: %q", op.XSort.Default)
	}
	if op.XSort.Direction != "" {
		fmt.Fprintf(buf, ", Direction: %q", op.XSort.Direction)
	}
	buf.WriteString("},\n")
}
