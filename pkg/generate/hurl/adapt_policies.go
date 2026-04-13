//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what adaptPolicies — pkg/parser/rego.Policy → internal/policy.Policy 어댑터
package hurl

import (
	internalpolicy "github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/pkg/parser/rego"
)

// adaptPolicies converts parsed pkg rego policies to internal policy types.
// Phase007 임시 어댑터. 장기: 내부 함수가 rego.Policy 를 직접 받도록 전환.
func adaptPolicies(src []rego.Policy) []*internalpolicy.Policy {
	if len(src) == 0 {
		return nil
	}
	dst := make([]*internalpolicy.Policy, 0, len(src))
	for _, p := range src {
		rules := make([]internalpolicy.AllowRule, 0, len(p.Rules))
		for _, r := range p.Rules {
			rules = append(rules, internalpolicy.AllowRule{
				Resource:  r.Resource,
				Actions:   r.Actions,
				RoleValue: r.RoleValue,
				UsesOwner: r.UsesOwner,
				UsesRole:  r.UsesRole,
			})
		}
		owns := make([]internalpolicy.OwnershipMapping, 0, len(p.Ownerships))
		for _, o := range p.Ownerships {
			owns = append(owns, internalpolicy.OwnershipMapping{
				Resource:  o.Resource,
				Table:     o.Table,
				Column:    o.Column,
				JoinTable: o.JoinTable,
				JoinFK:    o.JoinFK,
			})
		}
		dst = append(dst, &internalpolicy.Policy{
			File:       p.File,
			Rules:      rules,
			Ownerships: owns,
			ClaimsRefs: p.ClaimsRefs,
		})
	}
	return dst
}
