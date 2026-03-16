//ff:func feature=ssac-validate type=util control=iteration dimension=1
//ff:what Inputs에서 request 접두사 필드의 OpenAPI 일치를 검사한다
package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func collectInputRequestErrors(seq parser.Sequence, op OperationSymbol, ctx errCtx, used map[string]bool, errs []ValidationError) []ValidationError {
	for _, val := range seq.Inputs {
		if !strings.HasPrefix(val, "request.") {
			continue
		}
		field := val[len("request."):]
		used[field] = true
		if !op.RequestFields[field] {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("OpenAPI request에 %q 필드가 없습니다", field)))
		}
	}
	return errs
}
