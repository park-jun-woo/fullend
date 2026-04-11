//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what argUsesCurrentUser — Arg 목록에서 currentUser 소스 존재 여부 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func argUsesCurrentUser(args []ssac.Arg) bool {
	for _, arg := range args {
		if arg.Source == "currentUser" {
			return true
		}
	}
	return false
}
