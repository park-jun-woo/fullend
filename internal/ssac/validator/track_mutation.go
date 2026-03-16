//ff:func feature=ssac-validate type=util control=sequence topic=type-resolve
//ff:what put/delete 모델의 mutation 여부를 기록한다
package validator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func trackMutation(seq parser.Sequence, mutated map[string]bool) {
	if seq.Model == "" {
		return
	}
	modelName := strings.SplitN(seq.Model, ".", 2)[0]
	mutated[modelName] = true
}
