//ff:func feature=crosscheck type=rule control=sequence
//ff:what @call Input 필드 수와 Request 필드 수 일치 검증
package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// validateCallInputCount checks that input field count matches request field count.
func validateCallInputCount(ctx string, spec *funcspec.FuncSpec, seq ssacparser.Sequence) []CrossError {
	inputCount := len(seq.Inputs)
	reqFieldCount := len(spec.RequestFields)
	if inputCount != reqFieldCount {
		return []CrossError{{
			Rule:    "Func ↔ SSaC",
			Context: ctx,
			Message: fmt.Sprintf("@call Inputs %d개, Request 필드 %d개 (불일치)", inputCount, reqFieldCount),
			Level:   "ERROR",
		}}
	}
	return nil
}
