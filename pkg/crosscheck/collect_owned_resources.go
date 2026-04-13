//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectOwnedResources — 모든 Rego 정책에서 @ownership 선언된 리소스 수집
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func collectOwnedResources(fs *fullend.Fullstack) rule.StringSet {
	owned := make(rule.StringSet)
	for _, p := range fs.ParsedPolicies {
		for _, om := range p.Ownerships {
			owned[om.Resource] = true
		}
	}
	return owned
}
