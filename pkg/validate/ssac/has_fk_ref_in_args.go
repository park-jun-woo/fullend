//ff:func feature=rule type=util control=iteration dimension=1
//ff:what hasFKRefInArgs — Args/Inputs에서 FK 참조(다른 Model 변수)가 있는지 확인
package ssac

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func hasFKRefInArgs(seq parsessac.Sequence, declared map[string]bool, varTypes map[string]string, getModel string) bool {
	for _, arg := range seq.Args {
		ref := arg.Source
		if ref == "" || ref == "request" || ref == "currentUser" || ref == "query" || ref == "message" {
			continue
		}
		if declared[ref] && varTypes[ref] != getModel {
			return true
		}
	}
	for _, val := range seq.Inputs {
		if strings.HasPrefix(val, `"`) || parsessac.IsLiteral(val) {
			continue
		}
		ref := strings.SplitN(val, ".", 2)[0]
		if ref == "" || ref == "request" || ref == "currentUser" || ref == "query" || ref == "message" {
			continue
		}
		if declared[ref] && varTypes[ref] != getModel {
			return true
		}
	}
	return false
}
