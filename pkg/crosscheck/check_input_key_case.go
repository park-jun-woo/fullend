//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkInputKeyCase — SSaC input key가 대문자로 시작하는지 검증 (X-14)
package crosscheck

import (
	"unicode"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func checkInputKeyCase(funcName string, args []ssac.Arg) []CrossError {
	var errs []CrossError
	for _, arg := range args {
		if arg.Field != "" && len(arg.Field) > 0 && unicode.IsLower(rune(arg.Field[0])) {
			errs = append(errs, CrossError{Rule: "X-14", Context: funcName, Level: "ERROR",
				Message: "SSaC input key " + arg.Field + " should start with uppercase (Go convention)"})
		}
	}
	return errs
}
