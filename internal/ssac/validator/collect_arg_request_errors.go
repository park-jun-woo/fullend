//ff:func feature=ssac-validate type=util control=iteration dimension=1
//ff:what Args에서 request 소스 필드의 OpenAPI 일치를 검사한다
package validator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func collectArgRequestErrors(seq parser.Sequence, op OperationSymbol, ctx errCtx, used map[string]bool, errs []ValidationError) []ValidationError {
	for _, arg := range seq.Args {
		if arg.Source != "request" {
			continue
		}
		used[arg.Field] = true
		if !op.RequestFields[arg.Field] {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("OpenAPI request에 %q 필드가 없습니다", arg.Field)))
		}
	}
	return errs
}
