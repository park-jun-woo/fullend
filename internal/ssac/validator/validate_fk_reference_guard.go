//ff:func feature=ssac-validate type=rule control=iteration dimension=2 topic=type-resolve
//ff:what FK 참조 @get 후 @empty 가드 누락 검증

package validator

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// validateFKReferenceGuard는 FK 참조 @get 후 @empty 가드 누락을 검증한다.
// FK 참조: @get의 input이 이전 result 변수의 필드를 참조 (request/currentUser 아닌 경우).
// nil pointer dereference 방지를 위해 @empty 가드가 필요하다. @get!로 억제 가능.
func validateFKReferenceGuard(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
	declared := map[string]bool{}
	varTypes := map[string]string{} // var name → Model type
	if sf.Subscribe != nil {
		declared["message"] = true
	}

	for i, seq := range sf.Sequences {
		if seq.Type != parser.SeqGet || seq.Result == nil {
			if seq.Result != nil {
				declared[seq.Result.Var] = true
				varTypes[seq.Result.Var] = seq.Result.Type
			}
			continue
		}

		// 슬라이스/래퍼 결과는 nil dereference 위험 없음
		if strings.HasPrefix(seq.Result.Type, "[]") || seq.Result.Wrapper != "" {
			declared[seq.Result.Var] = true
			varTypes[seq.Result.Var] = seq.Result.Type
			continue
		}

		// input 중 이전 result 변수 참조가 있는지 확인
		// 같은 Model의 재쿼리(PK lookup)는 FK 참조가 아님
		hasFKRef := hasForeignKeyRef(seq, declared, varTypes)

		if seq.Result != nil {
			declared[seq.Result.Var] = true
			varTypes[seq.Result.Var] = seq.Result.Type
		}
		if !hasFKRef {
			continue
		}
		// 이후 시퀀스에 @empty 가드가 있는지 확인
		if hasEmptyGuardFor(sf.Sequences[i+1:], seq.Result.Var) {
			continue
		}
		ctx := errCtx{sf.FileName, sf.Name, i}
		errs = append(errs, ctx.err("@get", fmt.Sprintf("%q — FK 참조 조회 후 @empty 가드가 필요합니다", seq.Result.Var)))
	}

	return errs
}
