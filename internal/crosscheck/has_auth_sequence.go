//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what SSaC 함수 목록에서 @auth 시퀀스 존재 여부를 확인
package crosscheck

import ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"

func hasAuthSequence(funcs []ssacparser.ServiceFunc) bool {
	for _, sf := range funcs {
		if funcHasAuth(sf) {
			return true
		}
	}
	return false
}
