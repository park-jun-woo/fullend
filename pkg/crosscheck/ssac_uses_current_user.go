//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what ssacUsesCurrentUser — SSaC에서 currentUser 참조 여부 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

func ssacUsesCurrentUser(fs *fullend.Fullstack) bool {
	for _, fn := range fs.ServiceFuncs {
		if funcUsesCurrentUser(fn.Sequences) {
			return true
		}
	}
	return false
}
