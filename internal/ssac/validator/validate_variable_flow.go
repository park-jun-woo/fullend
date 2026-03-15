//ff:func feature=ssac-validate type=rule
//ff:what 변수 선언 전 참조 검증

package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateVariableFlow는 변수가 선언 전에 참조되지 않는지 검증한다.
func validateVariableFlow(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
	declared := map[string]bool{
		"currentUser": true,
	}
	if sf.Subscribe != nil {
		declared["message"] = true
	}

	for i, seq := range sf.Sequences {
		ctx := errCtx{sf.FileName, sf.Name, i}

		// guard Target 검증
		if seq.Target != "" {
			rootTarget := rootVar(seq.Target)
			if !declared[rootTarget] {
				errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%q 변수가 아직 선언되지 않았습니다", rootTarget)))
			}
		}

		// Args source 검증
		for _, arg := range seq.Args {
			ref := argVarRef(arg)
			if ref != "" && !declared[ref] {
				errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%q 변수가 아직 선언되지 않았습니다", ref)))
			}
		}

		// Inputs value 검증
		for _, val := range seq.Inputs {
			if strings.HasPrefix(val, `"`) {
				continue // 리터럴
			}
			if strings.HasPrefix(val, "config.") {
				errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("config.%s — config.* 입력은 지원하지 않습니다. func 내부에서 os.Getenv()를 사용하세요", val[len("config."):])))
				continue
			}
			// @publish에서 query 사용 금지
			if seq.Type == parser.SeqPublish && (val == "query" || strings.HasPrefix(val, "query.")) {
				errs = append(errs, ctx.err("@publish", "query는 HTTP 전용입니다 — @publish에서 사용할 수 없습니다"))
				continue
			}
			ref := rootVar(val)
			if ref != "request" && ref != "currentUser" && ref != "query" && ref != "" && !declared[ref] {
				errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%q 변수가 아직 선언되지 않았습니다", ref)))
			}
		}

		// @response Fields value 검증
		for _, val := range seq.Fields {
			if strings.HasPrefix(val, `"`) {
				continue // 리터럴
			}
			ref := rootVar(val)
			if ref != "" && !declared[ref] {
				errs = append(errs, ctx.err("@response", fmt.Sprintf("%q 변수가 아직 선언되지 않았습니다", ref)))
			}
		}

		// Result로 변수 선언
		if seq.Result != nil {
			declared[seq.Result.Var] = true
		}
	}

	return errs
}
