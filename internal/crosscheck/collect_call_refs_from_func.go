//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what 단일 SSaC 함수에서 @call 모델 참조를 수집
package crosscheck

import (
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func collectCallRefsFromFunc(fn ssacparser.ServiceFunc, referenced map[string]bool) {
	for _, seq := range fn.Sequences {
		if seq.Type == "call" && seq.Model != "" {
			referenced[strings.ToLower(seq.Model)] = true
		}
	}
}
