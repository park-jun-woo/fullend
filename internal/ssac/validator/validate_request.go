//ff:func feature=ssac-validate type=rule control=iteration dimension=1 topic=args-inputs
//ff:what request 필드 OpenAPI 일치 검증

package validator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateRequest는 request 필드가 OpenAPI와 일치하는지 검증한다.
func validateRequest(sf parser.ServiceFunc, st *SymbolTable) []ValidationError {
	var errs []ValidationError
	op, ok := st.Operations[sf.Name]
	if !ok {
		return nil
	}

	usedRequestFields := make(map[string]bool)
	for i, seq := range sf.Sequences {
		ctx := errCtx{sf.FileName, sf.Name, i}
		errs = collectArgRequestErrors(seq, op, ctx, usedRequestFields, errs)
		errs = collectInputRequestErrors(seq, op, ctx, usedRequestFields, errs)
	}

	// 역방향: OpenAPI → SSaC (path param 제외)
	pathParams := make(map[string]bool)
	for _, pp := range op.PathParams {
		pathParams[pp.Name] = true
	}
	for field := range op.RequestFields {
		if pathParams[field] {
			continue
		}
		if usedRequestFields[field] {
			continue
		}
		errs = append(errs, ValidationError{
			FileName: sf.FileName,
			FuncName: sf.Name,
			SeqIndex: -1,
			Tag:      "@request",
			Message:  fmt.Sprintf("OpenAPI request에 %q 필드가 있지만 SSaC에서 사용하지 않습니다", field),
			Level:    "WARNING",
		})
	}

	return errs
}
