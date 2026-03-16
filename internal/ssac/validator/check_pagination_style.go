//ff:func feature=ssac-validate type=rule control=selection topic=query-opts
//ff:what x-pagination style과 Result.Wrapper 타입 일치를 검사한다
package validator

import (
	"fmt"
)

// checkPaginationStyle은 단일 시퀀스에 대해 x-pagination style과 Result.Wrapper 타입의 일치를 검사한다.
func checkPaginationStyle(op OperationSymbol, wrapper string, ctx errCtx) []ValidationError {
	var errs []ValidationError

	if op.XPagination == nil {
		// x-pagination 없음 → Wrapper 사용 불가
		if wrapper != "" {
			errs = append(errs, ctx.err("@get", fmt.Sprintf("OpenAPI에 x-pagination이 없지만 %s[T] 타입을 사용했습니다. []T를 사용하세요", wrapper)))
		}
		return errs
	}

	// x-pagination 있음 → Wrapper 필수
	switch op.XPagination.Style {
	case "offset":
		if wrapper != "Page" {
			errs = append(errs, ctx.err("@get", "x-pagination style: offset이지만 반환 타입이 Page[T]가 아닙니다"))
		}
	case "cursor":
		if wrapper != "Cursor" {
			errs = append(errs, ctx.err("@get", "x-pagination style: cursor이지만 반환 타입이 Cursor[T]가 아닙니다"))
		}
	}

	return errs
}
