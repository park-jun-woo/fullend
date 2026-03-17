//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=currentuser
//ff:what 시퀀스의 Args에 currentUser 소스가 있는지 확인
package generator

import "github.com/park-jun-woo/fullend/internal/ssac/parser"

func hasCurrentUserArg(seq parser.Sequence) bool {
	for _, a := range seq.Args {
		if a.Source == "currentUser" {
			return true
		}
	}
	return false
}
