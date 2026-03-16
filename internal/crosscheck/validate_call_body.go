//ff:func feature=crosscheck type=rule control=sequence
//ff:what @call func의 본체 구현 여부 검증
package crosscheck

import "github.com/geul-org/fullend/internal/funcspec"

// validateCallBody checks that a func spec has a body implementation.
func validateCallBody(ctx string, spec *funcspec.FuncSpec) []CrossError {
	if !spec.HasBody {
		return []CrossError{{
			Rule:    "Func ↔ SSaC",
			Context: ctx,
			Message: "본체 미구현 (TODO)",
			Level:   "ERROR",
		}}
	}
	return nil
}
