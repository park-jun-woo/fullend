//ff:func feature=ssac-validate type=rule
//ff:what x-pagination style과 Result.Wrapper 타입 일치 검증

package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validatePaginationType은 x-pagination style과 Result.Wrapper 타입의 일치를 검증한다.
func validatePaginationType(sf parser.ServiceFunc, st *SymbolTable) []ValidationError {
	if st == nil {
		return nil
	}

	op, ok := st.Operations[sf.Name]
	if !ok {
		return nil
	}

	var errs []ValidationError
	for i, seq := range sf.Sequences {
		if seq.Result == nil || seq.Result.Wrapper == "" && !strings.HasPrefix(seq.Result.Type, "[]") {
			continue
		}
		ctx := errCtx{sf.FileName, sf.Name, i}

		if op.XPagination != nil {
			// x-pagination 있음 → Wrapper 필수
			switch op.XPagination.Style {
			case "offset":
				if seq.Result.Wrapper != "Page" {
					errs = append(errs, ctx.err("@get", "x-pagination style: offset이지만 반환 타입이 Page[T]가 아닙니다"))
				}
			case "cursor":
				if seq.Result.Wrapper != "Cursor" {
					errs = append(errs, ctx.err("@get", "x-pagination style: cursor이지만 반환 타입이 Cursor[T]가 아닙니다"))
				}
			}
		} else {
			// x-pagination 없음 → Wrapper 사용 불가
			if seq.Result.Wrapper != "" {
				errs = append(errs, ctx.err("@get", fmt.Sprintf("OpenAPI에 x-pagination이 없지만 %s[T] 타입을 사용했습니다. []T를 사용하세요", seq.Result.Wrapper)))
			}
		}
	}

	return errs
}
