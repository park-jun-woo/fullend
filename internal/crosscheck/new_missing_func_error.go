//ff:func feature=crosscheck type=util control=sequence topic=func-check
//ff:what 미구현 func에 대한 CrossError 생성
package crosscheck

import (
	"fmt"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// newMissingFuncError creates a CrossError for a missing func implementation.
func newMissingFuncError(ctx, key, pkg, camelName string, seq ssacparser.Sequence) CrossError {
	skeleton := generateSkeleton(pkg, camelName, seq)
	return CrossError{
		Rule:       "Func ↔ SSaC",
		Context:    ctx,
		Message:    fmt.Sprintf("@call %s — 구현 없음", key),
		Level:      "ERROR",
		Suggestion: skeleton,
	}
}
