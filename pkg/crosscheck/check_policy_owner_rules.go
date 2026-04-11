//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkPolicyOwnerRules — 단일 정책의 allow rule에서 resource_owner 사용 시 @ownership 누락 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/rego"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkPolicyOwnerRules(rules []rego.AllowRule, ownedResources rule.StringSet) []CrossError {
	var errs []CrossError
	for _, r := range rules {
		if r.UsesOwner && !ownedResources[r.Resource] {
			errs = append(errs, CrossError{Rule: "X-30", Context: r.Resource, Level: "ERROR",
				Message: "allow rule uses resource_owner but no @ownership annotation for " + r.Resource})
		}
	}
	return errs
}
