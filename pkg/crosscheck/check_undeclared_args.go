//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkUndeclaredArgs — 단일 시퀀스의 arg source 미선언 검증
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func checkUndeclaredArgs(funcName string, args []ssac.Arg, declared map[string]bool) []CrossError {
	var errs []CrossError
	for _, arg := range args {
		if arg.Source == "" || implicitSources[arg.Source] || declared[arg.Source] {
			continue
		}
		errs = append(errs, CrossError{Rule: "X-47", Context: funcName, Level: "WARNING",
			Message: "@call arg source " + arg.Source + " not declared"})
	}
	return errs
}
