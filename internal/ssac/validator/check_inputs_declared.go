//ff:func feature=ssac-validate type=util control=iteration dimension=1 topic=args-inputs
//ff:what Inputs 참조 변수의 선언 여부를 검사한다
package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func checkInputsDeclared(seq parser.Sequence, ctx errCtx, declared map[string]bool, errs []ValidationError) []ValidationError {
	for _, val := range seq.Inputs {
		if strings.HasPrefix(val, `"`) {
			continue
		}
		if strings.HasPrefix(val, "config.") {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("config.%s — config.* 입력은 지원하지 않습니다. func 내부에서 os.Getenv()를 사용하세요", val[len("config."):])))
			continue
		}
		if seq.Type == parser.SeqPublish && (val == "query" || strings.HasPrefix(val, "query.")) {
			errs = append(errs, ctx.err("@publish", "query는 HTTP 전용입니다 — @publish에서 사용할 수 없습니다"))
			continue
		}
		ref := rootVar(val)
		if ref != "request" && ref != "currentUser" && ref != "query" && ref != "" && !declared[ref] {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%q 변수가 아직 선언되지 않았습니다", ref)))
		}
	}
	return errs
}
