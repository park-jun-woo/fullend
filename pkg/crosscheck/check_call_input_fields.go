//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallInputFields — @call input 필드가 FuncRequest에 있는지 검증 (X-43)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkCallInputFields(funcName string, seq ssac.Sequence, reqFields []string) []CrossError {
	reqSet := make(rule.StringSet, len(reqFields))
	for _, f := range reqFields {
		reqSet[f] = true
	}
	var errs []CrossError
	for _, arg := range seq.Args {
		if arg.Field != "" && !reqSet[arg.Field] {
			errs = append(errs, CrossError{Rule: "X-43", Context: funcName + "/" + seq.Model, Level: "ERROR",
				Message: "@call input " + arg.Field + " not in func request"})
		}
	}
	return errs
}
