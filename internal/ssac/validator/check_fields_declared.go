//ff:func feature=ssac-validate type=util control=iteration dimension=1 topic=args-inputs
//ff:what Fields 참조 변수의 선언 여부를 검사한다
package validator

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func checkFieldsDeclared(seq parser.Sequence, ctx errCtx, declared map[string]bool, errs []ValidationError) []ValidationError {
	for _, val := range seq.Fields {
		if strings.HasPrefix(val, `"`) {
			continue
		}
		ref := rootVar(val)
		if ref == "" || declared[ref] {
			continue
		}
		errs = append(errs, ctx.err("@response", fmt.Sprintf("%q 변수가 아직 선언되지 않았습니다", ref)))
	}
	return errs
}
