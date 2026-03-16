//ff:func feature=ssac-gen type=generator control=sequence
//ff:what QueryOpts의 Pagination 설정 코드를 버퍼에 출력
package generator

import (
	"bytes"
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/validator"
)

func writePaginationConfig(buf *bytes.Buffer, op validator.OperationSymbol) {
	if op.XPagination == nil {
		return
	}
	fmt.Fprintf(buf, "\t\tPagination: &model.PaginationConfig{Style: %q, DefaultLimit: %d, MaxLimit: %d},\n",
		op.XPagination.Style, op.XPagination.DefaultLimit, op.XPagination.MaxLimit)
}
