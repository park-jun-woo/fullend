//ff:func feature=ssac-validate type=rule control=iteration dimension=1 topic=query-opts
//ff:what x-pagination style과 Result.Wrapper 타입 일치 검증
package validator

import (
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
		errs = append(errs, checkPaginationStyle(op, seq.Result.Wrapper, ctx)...)
	}

	return errs
}
