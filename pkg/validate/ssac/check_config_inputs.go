//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkConfigInputs — Inputs에서 config.* 사용 검출 (S-31)
package ssac

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkConfigInputs(file, funcName string, seqIdx int, inputs map[string]string) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, val := range inputs {
		if strings.HasPrefix(val, "config.") {
			errs = append(errs, validate.ValidationError{
				Rule: "S-31", File: file, Func: funcName, SeqIdx: seqIdx, Level: "ERROR",
				Message: "config.* input forbidden — use os.Getenv() inside func",
			})
		}
	}
	return errs
}
