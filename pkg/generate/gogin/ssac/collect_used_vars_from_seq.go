//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=http-handler
//ff:what 단일 시퀀스에서 참조되는 변수명을 수집
package ssac

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func collectUsedVarsFromSeq(seq ssacparser.Sequence, used map[string]bool) {
	if seq.Target != "" {
		used[rootVar(seq.Target)] = true
	}
	for _, val := range seq.Inputs {
		if isReservedInput(val) {
			continue
		}
		used[rootVar(val)] = true
	}
	for _, val := range seq.Fields {
		if !strings.HasPrefix(val, `"`) {
			used[rootVar(val)] = true
		}
	}
}
