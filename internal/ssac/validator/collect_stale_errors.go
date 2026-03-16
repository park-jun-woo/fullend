//ff:func feature=ssac-validate type=util control=iteration dimension=1
//ff:what response 필드에서 stale 변수 참조를 검출한다
package validator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func collectStaleErrors(seq parser.Sequence, getVars map[string]string, mutated map[string]bool, ctx errCtx, errs []ValidationError) []ValidationError {
	for field, val := range seq.Fields {
		ref := rootVar(val)
		modelName, ok := getVars[ref]
		if !ok || !mutated[modelName] {
			continue
		}
		errs = append(errs, ValidationError{
			FileName: ctx.fileName,
			FuncName: ctx.funcName,
			SeqIndex: ctx.seqIndex,
			Tag:      "@response",
			Message:  fmt.Sprintf("%q (필드 %q)가 %s 수정 이후 갱신 없이 response에 사용됩니다", ref, field, modelName),
			Level:    "WARNING",
		})
	}
	return errs
}
