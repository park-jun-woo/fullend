//ff:func feature=ssac-validate type=util control=iteration dimension=1 topic=type-resolve
//ff:what @get input이 FK 참조(다른 Model 변수)인지 확인
package validator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func hasForeignKeyRef(seq parser.Sequence, declared map[string]bool, varTypes map[string]string) bool {
	getModel := extractModel(seq.Model)
	for _, val := range seq.Inputs {
		if strings.HasPrefix(val, `"`) {
			continue
		}
		ref := rootVar(val)
		if ref == "request" || ref == "currentUser" || ref == "query" || ref == "message" || ref == "config" || ref == "" {
			continue
		}
		if declared[ref] && varTypes[ref] != getModel {
			return true
		}
	}
	return false
}
