//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 단일 SSaC 함수에 @auth 시퀀스가 있는지 확인
package crosscheck

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

func funcHasAuth(sf ssacparser.ServiceFunc) bool {
	for _, seq := range sf.Sequences {
		if seq.Type == "auth" {
			return true
		}
	}
	return false
}
