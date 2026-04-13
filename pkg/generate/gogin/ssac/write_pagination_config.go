//ff:func feature=ssac-gen type=generator control=sequence topic=query-opts
//ff:what QueryOpts의 Pagination 설정 코드를 버퍼에 출력
package ssac

import (
	"bytes"
	"fmt"

	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func writePaginationConfig(buf *bytes.Buffer, op validator.OperationSymbol) {
	if op.XPagination == nil {
		return
	}
	fmt.Fprintf(buf, "\t\tPagination: &model.PaginationConfig{Style: %q, DefaultLimit: %d, MaxLimit: %d},\n",
		op.XPagination.Style, op.XPagination.DefaultLimit, op.XPagination.MaxLimit)
}
