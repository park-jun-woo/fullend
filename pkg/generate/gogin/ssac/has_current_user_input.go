//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=currentuser
//ff:what 시퀀스의 Inputs에 currentUser. 접두사가 있는지 확인
package ssac

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func hasCurrentUserInput(seq ssacparser.Sequence) bool {
	for _, val := range seq.Inputs {
		if strings.HasPrefix(val, "currentUser.") {
			return true
		}
	}
	return false
}
