//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkOwnershipAnnotation — resource_owner 참조인데 @ownership 없는 경우 경고 (X-30)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkOwnershipAnnotation(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.ParsedPolicies) == 0 {
		return nil
	}
	ownedResources := collectOwnedResources(fs)
	var errs []CrossError
	for _, p := range fs.ParsedPolicies {
		errs = append(errs, checkPolicyOwnerRules(p.Rules, ownedResources)...)
	}
	return errs
}
