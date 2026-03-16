//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 단일 SSaC 함수에서 @auth (action, resource) 쌍을 수집
package crosscheck

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

func collectAuthPairsFromFunc(fn ssacparser.ServiceFunc, ssacPairs map[[2]string]bool) {
	for _, seq := range fn.Sequences {
		if seq.Type == "auth" {
			ssacPairs[[2]string{seq.Action, seq.Resource}] = true
		}
	}
}
