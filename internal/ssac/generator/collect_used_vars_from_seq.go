//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 단일 시퀀스에서 참조되는 변수명을 수집
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func collectUsedVarsFromSeq(seq parser.Sequence, used map[string]bool) {
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
