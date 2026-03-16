//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what @call Input 소스 변수가 이전 시퀀스에서 정의되었는지 검증
package crosscheck

import (
	"fmt"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// validateCallSourceVars checks that source variables in inputs are defined in prior sequences.
func validateCallSourceVars(ctx string, seq ssacparser.Sequence, definedVars map[string]string) []CrossError {
	var errs []CrossError
	for _, value := range seq.Inputs {
		parts := strings.SplitN(value, ".", 2)
		source := parts[0]
		if source == "request" || source == "currentUser" {
			continue
		}
		if strings.HasPrefix(value, "\"") {
			continue
		}
		if ssacparser.IsLiteral(value) {
			continue
		}
		if _, ok := definedVars[source]; !ok {
			errs = append(errs, CrossError{
				Rule:    "Func ↔ SSaC",
				Context: ctx,
				Message: fmt.Sprintf("arg source %q 미정의", source),
				Level:   "WARNING",
			})
		}
	}
	return errs
}
