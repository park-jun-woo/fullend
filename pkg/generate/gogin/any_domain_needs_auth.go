//ff:func feature=gen-gogin type=util control=iteration dimension=1
//ff:what 도메인 중 인증이 필요한 것이 있는지 확인한다

package gogin

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func anyDomainNeedsAuth(serviceFuncs []ssacparser.ServiceFunc, domains []string) bool {
	for _, d := range domains {
		if domainNeedsAuth(serviceFuncs, d) {
			return true
		}
	}
	return false
}
