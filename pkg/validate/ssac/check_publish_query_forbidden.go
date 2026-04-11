//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkPublishQueryForbidden — @publish Inputs에서 query 사용 감지
package ssac

import "github.com/park-jun-woo/fullend/pkg/validate"

func checkPublishQueryForbidden(file, funcName string, seqIdx int, inputs map[string]string) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, val := range inputs {
		if val == "query" {
			errs = append(errs, validate.ValidationError{
				Rule: "S-32", File: file, Func: funcName, SeqIdx: seqIdx, Level: "ERROR",
				Message: "@publish cannot use query",
			})
		}
	}
	return errs
}
