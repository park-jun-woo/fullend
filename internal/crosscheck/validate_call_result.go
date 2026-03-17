//ff:func feature=crosscheck type=rule control=sequence topic=func-check
//ff:what @call 결과와 Response 필드 매칭 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// validateCallResult checks result/response field matching.
func validateCallResult(ctx string, spec *funcspec.FuncSpec, seq ssacparser.Sequence) []CrossError {
	if seq.Result != nil && len(spec.ResponseFields) == 0 {
		return []CrossError{{
			Rule:    "Func ↔ SSaC",
			Context: ctx,
			Message: "@result 있지만 Response 필드 없음",
			Level:   "ERROR",
		}}
	}
	if seq.Result == nil && len(spec.ResponseFields) > 0 {
		return []CrossError{{
			Rule:    "Func ↔ SSaC",
			Context: ctx,
			Message: "@result 없지만 Response 필드 존재 (반환값 무시)",
			Level:   "WARNING",
		}}
	}
	return nil
}
