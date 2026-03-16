//ff:func feature=ssac-validate type=util control=sequence
//ff:what target 변수의 선언 여부를 검사한다
package validator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func checkTargetDeclared(seq parser.Sequence, ctx errCtx, declared map[string]bool, errs []ValidationError) []ValidationError {
	if seq.Target == "" {
		return errs
	}
	rootTarget := rootVar(seq.Target)
	if declared[rootTarget] {
		return errs
	}
	return append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%q 변수가 아직 선언되지 않았습니다", rootTarget)))
}
