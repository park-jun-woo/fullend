//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=currentuser
//ff:what 시퀀스의 Inputs에 currentUser. 접두사가 있는지 확인
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func hasCurrentUserInput(seq parser.Sequence) bool {
	for _, val := range seq.Inputs {
		if strings.HasPrefix(val, "currentUser.") {
			return true
		}
	}
	return false
}
