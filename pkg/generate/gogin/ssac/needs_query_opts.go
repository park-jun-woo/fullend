//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=query-opts
//ff:what ServiceFunc의 시퀀스에서 query 옵션이 필요한지 확인
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func needsQueryOpts(sf ssacparser.ServiceFunc, st *rule.Ground) bool {
	for _, seq := range sf.Sequences {
		if hasQueryInput(seq.Inputs) {
			return true
		}
	}
	return false
}
