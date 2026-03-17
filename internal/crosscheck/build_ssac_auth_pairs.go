//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=policy-check
//ff:what SSaC 함수에서 authorize (action, resource) 쌍을 수집
package crosscheck

import ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"

func buildSSaCAuthPairs(funcs []ssacparser.ServiceFunc) map[[2]string]bool {
	ssacPairs := make(map[[2]string]bool)
	for _, fn := range funcs {
		collectAuthPairsFromFunc(fn, ssacPairs)
	}
	return ssacPairs
}
