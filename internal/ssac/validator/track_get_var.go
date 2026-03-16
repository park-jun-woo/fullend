//ff:func feature=ssac-validate type=util control=sequence topic=type-resolve
//ff:what @get 결과 변수와 모델명을 추적한다
package validator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func trackGetVar(seq parser.Sequence, getVars map[string]string, mutated map[string]bool) {
	if seq.Result == nil || seq.Model == "" {
		return
	}
	modelName := strings.SplitN(seq.Model, ".", 2)[0]
	getVars[seq.Result.Var] = modelName
	mutated[modelName] = false
}
